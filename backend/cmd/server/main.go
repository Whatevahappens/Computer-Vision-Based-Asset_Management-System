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
	// Load config
	cfg := config.Load()

	// Set JWT secret
	middleware.SetJWTSecret(cfg.JWTSecret)

	// Set CV service URL
	service.SetCVServiceURL(cfg.CVServiceURL)

	// Share config with handlers
	handler.SetConfig(cfg)

	// Connect database
	database.Connect(cfg)

	// Run migrations
	database.Migrate()

	// Seed admin user
	database.SeedAdmin(cfg)

	// Setup Gin
	gin.SetMode(cfg.GinMode)
	r := gin.Default()

	// Setup routes
	router.SetupRoutes(r)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
