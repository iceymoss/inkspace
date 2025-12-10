package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Work struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Title         string         `gorm:"size:200;not null" json:"title" binding:"required"`
	Type          string         `gorm:"size:50;not null;default:'project';index:idx_type" json:"type"` // project, photography, video, etc.
	Metadata      string         `gorm:"type:text" json:"metadata"`                                     // JSON: 类型专属元数据
	DailyQuota    bool           `gorm:"default:false" json:"daily_quota"`                              // 是否受每日配额限制
	Description   string         `gorm:"type:text" json:"description"`
	Cover         string         `gorm:"size:255" json:"cover"`
	Images        string         `gorm:"type:text" json:"images"` // JSON array of image URLs
	Link          string         `gorm:"size:255" json:"link"`
	GithubURL     string         `gorm:"size:255" json:"github_url"`
	DemoURL       string         `gorm:"size:255" json:"demo_url"`
	TechStack     string         `gorm:"size:500" json:"tech_stack"` // 技术栈，逗号分隔
	AuthorID      uint           `gorm:"index:idx_author_id;not null;default:1" json:"author_id"`
	Author        *User          `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Sort          int            `gorm:"default:0;index:idx_status_sort" json:"sort"`
	ViewCount     int            `gorm:"default:0;not null" json:"view_count"`
	CommentCount  int            `gorm:"default:0;not null" json:"comment_count"`
	LikeCount     int            `gorm:"default:0;not null" json:"like_count"`
	FavoriteCount int            `gorm:"default:0;not null" json:"favorite_count"`
	Status        int            `gorm:"default:1;index:idx_status_sort" json:"status"` // 0: draft, 1: published, 2: pending, 3: rejected
	IsRecommend   bool           `gorm:"default:false" json:"is_recommend"`
	AuditMessage  string         `gorm:"type:text" json:"audit_message"` // 审核消息（审核通过或拒绝的原因）
}

// PhotoItem 单张照片信息（摄影作品中的一张照片）
type PhotoItem struct {
	URL         string             `json:"url"`         // 照片URL
	Description string             `json:"description"` // 照片描述
	Metadata    *PhotoItemMetadata `json:"metadata"`    // 照片参数
}

// PhotoItemMetadata 单张照片的EXIF参数
type PhotoItemMetadata struct {
	Camera       string `json:"camera,omitempty"`        // 相机型号
	Lens         string `json:"lens,omitempty"`          // 镜头
	FocalLength  string `json:"focal_length,omitempty"`  // 焦段
	Aperture     string `json:"aperture,omitempty"`      // 光圈
	ShutterSpeed string `json:"shutter_speed,omitempty"` // 快门速度
	ISO          string `json:"iso,omitempty"`           // ISO
}

// PhotographyMetadata 摄影作品（相册）元数据
type PhotographyMetadata struct {
	Location     string `json:"location,omitempty"`      // 拍摄地点
	ShootingDate string `json:"shooting_date,omitempty"` // 拍摄日期
	PhotoCount   int    `json:"photo_count"`             // 照片数量
}

type WorkRequest struct {
	Title       string                 `json:"title" binding:"required,max=200"`
	Type        string                 `json:"type" binding:"required,oneof=project photography"`
	Metadata    map[string]interface{} `json:"metadata"` // 动态元数据
	DailyQuota  bool                   `json:"daily_quota"`
	Description string                 `json:"description"`
	Cover       string                 `json:"cover"`
	Images      []PhotoItem            `json:"images"` // 对于项目：简单字符串；对于摄影：PhotoItem对象
	Link        string                 `json:"link" binding:"max=255"`
	GithubURL   string                 `json:"github_url"`
	DemoURL     string                 `json:"demo_url"`
	TechStack   string                 `json:"tech_stack"`
	// AuthorID 字段已移除，必须从token中提取，不允许前端传递
	Sort        int                    `json:"sort"`
	Status      int                    `json:"status"`
	IsRecommend bool                   `json:"is_recommend"`
}

type WorkResponse struct {
	ID            uint                   `json:"id"`
	Title         string                 `json:"title"`
	Type          string                 `json:"type"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"` // 动态元数据
	DailyQuota    bool                   `json:"daily_quota"`
	Description   string                 `json:"description"`
	Cover         string                 `json:"cover"`
	Images        []PhotoItem            `json:"images"` // 照片数组（包含参数）
	Link          string                 `json:"link"`
	GithubURL     string                 `json:"github_url"`
	DemoURL       string                 `json:"demo_url"`
	TechStack     string                 `json:"tech_stack"`
	AuthorID      uint                   `json:"author_id"`
	Author        *UserResponse          `json:"author,omitempty"`
	Sort          int                    `json:"sort"`
	ViewCount     int                    `json:"view_count"`
	CommentCount  int                    `json:"comment_count"`
	LikeCount     int                    `json:"like_count"`
	FavoriteCount int                    `json:"favorite_count"`
	Status        int                    `json:"status"`
	IsPublished   bool                   `json:"is_published"`      // 是否已发布（status=1），用于前端判断是否可以评论
	IsRecommend   bool                   `json:"is_recommend"`
	AuditMessage  string                 `json:"audit_message,omitempty"` // 审核消息（审核通过或拒绝的原因）
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

func (w *Work) ToResponse() *WorkResponse {
	resp := &WorkResponse{
		ID:            w.ID,
		Title:         w.Title,
		Type:          w.Type,
		DailyQuota:    w.DailyQuota,
		Description:   w.Description,
		Cover:         w.Cover,
		Link:          w.Link,
		GithubURL:     w.GithubURL,
		DemoURL:       w.DemoURL,
		TechStack:     w.TechStack,
		AuthorID:      w.AuthorID,
		Sort:          w.Sort,
		ViewCount:     w.ViewCount,
		CommentCount:  w.CommentCount,
		LikeCount:     w.LikeCount,
		FavoriteCount: w.FavoriteCount,
		Status:        w.Status,
		IsPublished:   w.Status == 1, // 只有已发布（status=1）的作品才允许评论
		IsRecommend:   w.IsRecommend,
		AuditMessage:  w.AuditMessage,
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}

	// 解析 Metadata
	if w.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(w.Metadata), &metadata); err == nil {
			resp.Metadata = metadata
		}
	}

	// 解析 Images
	if w.Images != "" {
		var images []PhotoItem
		if err := json.Unmarshal([]byte(w.Images), &images); err == nil {
			resp.Images = images
		}
	}

	// 处理 Author
	if w.Author != nil {
		resp.Author = w.Author.ToResponse()
	}

	return resp
}
