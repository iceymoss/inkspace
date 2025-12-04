package service

import (
	"gorm.io/gorm"
	"mysite/internal/database"
	"mysite/internal/models"
)

type SettingService struct{}

func NewSettingService() *SettingService {
	return &SettingService{}
}

func (s *SettingService) Get(key string) (*models.Setting, error) {
	var setting models.Setting
	if err := database.DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (s *SettingService) Set(req *models.SettingRequest) (*models.Setting, error) {
	var setting models.Setting
	err := database.DB.Where("`key` = ?", req.Key).First(&setting).Error

	if err != nil {
		// 不存在，创建新记录
		setting = models.Setting{
			Key:         req.Key,
			Value:       req.Value,
			Type:        req.Type,
			Description: req.Description,
			Group:       req.Group,
			IsPublic:    req.IsPublic,
		}
		if err := database.DB.Create(&setting).Error; err != nil {
			return nil, err
		}
	} else {
		// 已存在，更新
		setting.Value = req.Value
		setting.Type = req.Type
		setting.Description = req.Description
		setting.Group = req.Group
		setting.IsPublic = req.IsPublic
		if err := database.DB.Save(&setting).Error; err != nil {
			return nil, err
		}
	}

	return &setting, nil
}

func (s *SettingService) Delete(key string) error {
	return database.DB.Where("`key` = ?", key).Delete(&models.Setting{}).Error
}

func (s *SettingService) GetAll() ([]*models.Setting, error) {
	var settings []*models.Setting
	if err := database.DB.Order("`group`, `key`").Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *SettingService) GetPublic() ([]*models.Setting, error) {
	var settings []*models.Setting
	if err := database.DB.Where("is_public = ?", true).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *SettingService) GetByGroup(group string) ([]*models.Setting, error) {
	var settings []*models.Setting
	if err := database.DB.Where("`group` = ?", group).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// BatchSet 批量设置
func (s *SettingService) BatchSet(settings map[string]string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		for key, value := range settings {
			var setting models.Setting
			err := tx.Where("`key` = ?", key).First(&setting).Error
			if err != nil {
				// 创建新记录
				setting = models.Setting{
					Key:   key,
					Value: value,
					Type:  "string",
				}
				if err := tx.Create(&setting).Error; err != nil {
					return err
				}
			} else {
				// 更新现有记录
				if err := tx.Model(&setting).Update("value", value).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
