package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/frostyeti/hyprship/apps/api/config"
	v1 "github.com/frostyeti/hyprship/apps/api/routes/v1"
	"github.com/frostyeti/hyprship/apps/api/telemetry"
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
	v1.RegisterRoutes(apiV1)

	addr := cfg.Addr
	if addr == "" {
		addr = fmt.Sprintf(":%d", cfg.Port)
	}

	slog.Info("Starting API server", "addr", addr, "env", cfg.Env)
	if err := r.Run(addr); err != nil {
		slog.Error("Server failed", "error", err)
	}
}
