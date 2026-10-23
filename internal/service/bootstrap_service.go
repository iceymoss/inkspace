package service

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"time"
	"unicode/utf8"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrBootstrapUnavailable = errors.New("系统初始化不可用，请联系维护者")
	ErrBootstrapCompleted   = errors.New("系统已初始化，请登录")
	ErrBootstrapSecret      = errors.New("初始化密钥错误或未配置")
	ErrBootstrapConflict    = errors.New("用户名或邮箱已被占用，请更换")
	ErrBootstrapPassword    = errors.New("密码不能超过72个字节")
)

type BootstrapService struct{}

func NewBootstrapService() *BootstrapService { return &BootstrapService{} }

func bootstrapSecret() (string, bool) {
	if config.AppConfig == nil {
		return "", false
	}
	secret := config.AppConfig.Bootstrap.Secret
	return secret, len([]byte(secret)) >= 32
}

func (s *BootstrapService) Status() (*models.AdminBootstrapStatus, error) {
	if database.DB == nil {
		return nil, ErrBootstrapUnavailable
	}
	var state models.SystemInitialization
	if err := database.DB.First(&state, 1).Error; err != nil {
		return nil, ErrBootstrapUnavailable
	}
	var total int64
	if err := database.DB.Unscoped().Model(&models.User{}).Where("role = ?", "admin").Count(&total).Error; err != nil {
		return nil, ErrBootstrapUnavailable
	}
	result := &models.AdminBootstrapStatus{}
	switch state.State {
	case models.InitializationPending:
		if total > 0 {
			result.State = "unavailable"
		} else {
			result.State = "pending"
			_, result.SecretConfigured = bootstrapSecret()
		}
	case models.InitializationLegacyClosed, models.InitializationCompleted:
		var loginCapable int64
		if err := database.DB.Model(&models.User{}).Where("role = ? AND status = 1", "admin").Count(&loginCapable).Error; err != nil {
			return nil, ErrBootstrapUnavailable
		}
		if loginCapable > 0 {
			result.State = "initialized"
		} else {
			result.State = "unavailable"
		}
	default:
		result.State = "unavailable"
	}
	return result, nil
}

func (s *BootstrapService) Create(req *models.AdminBootstrapRequest) error {
	if database.DB == nil {
		return ErrBootstrapUnavailable
	}
	var current models.SystemInitialization
	if err := database.DB.First(&current, 1).Error; err != nil {
		return ErrBootstrapUnavailable
	}
	if current.State != models.InitializationPending {
		return ErrBootstrapCompleted
	}
	secret, configured := bootstrapSecret()
	if !configured || req.Secret == "" {
		return ErrBootstrapSecret
	}
	provided := sha256.Sum256([]byte(req.Secret))
	expected := sha256.Sum256([]byte(secret))
	if subtle.ConstantTimeCompare(provided[:], expected[:]) != 1 {
		return ErrBootstrapSecret
	}
	usernameLength := utf8.RuneCountInString(req.Username)
	emailLength := utf8.RuneCountInString(req.Email)
	passwordLength := utf8.RuneCountInString(req.Password)
	if req.Password != req.ConfirmPassword || emailLength > 100 || usernameLength < 3 || usernameLength > 50 || passwordLength < 6 || passwordLength > 50 {
		return errors.New("初始化信息无效")
	}
	if len([]byte(req.Password)) > 72 {
		return ErrBootstrapPassword
	}
	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return ErrBootstrapUnavailable
	}
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var state models.SystemInitialization
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&state, 1).Error; err != nil {
			return ErrBootstrapUnavailable
		}
		if state.State != models.InitializationPending {
			return ErrBootstrapCompleted
		}
		var admins int64
		if err := tx.Unscoped().Model(&models.User{}).Where("role = ?", "admin").Count(&admins).Error; err != nil {
			return ErrBootstrapUnavailable
		}
		if admins != 0 {
			return ErrBootstrapUnavailable
		}
		user := models.User{Username: req.Username, Email: req.Email, Nickname: req.Username, Password: password, Role: "admin", Status: 1}
		if err := tx.Create(&user).Error; err != nil {
			var mysqlErr *mysqlDriver.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return ErrBootstrapConflict
			}
			return ErrBootstrapUnavailable
		}
		now := time.Now()
		update := tx.Model(&state).Where("state = ?", models.InitializationPending).
			Updates(map[string]interface{}{"state": models.InitializationCompleted, "completed_at": now})
		if update.Error != nil || update.RowsAffected != 1 {
			return ErrBootstrapUnavailable
		}
		return nil
	})
}
