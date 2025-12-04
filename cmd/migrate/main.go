package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"mysite/internal/config"
	"mysite/internal/database"
)

func main() {
	// 定义命令行参数
	var (
		action = flag.String("action", "migrate", "操作类型: migrate, drop, sync, health, tables")
		force  = flag.Bool("force", false, "强制执行（危险操作）")
	)
	flag.Parse()

	// 加载配置
	if err := config.Init(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库连接（不执行迁移）
	if err := initDBWithoutMigrate(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 执行对应的操作
	switch *action {
	case "migrate":
		handleMigrate()
	case "drop":
		handleDrop(*force)
	case "sync":
		handleSync()
	case "health":
		handleHealth()
	case "tables":
		handleTables()
	case "info":
		handleInfo()
	default:
		fmt.Printf("未知的操作: %s\n", *action)
		printUsage()
		os.Exit(1)
	}
}

// initDBWithoutMigrate 初始化数据库连接但不执行迁移
func initDBWithoutMigrate() error {
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	var err error
	database.DB, err = database.InitConnection(dsn)
	return err
}

// handleMigrate 执行数据库迁移
func handleMigrate() {
	fmt.Println("========================================")
	fmt.Println("开始数据库迁移...")
	fmt.Println("========================================")

	if err := database.MigrateDB(); err != nil {
		log.Fatalf("❌ 迁移失败: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("✅ 数据库迁移完成！")
	fmt.Println("========================================")
}

// handleDrop 删除所有表
func handleDrop(force bool) {
	if !force {
		fmt.Println("⚠️  警告: 此操作将删除所有表和数据！")
		fmt.Println("如果确定要执行，请使用 -force 参数")
		os.Exit(1)
	}

	fmt.Println("========================================")
	fmt.Println("⚠️  正在删除所有表...")
	fmt.Println("========================================")

	if err := database.DropAllTables(); err != nil {
		log.Fatalf("❌ 删除失败: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("✅ 所有表已删除")
	fmt.Println("========================================")
}

// handleSync 同步计数器
func handleSync() {
	fmt.Println("========================================")
	fmt.Println("开始同步计数器...")
	fmt.Println("========================================")

	if err := database.SyncCounters(); err != nil {
		log.Fatalf("❌ 同步失败: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("✅ 计数器同步完成！")
	fmt.Println("========================================")
}

// handleHealth 检查数据库健康状态
func handleHealth() {
	fmt.Println("========================================")
	fmt.Println("检查数据库健康状态...")
	fmt.Println("========================================")

	if err := database.CheckDatabaseHealth(); err != nil {
		fmt.Printf("❌ 健康检查失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("========================================")
	fmt.Println("✅ 数据库状态正常！")
	fmt.Println("========================================")
}

// handleTables 列出所有表
func handleTables() {
	fmt.Println("========================================")
	fmt.Println("数据库表列表")
	fmt.Println("========================================")

	tables, err := database.GetAllTables()
	if err != nil {
		log.Fatalf("❌ 获取表列表失败: %v", err)
	}

	for i, table := range tables {
		fmt.Printf("%d. %s\n", i+1, table)
	}

	fmt.Println("========================================")
	fmt.Printf("共 %d 个表\n", len(tables))
	fmt.Println("========================================")
}

// handleInfo 显示表信息
func handleInfo() {
	if flag.NArg() < 1 {
		fmt.Println("请指定表名: go run cmd/migrate/main.go -action=info <table_name>")
		os.Exit(1)
	}

	tableName := flag.Arg(0)
	fmt.Println("========================================")
	fmt.Printf("表信息: %s\n", tableName)
	fmt.Println("========================================")

	info, err := database.GetTableInfo(tableName)
	if err != nil {
		log.Fatalf("❌ 获取表信息失败: %v", err)
	}

	fmt.Printf("%-20s %-20s %-10s %-10s %-20s\n", "字段名", "类型", "允许空", "键", "默认值")
	fmt.Println(string(make([]byte, 80)))

	for _, row := range info {
		field := row["Field"]
		typ := row["Type"]
		null := row["Null"]
		key := row["Key"]
		def := row["Default"]

		fmt.Printf("%-20v %-20v %-10v %-10v %-20v\n", field, typ, null, key, def)
	}

	fmt.Println("========================================")
}

// printUsage 打印使用说明
func printUsage() {
	fmt.Println(`
数据库管理工具

用法:
  go run cmd/migrate/main.go -action=<action> [options]

操作类型:
  migrate   执行数据库迁移（创建表、索引等）
  drop      删除所有表（需要 -force 参数）
  sync      同步所有计数器字段
  health    检查数据库健康状态
  tables    列出所有表
  info      显示指定表的信息

选项:
  -force    强制执行危险操作

示例:
  # 执行迁移
  go run cmd/migrate/main.go -action=migrate

  # 同步计数器
  go run cmd/migrate/main.go -action=sync

  # 检查健康状态
  go run cmd/migrate/main.go -action=health

  # 列出所有表
  go run cmd/migrate/main.go -action=tables

  # 查看表信息
  go run cmd/migrate/main.go -action=info users

  # 删除所有表（危险！）
  go run cmd/migrate/main.go -action=drop -force
`)
}

