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

	// Check if parent comment exists and get root_id
	var rootID *uint
	if req.ParentID != nil {
		var parent models.Comment
		if err := database.DB.First(&parent, *req.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父评论不存在")
			}
			return nil, err
		}
		// 如果父评论有root_id，则使用父评论的root_id，否则使用父评论的id
		if parent.RootID != nil {
			rootID = parent.RootID
		} else {
			rootID = &parent.ID
		}
	}

	comment := &models.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID,
		Content:   req.Content,
		ParentID:  req.ParentID,
		RootID:    rootID,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Website:   req.Website,
		Status:    1, // Auto approve
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建评论
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		// 更新文章评论数
		if err := tx.Model(&models.Article{}).
			Where("id = ?", req.ArticleID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return err
		}

		// 更新用户评论数
		if userID > 0 {
			if err := tx.Model(&models.User{}).
				Where("id = ?", userID).
				UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
				return err
			}
		}

		// 更新父评论回复数
		if req.ParentID != nil {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", *req.ParentID).
				UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 发送通知（异步，不阻塞主流程）
	go func() {
		notificationService := NewNotificationService()
		
		if req.ParentID != nil {
			// 回复评论：需要通知两个对象
			var parentComment models.Comment
			if err := database.DB.First(&parentComment, *req.ParentID).Error; err == nil {
				// 1. 通知被回复的评论作者
				if parentComment.UserID > 0 && userID > 0 && parentComment.UserID != userID {
					_ = notificationService.CreateReplyNotification(
						userID,
						parentComment.UserID,
						req.ArticleID,
						*req.ParentID,
						req.Content,
					)
				}
				
				// 2. 通知文章作者（如果文章作者不是回复者本人，也不是被回复的评论作者）
				if article.AuthorID > 0 && userID > 0 && 
				   article.AuthorID != userID && 
				   article.AuthorID != parentComment.UserID {
					_ = notificationService.CreateCommentNotification(
						userID,
						article.AuthorID,
						req.ArticleID,
						req.Content,
					)
				}
			}
		} else {
			// 评论文章：通知文章作者
			if article.AuthorID > 0 && userID > 0 && article.AuthorID != userID {
				_ = notificationService.CreateCommentNotification(
					userID,
					article.AuthorID,
					req.ArticleID,
					req.Content,
				)
			}
		}
	}()

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

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除评论
		if err := tx.Delete(&comment).Error; err != nil {
			return err
		}

		// 更新文章评论数
		if err := tx.Model(&models.Article{}).
			Where("id = ?", comment.ArticleID).
			UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
			return err
		}

		// 更新用户评论数
		if comment.UserID > 0 {
			if err := tx.Model(&models.User{}).
				Where("id = ?", comment.UserID).
				UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
				return err
			}
		}

		// 更新父评论回复数
		if comment.ParentID != nil {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", *comment.ParentID).
				UpdateColumn("reply_count", gorm.Expr("reply_count - ?", 1)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (s *CommentService) GetList(query *models.CommentListQuery) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	db := database.DB.Model(&models.Comment{})

	// Filter by article
	if query.ArticleID > 0 {
		db = db.Where("article_id = ?", query.ArticleID)
	}

	// Filter by user
	if query.UserID > 0 {
		db = db.Where("user_id = ?", query.UserID)
	}

	// Filter by status
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	} else {
		// Default: only show approved
		// But if user_id is specified (viewing own comments), show all statuses
		if query.UserID == 0 {
			db = db.Where("status = ?", 1)
		}
		// If user_id > 0, don't filter by status (show all: pending, approved, rejected)
	}

	// Only get root comments (when filtering by article, not when filtering by user)
	if query.ArticleID > 0 {
		db = db.Where("parent_id IS NULL")
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get list
	offset := (query.Page - 1) * query.PageSize
	db = db.Preload("User").Preload("Article").Order("created_at DESC").Offset(offset).Limit(query.PageSize)

	if err := db.Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// Load replies for each comment and attach to comment
	for _, comment := range comments {
		var replies []*models.Comment
		if err := database.DB.Where("parent_id = ?", comment.ID).
			Preload("User").
			Preload("Article").
			Order("created_at ASC").
			Find(&replies).Error; err == nil {
			// 将回复添加到评论（通过反射或手动设置）
			// 由于 Comment 模型没有 Replies 字段，我们将在 handler 中处理
		}
	}

	return comments, total, nil
}

func (s *CommentService) UpdateStatus(id uint, status int) error {
	return database.DB.Model(&models.Comment{}).Where("id = ?", id).Update("status", status).Error
}

