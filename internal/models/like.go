package models

import (
	"time"

	"gorm.io/gorm"
)

// Like 点赞
type Like struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID    uint  `gorm:"not null;index:idx_user_target" json:"user_id"`
	ArticleID *uint `gorm:"index:idx_user_target;index:idx_article" json:"article_id,omitempty"` // 可为空
	WorkID    *uint `gorm:"index:idx_user_target;index:idx_work" json:"work_id,omitempty"`       // 可为空

	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Article *Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	Work    *Work    `gorm:"foreignKey:WorkID" json:"work,omitempty"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}
