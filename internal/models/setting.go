package models

import (
	"time"

	"gorm.io/gorm"
)

// Setting 系统配置表
type Setting struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Key         string         `gorm:"uniqueIndex;size:100;not null" json:"key"`     // 配置键
	Value       string         `gorm:"type:text" json:"value"`                       // 配置值
	Type        string         `gorm:"size:20;default:'string'" json:"type"`         // 类型: string, int, bool, json
	Description string         `gorm:"size:200" json:"description"`                  // 描述
	Group       string         `gorm:"size:50;default:'general';index" json:"group"` // 分组
	IsPublic    bool           `gorm:"default:false" json:"is_public"`               // 是否公开（前端可访问）
}

type SettingRequest struct {
	Key         string `json:"key" binding:"required,max=100"`
	Value       string `json:"value"`
	Type        string `json:"type" binding:"oneof=string int bool json"`
	Description string `json:"description" binding:"max=200"`
	Group       string `json:"group" binding:"max=50"`
	IsPublic    bool   `json:"is_public"`
}

type SettingResponse struct {
	ID          uint      `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Group       string    `json:"group"`
	IsPublic    bool      `json:"is_public"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 预定义的配置键
const (
	SettingSiteName        = "site_name"        // 网站名称
	SettingSiteDescription = "site_description" // 网站描述
	SettingSiteKeywords    = "site_keywords"    // 网站关键词
	SettingSiteICP         = "site_icp"         // 备案号
	SettingSiteCopyright   = "site_copyright"   // 版权信息
	SettingSiteLogo        = "site_logo"        // 网站Logo
	SettingSiteFavicon     = "site_favicon"     // 网站图标
	SettingCommentAudit          = "comment_audit"           // 评论是否需要审核
	SettingArticleCommentEnabled = "article_comment_enabled" // 是否开放文章评论
	SettingWorkCommentEnabled    = "work_comment_enabled"    // 是否开放作品评论
	SettingWorkAudit             = "work_audit"              // 作品是否需要审核
	SettingRegisterEnabled       = "register_enabled"        // 是否开放注册
	SettingUploadMaxSize         = "upload_max_size"          // 上传文件最大大小
	SettingCodeTheme       = "code_theme"       // Markdown 代码高亮主题
	SettingMarkdownTheme   = "markdown_theme"   // Markdown 主题风格（light/dark）
	SettingSiteTheme       = "site_theme"       // 网站整体主题（day/night/holiday/mourning）
)

func (s *Setting) ToResponse() *SettingResponse {
	return &SettingResponse{
		ID:          s.ID,
		Key:         s.Key,
		Value:       s.Value,
		Type:        s.Type,
		Description: s.Description,
		Group:       s.Group,
		IsPublic:    s.IsPublic,
		UpdatedAt:   s.UpdatedAt,
	}
}
