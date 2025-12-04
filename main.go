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

	// Setup user router (不包含管理API)
	r := router.SetupUserRouter()

	// Start server
	port := config.AppConfig.Server.Port
	log.Printf("用户服务启动在端口 %d...", port)
	log.Println("============================================")
	log.Println("提示: 建议使用以下命令启动服务:")
	log.Println("  用户服务: make dev 或 go run cmd/server/main.go")
	log.Println("  管理后台: make dev-admin 或 go run cmd/admin/main.go")
	log.Println("============================================")
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

