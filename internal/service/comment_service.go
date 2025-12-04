package service

import (
	"errors"

	"mysite/internal/database"
	"mysite/internal/models"

	"gorm.io/gorm"
)

type CommentService struct{}

func NewCommentService() *CommentService {
	return &CommentService{}
}

func (s *CommentService) Create(req *models.CommentRequest, userID uint) (*models.Comment, error) {
	// Check if article exists
	var article models.Article
	if err := database.DB.First(&article, req.ArticleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	// Check if parent comment exists
	if req.ParentID != nil {
		var parent models.Comment
		if err := database.DB.First(&parent, *req.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父评论不存在")
			}
			return nil, err
		}
	}

	comment := &models.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID,
		Content:   req.Content,
		ParentID:  req.ParentID,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Website:   req.Website,
		Status:    1, // Auto approve
	}

	if err := database.DB.Create(comment).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) Delete(id uint, userID uint, role string) error {
	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		return err
	}

	// Check permission
	if role != "admin" && comment.UserID != userID {
		return errors.New("无权限删除")
	}

	return database.DB.Delete(&comment).Error
}

func (s *CommentService) GetList(query *models.CommentListQuery) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	db := database.DB.Model(&models.Comment{})

	// Filter by article
	if query.ArticleID > 0 {
		db = db.Where("article_id = ?", query.ArticleID)
	}

	// Filter by status
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	} else {
		// Default: only show approved
		db = db.Where("status = ?", 1)
	}

	// Only get root comments
	db = db.Where("parent_id IS NULL")

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get list
	offset := (query.Page - 1) * query.PageSize
	db = db.Preload("User").Order("created_at DESC").Offset(offset).Limit(query.PageSize)

	if err := db.Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// Load replies
	for _, comment := range comments {
		var replies []*models.Comment
		if err := database.DB.Where("parent_id = ?", comment.ID).Preload("User").Order("created_at ASC").Find(&replies).Error; err == nil {
			// Convert to response format with replies
		}
	}

	return comments, total, nil
}

func (s *CommentService) UpdateStatus(id uint, status int) error {
	return database.DB.Model(&models.Comment{}).Where("id = ?", id).Update("status", status).Error
}

