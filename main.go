package main

import (
	"fmt"
	"log"

	"mysite/internal/config"
	"mysite/internal/database"
	"mysite/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis
	if err := database.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Set Gin mode
	gin.SetMode(config.AppConfig.Server.Mode)

	// Setup router
	r := router.SetupRouter()

	// Start server
	port := config.AppConfig.Server.Port
	log.Printf("Server starting on port %d...", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

