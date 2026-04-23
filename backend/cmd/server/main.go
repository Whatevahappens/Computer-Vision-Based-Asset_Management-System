package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/router"
	"backend/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()

	middleware.SetJWTSecret(cfg.JWTSecret)

	service.SetCVServiceURL(cfg.CVServiceURL)

	// Share config with handlers
	handler.SetConfig(cfg)

	database.Connect(cfg)

	database.Migrate()

	database.SeedAdmin(cfg)

	gin.SetMode(cfg.GinMode)
	r := gin.Default()

	router.SetupRoutes(r)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
