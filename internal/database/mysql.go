package database

import (
	"fmt"
	"time"

	"mysite/internal/config"
	"mysite/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
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
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// 命名策略
		NamingStrategy: nil, // 使用默认策略
		// 禁用外键约束（由应用层维护）
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 执行数据库迁移
	if err := MigrateDB(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 检查数据库健康状态
	if err := CheckDatabaseHealth(); err != nil {
		return fmt.Errorf("数据库健康检查失败: %w", err)
	}

	return nil
}

// HealthCheck 数据库健康检查（别名）
func HealthCheck() error {
	return CheckDatabaseHealth()
}

//// CheckDatabaseHealth 检查数据库健康状态
//func CheckDatabaseHealth() error {
//	sqlDB, err := DB.DB()
//	if err != nil {
//		return err
//	}
//	return sqlDB.Ping()
//}

func autoMigrate() error {
	return DB.AutoMigrate(
		// 核心表
		&models.User{},
		&models.Article{},
		&models.Category{},
		&models.Tag{},
		&models.Comment{},
		&models.Work{},
		&models.Link{},
		&models.Setting{},
		&models.Attachment{},
		&models.Like{},
		&models.ArticleFavorite{},
		&models.Favorite{},
		&models.UserFollow{},
		&models.Notification{},
		&models.Subscription{},
		&models.AdPosition{},
		&models.Advertisement{},
		&models.AdPlacement{},
		// 日志表
		&models.VisitLog{},
		&models.VisitLogSummary{},
	)
}
