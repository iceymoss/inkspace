package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/gorm"
)

const initializationLock = "inkspace_system_initialization_v1"

// WithInitializationLock serializes pre-migration checks and schema changes across all binaries.
func WithInitializationLock(action func(*gorm.DB) error) error {
	return withInitializationLock(DB, action)
}

func withInitializationLock(db *gorm.DB, action func(*gorm.DB) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	return db.WithContext(ctx).Connection(func(conn *gorm.DB) (err error) {
		// Connection pins one SQL connection; WithContext would clone it and lose the pin.
		var acquired int
		if err := conn.Session(&gorm.Session{NewDB: true}).Raw("SELECT GET_LOCK(?, 300)", initializationLock).Scan(&acquired).Error; err != nil {
			return err
		}
		if acquired != 1 {
			return errors.New("无法获取初始化互斥锁")
		}
		defer func() {
			var released int
			releaseErr := conn.Session(&gorm.Session{NewDB: true}).Raw("SELECT RELEASE_LOCK(?)", initializationLock).Scan(&released).Error
			if releaseErr != nil {
				err = errors.Join(err, fmt.Errorf("释放初始化互斥锁失败: %w", releaseErr))
			} else if released != 1 {
				err = errors.Join(err, fmt.Errorf("释放初始化互斥锁失败: result=%d", released))
			}
		}()
		return action(conn)
	})
}

// checkInitializationConsistency validates account invariants after users has been migrated.
func checkInitializationConsistency(conn *gorm.DB) error {
	var state models.SystemInitialization
	if err := conn.Session(&gorm.Session{NewDB: true}).First(&state, 1).Error; err != nil {
		return fmt.Errorf("读取初始化状态失败: %w", err)
	}
	var administrators int64
	if err := conn.Session(&gorm.Session{NewDB: true}).Unscoped().Model(&models.User{}).
		Where("role = ?", "admin").Count(&administrators).Error; err != nil {
		return fmt.Errorf("检查管理员账号失败: %w", err)
	}
	if state.State == models.InitializationPending && administrators != 0 {
		return errors.New("待初始化系统中存在历史管理员，需要人工检查")
	}
	if state.State == models.InitializationLegacyClosed || state.State == models.InitializationCompleted {
		var loginCapable int64
		if err := conn.Session(&gorm.Session{NewDB: true}).Model(&models.User{}).
			Where("role = ? AND status = 1", "admin").Count(&loginCapable).Error; err != nil {
			return fmt.Errorf("检查可登录管理员失败: %w", err)
		}
		if loginCapable == 0 {
			log.Printf("警告: 系统初始化状态为 %s，但没有可登录管理员；公开初始化入口保持关闭", state.State)
		}
	}
	return nil
}

// PrepareInitialization records installation provenance before AutoMigrate creates users.
// A missing or corrupt state row in an existing state table is never repaired automatically.
func PrepareInitialization() error {
	return WithInitializationLock(prepareInitialization)
}

func prepareInitialization(conn *gorm.DB) error {
	var tables []string
	if err := conn.Session(&gorm.Session{NewDB: true}).Raw(`SELECT table_name FROM information_schema.tables
		WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE'`).Scan(&tables).Error; err != nil {
		return fmt.Errorf("检查安装来源失败: %w", err)
	}
	stateTable := false
	legacy := false
	for _, table := range tables {
		if table == "system_initializations" {
			stateTable = true
		} else {
			legacy = true
		}
	}
	if stateTable {
		var records []models.SystemInitialization
		if err := conn.Session(&gorm.Session{NewDB: true}).Find(&records).Error; err != nil {
			return fmt.Errorf("读取初始化状态失败: %w", err)
		}
		if len(records) != 1 || records[0].ID != 1 {
			return errors.New("初始化状态记录异常，需要人工检查")
		}
		switch records[0].State {
		case models.InitializationPending, models.InitializationLegacyClosed, models.InitializationCompleted:
			return nil
		default:
			return errors.New("未知初始化状态，需要人工检查")
		}
	}
	if err := conn.AutoMigrate(&models.SystemInitialization{}); err != nil {
		return fmt.Errorf("创建初始化状态表失败: %w", err)
	}
	state := models.InitializationPending
	if legacy {
		state = models.InitializationLegacyClosed
	}
	if err := conn.Session(&gorm.Session{NewDB: true}).Create(&models.SystemInitialization{ID: 1, State: state}).Error; err != nil {
		return fmt.Errorf("保存初始化状态失败: %w", err)
	}
	return nil
}
