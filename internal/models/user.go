package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Username      string         `gorm:"uniqueIndex;size:50;not null" json:"username" binding:"required,min=3,max=50"`
	Password      string         `gorm:"size:255;not null" json:"-"`
	Email         string         `gorm:"uniqueIndex;size:100;not null" json:"email" binding:"required,email"`
	Nickname      string         `gorm:"size:50" json:"nickname"`
	Avatar        string         `gorm:"size:255" json:"avatar"`
	Bio           string         `gorm:"size:500" json:"bio"`
	Role          string         `gorm:"size:20;default:'user';index:idx_role_status" json:"role"` // admin, user
	Status        int            `gorm:"default:1;index:idx_role_status" json:"status"`            // 1: active, 0: inactive
	LastLoginAt   *time.Time     `gorm:"type:datetime(3)" json:"last_login_at"`
	LastLoginIP    string         `gorm:"size:50" json:"last_login_ip"`
	ArticleCount   int            `gorm:"default:0;not null" json:"article_count"`
	WorkCount      int            `gorm:"default:0;not null" json:"work_count"`      // 作品数
	CommentCount   int            `gorm:"default:0;not null" json:"comment_count"`
	FollowingCount int            `gorm:"default:0;not null" json:"following_count"` // 关注数
	FollowerCount  int            `gorm:"default:0;not null" json:"follower_count"`  // 粉丝数
	FavoriteCount  int            `gorm:"default:0;not null" json:"favorite_count"`  // 收藏数
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"max=50"`
}

type UserUpdateRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Email    string `json:"email" binding:"email"`
	Bio      string `json:"bio" binding:"max=500"`
	Avatar   string `json:"avatar"`
}

type PasswordChangeRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// UserResponse 用户响应（完整信息，用于自己的profile）
type UserResponse struct {
	ID             uint      `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Nickname       string    `json:"nickname"`
	Avatar         string    `json:"avatar"`
	Bio            string    `json:"bio"`
	Role           string    `json:"role"`
	Status         int       `json:"status"`
	ArticleCount   int       `json:"article_count"`
	WorkCount      int       `json:"work_count"`
	CommentCount   int       `json:"comment_count"`
	FollowingCount int       `json:"following_count"`
	FollowerCount  int       `json:"follower_count"`
	FavoriteCount  int       `json:"favorite_count"`
	CreatedAt      time.Time `json:"created_at"`
}

// PublicUserResponse 公开用户响应（用于查看他人主页，不包含敏感信息）
type PublicUserResponse struct {
	ID             uint      `json:"id"`
	Username       string    `json:"username"`
	Nickname       string    `json:"nickname"`
	Avatar         string    `json:"avatar"`
	Bio            string    `json:"bio"`
	ArticleCount   int       `json:"article_count"`
	WorkCount      int       `json:"work_count"`
	CommentCount   int       `json:"comment_count"`
	FollowingCount int       `json:"following_count"`
	FollowerCount  int       `json:"follower_count"`
	FavoriteCount  int       `json:"favorite_count"`
	CreatedAt      time.Time `json:"created_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		Nickname:       u.Nickname,
		Avatar:         u.Avatar,
		Bio:            u.Bio,
		Role:           u.Role,
		Status:         u.Status,
		ArticleCount:   u.ArticleCount,
		WorkCount:      u.WorkCount,
		CommentCount:   u.CommentCount,
		FollowingCount: u.FollowingCount,
		FollowerCount:  u.FollowerCount,
		FavoriteCount:  u.FavoriteCount,
		CreatedAt:      u.CreatedAt,
	}
}

// ToPublicResponse 转换为公开响应（不包含Email、Role、Status等敏感信息）
func (u *User) ToPublicResponse() *PublicUserResponse {
	return &PublicUserResponse{
		ID:             u.ID,
		Username:       u.Username,
		Nickname:       u.Nickname,
		Avatar:         u.Avatar,
		Bio:            u.Bio,
		ArticleCount:   u.ArticleCount,
		WorkCount:      u.WorkCount,
		CommentCount:   u.CommentCount,
		FollowingCount: u.FollowingCount,
		FollowerCount:  u.FollowerCount,
		FavoriteCount:  u.FavoriteCount,
		CreatedAt:      u.CreatedAt,
	}
}

