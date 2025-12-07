package service

import (
	"errors"
	"log"

	"mysite/internal/database"
	"mysite/internal/models"

	"gorm.io/gorm"
)

type CommentService struct {
	notificationService *NotificationService
}

func NewCommentService() *CommentService {
	return &CommentService{
		notificationService: NewNotificationService(),
	}
}

func (s *CommentService) Create(req *models.CommentRequest, userID uint) (*models.Comment, error) {
	// 验证：必须指定文章ID或作品ID其中之一
	if req.ArticleID == nil && req.WorkID == nil {
		return nil, errors.New("必须指定文章或作品")
	}
	if req.ArticleID != nil && req.WorkID != nil {
		return nil, errors.New("不能同时评论文章和作品")
	}

	// Check if article/work exists and load for notifications
	var article *models.Article
	if req.ArticleID != nil && *req.ArticleID > 0 {
		article = &models.Article{}
		if err := database.DB.First(article, *req.ArticleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("文章不存在")
			}
			return nil, err
		}
	}

	if req.WorkID != nil && *req.WorkID > 0 {
		var work models.Work
		if err := database.DB.First(&work, *req.WorkID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("作品不存在")
			}
			return nil, err
		}
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

	// 处理 article_id 和 work_id：确保只有一个有值，另一个为 nil
	var articleID *uint
	var workID *uint

	if req.ArticleID != nil && *req.ArticleID > 0 {
		articleID = req.ArticleID
	}

	if req.WorkID != nil && *req.WorkID > 0 {
		workID = req.WorkID
	}

	comment := &models.Comment{
		ArticleID: articleID,
		WorkID:    workID,
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
		if req.ArticleID != nil && *req.ArticleID > 0 {
			if err := tx.Model(&models.Article{}).
				Where("id = ?", *req.ArticleID).
				UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
				return err
			}
		}

		// 更新作品评论数
		if req.WorkID != nil && *req.WorkID > 0 {
			if err := tx.Model(&models.Work{}).
				Where("id = ?", *req.WorkID).
				UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
				return err
			}
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
	if req.ArticleID != nil && *req.ArticleID > 0 && article != nil {
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
							nil,
							comment.ID,
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
							nil,
							comment.ID,
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
						nil,
						comment.ID,
					)
				}
			}
		}()
	} else if workID != nil && *workID > 0 {
		// 作品评论通知
		go func() {
			notificationService := NewNotificationService()
			var work models.Work
			if err := database.DB.First(&work, *workID).Error; err != nil {
				log.Printf("❌ 获取作品信息失败 (ID: %d): %v", *workID, err)
				return
			}

			if req.ParentID != nil {
				// 回复评论：只通知被回复的评论作者
				var parentComment models.Comment
				if err := database.DB.First(&parentComment, *req.ParentID).Error; err == nil {
					// 通知被回复的评论作者（回复通知）
					if parentComment.UserID > 0 && userID > 0 && parentComment.UserID != userID {
						_ = notificationService.CreateReplyNotification(
							userID,
							parentComment.UserID,
							nil,
							workID,
							comment.ID,
						)
					}
					// 注意：回复评论时，不通知作品作者（除非作品作者就是被回复的评论作者，上面已经通知了）
				}
			} else {
				// 评论作品：通知作品作者
				if work.AuthorID > 0 && userID > 0 && work.AuthorID != userID {
					_ = notificationService.CreateCommentNotification(
						userID,
						work.AuthorID,
						nil,
						workID,
						comment.ID,
					)
				}
			}
		}()
	}

	return comment, nil
}

func (s *CommentService) Delete(id uint, userID uint, role string) error {
	// 先查询评论以获取相关信息
	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}

	// Check permission: 评论作者、文章/作品作者或管理员可以删除
	if role != "admin" {
		isCommentAuthor := comment.UserID > 0 && comment.UserID == userID
		isContentAuthor := false

		// 检查是否是文章作者（使用WHERE条件确保权限）
		if comment.ArticleID != nil && *comment.ArticleID > 0 {
			var count int64
			if err := database.DB.Model(&models.Article{}).
				Where("id = ? AND author_id = ?", *comment.ArticleID, userID).
				Count(&count).Error; err == nil && count > 0 {
				isContentAuthor = true
			}
		}

		// 检查是否是作品作者（使用WHERE条件确保权限）
		if !isContentAuthor && comment.WorkID != nil && *comment.WorkID > 0 {
			var count int64
			if err := database.DB.Model(&models.Work{}).
				Where("id = ? AND author_id = ?", *comment.WorkID, userID).
				Count(&count).Error; err == nil && count > 0 {
				isContentAuthor = true
			}
		}

		// 如果既不是评论作者，也不是内容作者，则无权限
		if !isCommentAuthor && !isContentAuthor {
			if comment.UserID == 0 {
				// 游客评论，只有内容作者或管理员可以删除
				return errors.New("无权限删除游客评论")
			}
			return errors.New("只能删除自己发表的评论或自己内容下的评论")
		}
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除评论 - 使用WHERE条件确保权限
		// 管理员可以删除任何评论，普通用户只能删除自己的评论或自己内容下的评论
		deleteQuery := tx.Where("id = ?", id)
		if role != "admin" {
			// 非管理员：只能删除自己的评论，或者自己内容下的评论
			// 这里已经在上面检查过权限，所以直接删除
			// 但为了安全，我们再次确认：如果是评论作者，使用user_id条件
			if comment.UserID > 0 && comment.UserID == userID {
				// 评论作者删除自己的评论
				deleteQuery = deleteQuery.Where("user_id = ?", userID)
			} else {
				// 内容作者删除自己内容下的评论（已经在上面验证过）
				// 这里直接删除，因为权限已经验证
			}
		}
		
		if err := deleteQuery.Delete(&models.Comment{}).Error; err != nil {
			return err
		}

		// 更新文章评论数
		if comment.ArticleID != nil && *comment.ArticleID > 0 {
			if err := tx.Model(&models.Article{}).
				Where("id = ?", *comment.ArticleID).
				UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
				return err
			}
		}

		// 更新作品评论数
		if comment.WorkID != nil && *comment.WorkID > 0 {
			if err := tx.Model(&models.Work{}).
				Where("id = ?", *comment.WorkID).
				UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
				return err
			}
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
	if query.ArticleID != nil && *query.ArticleID > 0 {
		db = db.Where("article_id = ?", *query.ArticleID)
	}

	// Filter by work
	if query.WorkID != nil && *query.WorkID > 0 {
		db = db.Where("work_id = ?", *query.WorkID)
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

	// Only get root comments (when filtering by article/work, not when filtering by user)
	if (query.ArticleID != nil && *query.ArticleID > 0) || (query.WorkID != nil && *query.WorkID > 0) {
		db = db.Where("parent_id IS NULL")
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get list
	offset := (query.Page - 1) * query.PageSize
	db = db.Preload("User").Preload("Article").Preload("Work").Order("created_at DESC").Offset(offset).Limit(query.PageSize)

	if err := db.Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// Load replies for each comment and attach to comment
	for _, comment := range comments {
		var replies []*models.Comment
		if err := database.DB.Where("parent_id = ?", comment.ID).
			Preload("User").
			Preload("Article").
			Preload("Work").
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
