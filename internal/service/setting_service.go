package service

import (
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/gorm"
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
		updateData := map[string]interface{}{
			"value":       req.Value,
			"type":        req.Type,
			"description": req.Description,
			"group":       req.Group,
			"is_public":   req.IsPublic,
		}
		if err := database.DB.Model(&models.Setting{}).
			Where("`key` = ?", req.Key).
			Updates(updateData).Error; err != nil {
			return nil, err
		}
		// 重新加载以返回最新数据
		if err := database.DB.Where("`key` = ?", req.Key).First(&setting).Error; err != nil {
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

			// 根据 key 确定 group 和 is_public
			group := setting.Group       // 默认使用原有分组
			isPublic := setting.IsPublic // 默认使用原有公开设置

			// 预定义的配置分组（仅在创建新记录时设置）
			if err != nil {
				group = "general"
				isPublic = false

				if key == models.SettingSiteName || key == models.SettingSiteDescription ||
					key == models.SettingSiteKeywords || key == models.SettingSiteICP ||
					key == models.SettingSiteCopyright || key == models.SettingSiteLogo ||
					key == models.SettingSiteFavicon || key == "link_apply_email" ||
					key == "link_apply_title" || key == "link_apply_description" ||
					key == "admin_backend_url" {
					group = "site"
					isPublic = true
				} else if key == models.SettingCommentAudit || key == models.SettingRegisterEnabled ||
					key == models.SettingArticleCommentEnabled || key == models.SettingWorkCommentEnabled ||
					key == models.SettingWorkAudit {
					group = "feature"
					isPublic = false
				} else if key == models.SettingCodeTheme || key == models.SettingMarkdownTheme {
					group = "markdown"
					isPublic = true // Markdown 相关设置需要公开，前端才能使用
				} else if key == models.SettingSiteTheme {
					group = "theme"
					isPublic = true // 主题设置需要公开，前端才能使用
				} else if key == "holiday_type" || key == "holiday_bg_primary" ||
					key == "holiday_bg_secondary" || key == "holiday_text_primary" ||
					key == "holiday_primary" {
					group = "theme"
					isPublic = true // 节假日主题设置需要公开，前端才能使用
				}
			} else {
				// 更新现有记录时，如果是 code_theme 或 site_theme，确保设置正确
				if key == models.SettingCodeTheme {
					group = "markdown"
					isPublic = true
				} else if key == models.SettingSiteTheme {
					group = "theme"
					isPublic = true
				} else if key == "holiday_type" || key == "holiday_bg_primary" ||
					key == "holiday_bg_secondary" || key == "holiday_text_primary" ||
					key == "holiday_primary" {
					group = "theme"
					isPublic = true
				} else if key == "home_carousel" {
					group = "carousel"
					isPublic = true
				} else if key == "about_page" {
					group = "about"
					isPublic = true
				} else if key == "about_page" {
					group = "about"
					isPublic = true
				}
			}

			if err != nil {
				// 创建新记录
				setting = models.Setting{
					Key:      key,
					Value:    value,
					Type:     "string",
					Group:    group,
					IsPublic: isPublic,
				}
				if err := tx.Create(&setting).Error; err != nil {
					return err
				}
			} else {
				// 更新现有记录（保留原有的 group 和 is_public，除非是 code_theme 或 site_theme）
				updateData := map[string]interface{}{
					"value": value,
				}
				if key == models.SettingCodeTheme || key == models.SettingSiteTheme {
					updateData["group"] = group
					updateData["is_public"] = isPublic
				}
				if err := tx.Model(&models.Setting{}).
					Where("`key` = ?", key).
					Updates(updateData).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
