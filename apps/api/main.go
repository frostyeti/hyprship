package main

import (
	"fmt"
	"log"

	"github.com/frostyeti/hyprship/apps/api/config"
	v1 "github.com/frostyeti/hyprship/apps/api/routes/v1"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.Env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	v1.RegisterRoutes(apiV1)

	addr := cfg.Addr
	if addr == "" {
		addr = fmt.Sprintf(":%d", cfg.Port)
	}

	log.Printf("Starting API server on %s (env: %s)", addr, cfg.Env)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
