package main

import (
	"fmt"
	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载管理后台配置
	// 尝试加载admin.yaml，如果不存在则使用默认config.yaml
	if err := config.InitWithFile("admin"); err != nil {
		log.Printf("警告: 无法加载 admin.yaml，使用默认配置: %v\n", err)
		if err := config.Init(); err != nil {
			log.Fatalf("加载配置失败: %v", err)
		}
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化Redis
	if err := database.InitRedis(); err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}

	// 数据库初始化完成（Init()中已经包含了健康检查）
	log.Println("数据库初始化完成")

	// 设置Gin模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 设置管理后台路由
	r := router.SetupAdminRouter()

	// 获取端口（优先使用admin配置，否则使用server配置+1）
	port := config.AppConfig.Server.Port + 1
	if config.AppConfig.Admin.Port > 0 {
		port = config.AppConfig.Admin.Port
	}

	log.Printf("管理后台服务启动在端口 %d...\n", port)
	log.Println("===========================================")
	log.Printf("管理后台API: http://localhost:%d\n", port)
	log.Println("登录地址: http://localhost:3001/admin/login")
	log.Println("默认账号: admin / admin123")
	log.Println("===========================================")

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
