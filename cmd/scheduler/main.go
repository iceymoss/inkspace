package main

import (
	"log"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/scheduler"
	"github.com/iceymoss/inkspace/internal/utils"
)

func main() {
	// 初始化日志
	utils.InitLogger()

	log.Println("========================================")
	log.Println("启动定时任务调度器...")
	log.Println("========================================")

	// 加载配置
	if err := config.Init(); err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}

	// 初始化Redis
	if err := database.InitRedis(); err != nil {
		log.Fatalf("❌ Redis连接失败: %v", err)
	}

	log.Println("✅ 数据库和Redis连接成功")

	// 创建调度器
	sched := scheduler.NewScheduler()

	// 注册任务
	sched.RegisterTask("hot_articles", scheduler.NewHotArticlesTask(), 3*time.Minute)
	sched.RegisterTask("hot_works", scheduler.NewHotWorksTask(), 3*time.Minute)
	// 注册榜单生成任务（每天检查一次，在周日、月初、年初生成榜单）
	sched.RegisterTask("article_rank", scheduler.NewArticleRankTask(), 24*time.Hour)

	log.Println("========================================")
	log.Println("✅ 定时任务调度器启动成功")
	log.Println("已注册的任务:")
	for name, info := range sched.GetTasks() {
		log.Printf("  - %s (间隔: %s)", name, info.Interval)
	}
	log.Println("========================================")

	// 启动调度器
	sched.Start()

	// 阻止程序退出
	select {}
}
