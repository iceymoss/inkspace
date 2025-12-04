package service

import (
	"errors"

	"mysite/internal/database"
	"mysite/internal/models"

	"gorm.io/gorm"
)

type LikeService struct{}

func NewLikeService() *LikeService {
	return &LikeService{}
}

// LikeArticle 点赞文章
func (s *LikeService) LikeArticle(userID, articleID uint, ip string) error {
	// 检查文章是否存在
	var article models.Article
	if err := database.DB.First(&article, articleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	// 检查是否已点赞
	var count int64
	query := database.DB.Model(&models.ArticleLike{}).Where("article_id = ?", articleID)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("user_id IS NULL AND ip = ?", ip)
	}
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("已经点赞过该文章")
	}

	// 创建点赞记录
	like := &models.ArticleLike{
		ArticleID: articleID,
		IP:        ip,
	}
	// 如果 userID > 0，设置 UserID，否则为 nil（游客）
	if userID > 0 {
		like.UserID = &userID
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建点赞记录
		if err := tx.Create(like).Error; err != nil {
			return err
		}

		// 更新文章点赞数
		if err := tx.Model(&models.Article{}).
			Where("id = ?", articleID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	// 发送通知（异步，不阻塞主流程）
	if err == nil && userID > 0 {
		go func() {
			notificationService := NewNotificationService()
			_ = notificationService.CreateLikeNotification(userID, article.AuthorID, articleID)
		}()
	}

	return err
}

// UnlikeArticle 取消点赞文章
func (s *LikeService) UnlikeArticle(userID, articleID uint, ip string) error {
	// 查找点赞记录
	var like models.ArticleLike
	query := database.DB.Where("article_id = ?", articleID)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("user_id IS NULL AND ip = ?", ip)
	}

	if err := query.First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未点赞该文章")
		}
		return err
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除点赞记录
		if err := tx.Delete(&like).Error; err != nil {
			return err
		}

		// 更新文章点赞数
		if err := tx.Model(&models.Article{}).
			Where("id = ?", articleID).
			UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// IsArticleLiked 检查是否已点赞文章
func (s *LikeService) IsArticleLiked(userID, articleID uint, ip string) (bool, error) {
	var count int64
	query := database.DB.Model(&models.ArticleLike{}).Where("article_id = ?", articleID)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("user_id IS NULL AND ip = ?", ip)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// LikeComment 点赞评论
func (s *LikeService) LikeComment(userID, commentID uint, ip string) error {
	// 检查评论是否存在
	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}

	// 检查是否已点赞
	var count int64
	query := database.DB.Model(&models.CommentLike{}).Where("comment_id = ?", commentID)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("user_id IS NULL AND ip = ?", ip)
	}
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("已经点赞过该评论")
	}

	// 创建点赞记录
	like := &models.CommentLike{
		CommentID: commentID,
		IP:        ip,
	}
	// 如果 userID > 0，设置 UserID，否则为 nil（游客）
	if userID > 0 {
		like.UserID = &userID
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建点赞记录
		if err := tx.Create(like).Error; err != nil {
			return err
		}

		// 更新评论点赞数
		if err := tx.Model(&models.Comment{}).
			Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	// 发送通知（异步，不阻塞主流程）
	// 注意：只有登录用户点赞时才发送通知，且评论作者必须存在且不是自己
	if err == nil && userID > 0 && comment.UserID > 0 && userID != comment.UserID {
		go func() {
			notificationService := NewNotificationService()
			_ = notificationService.CreateCommentLikeNotification(userID, comment.UserID, comment.ArticleID, commentID)
		}()
	}

	return err
}

// UnlikeComment 取消点赞评论
func (s *LikeService) UnlikeComment(userID, commentID uint, ip string) error {
	// 查找点赞记录
	var like models.CommentLike
	query := database.DB.Where("comment_id = ?", commentID)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("user_id IS NULL AND ip = ?", ip)
	}

	if err := query.First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未点赞该评论")
		}
		return err
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除点赞记录
		if err := tx.Delete(&like).Error; err != nil {
			return err
		}

		// 更新评论点赞数
		if err := tx.Model(&models.Comment{}).
			Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// IsCommentLiked 检查是否已点赞评论
func (s *LikeService) IsCommentLiked(userID, commentID uint, ip string) (bool, error) {
	var count int64
	query := database.DB.Model(&models.CommentLike{}).Where("comment_id = ?", commentID)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("user_id IS NULL AND ip = ?", ip)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

