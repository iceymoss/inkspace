package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	ArticleID  *uint          `gorm:"index:idx_article_id" json:"article_id"`
	Article    *Article       `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	WorkID     *uint          `gorm:"index:idx_work_id" json:"work_id"`
	Work       *Work          `gorm:"foreignKey:WorkID" json:"work,omitempty"`
	UserID     uint           `gorm:"index:idx_user_id" json:"user_id"`
	User       *User          `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	Content    string         `gorm:"type:text;not null" json:"content" binding:"required"`
	ParentID   *uint          `gorm:"index:idx_parent_id" json:"parent_id"`
	Parent     *Comment       `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"parent,omitempty"`
	RootID     *uint          `gorm:"index:idx_root_id" json:"root_id"` // 根评论ID，方便查询
	ReplyToID  *uint          `json:"reply_to_id"` // 回复的评论ID
	Nickname   string         `gorm:"size:50" json:"nickname"`
	Email      string         `gorm:"size:100" json:"email"`
	Website    string         `gorm:"size:200" json:"website"`
	IP         string         `gorm:"size:50" json:"ip"`
	UserAgent  string         `gorm:"size:255" json:"user_agent"`
	Status     int            `gorm:"default:1;index:idx_status_created" json:"status"` // 1: approved, 0: pending, -1: rejected
	LikeCount  int            `gorm:"default:0;not null" json:"like_count"`
	ReplyCount int            `gorm:"default:0;not null" json:"reply_count"`
}

type CommentRequest struct {
	ArticleID *uint  `json:"article_id"`
	WorkID    *uint  `json:"work_id"`
	Content   string `json:"content" binding:"required,max=500"`
	ParentID  *uint  `json:"parent_id"`
	Nickname  string `json:"nickname" binding:"omitempty,max=50"`
	Email     string `json:"email" binding:"omitempty,email,max=100"`
	Website   string `json:"website" binding:"omitempty,max=200"`
}

type CommentListQuery struct {
	Page      int   `form:"page,default=1"`
	PageSize  int   `form:"page_size,default=10"`
	ArticleID *uint `form:"article_id"`
	WorkID    *uint `form:"work_id"`
	UserID    uint  `form:"user_id"`
	Status    *int  `form:"status"`
	ShowAll   bool  `form:"show_all"`   // 是否显示所有状态的评论（管理后台使用）
	Type      string `form:"type"`      // 评论类型：'article' 只显示文章评论，'work' 只显示作品评论
}

type CommentResponse struct {
	ID         uint              `json:"id"`
	ArticleID  *uint             `json:"article_id"`
	Article    *ArticleResponse  `json:"article,omitempty"`
	WorkID     *uint             `json:"work_id"`
	Work       *WorkResponse     `json:"work,omitempty"`
	UserID     uint              `json:"user_id"`
	User       *UserResponse     `json:"user,omitempty"`
	Content    string            `json:"content"`
	ParentID   *uint             `json:"parent_id"`
	RootID     *uint             `json:"root_id"`
	Nickname   string            `json:"nickname"`
	Email      string            `json:"email"`
	Website    string            `json:"website"`
	Status     int               `json:"status"`
	LikeCount  int               `json:"like_count"`
	ReplyCount int               `json:"reply_count"`
	CreatedAt  time.Time         `json:"created_at"`
	Replies    []CommentResponse `json:"replies,omitempty"`
}

func (c *Comment) ToResponse() *CommentResponse {
	resp := &CommentResponse{
		ID:         c.ID,
		ArticleID:  c.ArticleID,
		WorkID:     c.WorkID,
		UserID:     c.UserID,
		Content:    c.Content,
		ParentID:   c.ParentID,
		RootID:     c.RootID,
		Nickname:   c.Nickname,
		Email:      c.Email,
		Website:    c.Website,
		Status:     c.Status,
		LikeCount:  c.LikeCount,
		ReplyCount: c.ReplyCount,
		CreatedAt:  c.CreatedAt,
	}

	if c.User != nil {
		resp.User = c.User.ToResponse()
	}

	if c.Article != nil {
		// 只返回文章的基本信息（标题等），不返回完整内容
		articleResp := c.Article.ToResponse()
		articleResp.Content = "" // 不返回文章内容
		resp.Article = articleResp
	}

	if c.Work != nil {
		// 简单返回基本信息，不需要完整的 ToResponse
		// 完整的转换由 work_handler 处理
		var images []PhotoItem
		if c.Work.Images != "" {
			json.Unmarshal([]byte(c.Work.Images), &images)
		}
		resp.Work = &WorkResponse{
			ID:           c.Work.ID,
			Title:        c.Work.Title,
			Description:  c.Work.Description,
			Cover:        c.Work.Cover,
			Images:       images,
			CommentCount: c.Work.CommentCount,
		}
	}

	return resp
}

