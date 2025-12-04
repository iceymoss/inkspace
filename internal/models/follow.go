package models

import (
	"time"

	"gorm.io/gorm"
)

// UserFollow 用户关注关系表
type UserFollow struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	FollowerID uint           `gorm:"index:idx_follower_following,priority:1;index:idx_follower_id;not null" json:"follower_id"` // 关注者ID（粉丝）
	Follower   *User          `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"follower,omitempty"`
	FollowingID uint          `gorm:"index:idx_follower_following,priority:2;index:idx_following_id;not null" json:"following_id"` // 被关注者ID
	Following  *User          `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE" json:"following,omitempty"`
}

// 组合唯一索引：防止重复关注
func (UserFollow) TableName() string {
	return "user_follows"
}

type FollowRequest struct {
	FollowingID uint `json:"following_id" binding:"required"` // 要关注的用户ID
}

type FollowListQuery struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

type FollowResponse struct {
	ID          uint          `json:"id"`
	FollowerID  uint          `json:"follower_id"`
	FollowingID uint          `json:"following_id"`
	User        *UserResponse `json:"user,omitempty"` // 用户信息（根据查询类型返回follower或following）
	CreatedAt   time.Time     `json:"created_at"`
}

// FollowStats 关注统计
type FollowStats struct {
	FollowingCount int  `json:"following_count"` // 关注数
	FollowerCount  int  `json:"follower_count"`  // 粉丝数
	IsFollowing    bool `json:"is_following"`    // 当前用户是否已关注
	IsFollower     bool `json:"is_follower"`     // 对方是否关注了当前用户（互关）
}

