package main

import (
	"fmt"
	"log"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/router"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	utils.InitLogger()

	// 加载配置
	if err := config.Init(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化Redis
	if err := database.InitRedis(); err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}

	// 数据库迁移（Init()中已经包含了健康检查）
	log.Println("数据库初始化完成")

	// 设置Gin模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 设置用户服务路由
	r := router.SetupUserRouter()

	// 启动服务器
	port := config.AppConfig.Server.Port
	log.Printf("用户服务启动在端口 %d...\n", port)
	log.Println("============================================")
	log.Println("提示: 建议使用以下命令启动服务:")
	log.Println("  用户服务: go run cmd/server/main.go")
	log.Println("  管理后台: go run cmd/admin/main.go")
	log.Println("============================================")

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
