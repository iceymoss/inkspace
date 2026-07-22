package service

import (
	"errors"
	"time"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserAppearanceService struct{}

func NewUserAppearanceService() *UserAppearanceService {
	return &UserAppearanceService{}
}

func (s *UserAppearanceService) Get(userID uint) (*models.UserAppearanceResponse, error) {
	var appearance models.UserAppearance
	if err := database.DB.Where("user_id = ?", userID).First(&appearance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.DefaultUserAppearanceResponse(), nil
		}
		return nil, err
	}

	return appearance.ToResponse(), nil
}

func (s *UserAppearanceService) Save(userID uint, req *models.UserAppearanceRequest) (*models.UserAppearanceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	appearance := models.UserAppearance{
		UserID:      userID,
		UITheme:     req.UITheme,
		ColorScheme: req.ColorScheme,
	}
	if err := database.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"ui_theme":     req.UITheme,
			"color_scheme": req.ColorScheme,
			"updated_at":   time.Now(),
		}),
	}).Create(&appearance).Error; err != nil {
		return nil, err
	}

	return appearance.ToResponse(), nil
}
