package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/frostyeti/hyprship/apps/api/config"
	"github.com/frostyeti/hyprship/apps/api/internal/crypto"
	"github.com/frostyeti/hyprship/apps/api/internal/db"
	v1 "github.com/frostyeti/hyprship/apps/api/internal/routes/v1"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/frostyeti/hyprship/apps/api/internal/svc/identity"
	"github.com/frostyeti/hyprship/apps/api/internal/telemetry"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger := telemetry.InitLogger(cfg)
	slog.SetDefault(logger)

	ctx := context.Background()
	shutdown, err := telemetry.InitOtel(ctx, cfg)
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry", "error", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			slog.Error("Failed to shutdown OpenTelemetry", "error", err)
		}
	}()

	// Init Database
	sqlDB, err := db.Connect(cfg.Database)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	// Setup Stores
	storeFactory, err := stores.NewStoreFactory(cfg.Database.Driver, sqlDB)
	if err != nil {
		slog.Error("Failed to initialize stores", "error", err)
		os.Exit(1)
	}

	// Setup Services
	identitySvc := identity.NewIdentityService(
		cfg,
		storeFactory.UserStore,
		storeFactory.UserPasswordAuthStore,
		storeFactory.UserSessionStore,
		storeFactory.UserPasswordHistoryStore,
		crypto.GetPasswordHasher("pbkdf2"),
	)

	if cfg.Env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(sloggin.New(logger))
	r.Use(gin.Recovery())

	if cfg.Otel != nil && cfg.Otel.Enabled {
		r.Use(otelgin.Middleware("hyprship-api"))
	}

	apiV1 := r.Group("/api/v1")
	v1.RegisterRoutes(apiV1, identitySvc, storeFactory.UserStore, storeFactory.RoleStore, storeFactory.UserSessionStore)

	addr := cfg.Addr
	if addr == "" {
		addr = fmt.Sprintf(":%d", cfg.Port)
	}

	slog.Info("Starting API server", "addr", addr, "env", cfg.Env)

	// Graceful shutdown
	go func() {
		if err := r.Run(addr); err != nil {
			slog.Error("Server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")
}
