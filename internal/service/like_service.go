package service

import (
	"mysite/internal/database"
	"mysite/internal/models"

	"gorm.io/gorm"
)

type LikeService struct {
	notificationService *NotificationService
}

func NewLikeService() *LikeService {
	return &LikeService{
		notificationService: NewNotificationService(),
	}
}

// LikeWork 点赞作品
func (s *LikeService) LikeWork(userID, workID uint) error {
	// 检查是否已点赞
	var existingLike models.Like
	err := database.DB.Where("user_id = ? AND work_id = ?", userID, workID).First(&existingLike).Error
	if err == nil {
		// 已点赞，取消点赞
		if err := database.DB.Delete(&existingLike).Error; err != nil {
			return err
		}

		// 减少作品点赞数
		return database.DB.Model(&models.Work{}).Where("id = ?", workID).
			UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
	}

	// 未点赞，添加点赞
	like := &models.Like{
		UserID: userID,
		WorkID: &workID,
	}
	if err := database.DB.Create(like).Error; err != nil {
		return err
	}

	// 增加作品点赞数
	if err := database.DB.Model(&models.Work{}).Where("id = ?", workID).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		return err
	}

	// 发送通知给作品作者
	var work models.Work
	if err := database.DB.First(&work, workID).Error; err == nil {
		if work.AuthorID != userID {
			go s.notificationService.CreateLikeNotification(userID, work.AuthorID, nil, &workID)
		}
	}

	return nil
}

// LikeArticle 点赞文章
func (s *LikeService) LikeArticle(userID, articleID uint) error {
	// 检查是否已点赞
	var existingLike models.Like
	err := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&existingLike).Error
	if err == nil {
		// 已点赞，取消点赞
		if err := database.DB.Delete(&existingLike).Error; err != nil {
			return err
		}

		// 减少文章点赞数
		return database.DB.Model(&models.Article{}).Where("id = ?", articleID).
			UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
	}

	// 未点赞，添加点赞
	like := &models.Like{
		UserID:    userID,
		ArticleID: &articleID,
	}
	if err := database.DB.Create(like).Error; err != nil {
		return err
	}

	// 增加文章点赞数
	if err := database.DB.Model(&models.Article{}).Where("id = ?", articleID).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		return err
	}

	// 发送通知给文章作者
	var article models.Article
	if err := database.DB.First(&article, articleID).Error; err == nil {
		if article.AuthorID != userID {
			go s.notificationService.CreateLikeNotification(userID, article.AuthorID, &articleID, nil)
		}
	}

	return nil
}

// CheckWorkLiked 检查用户是否已点赞作品
func (s *LikeService) CheckWorkLiked(userID, workID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Like{}).
		Where("user_id = ? AND work_id = ?", userID, workID).
		Count(&count).Error

	return count > 0, err
}

// CheckArticleLiked 检查用户是否已点赞文章
func (s *LikeService) CheckArticleLiked(userID, articleID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Like{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error

	return count > 0, err
}

// GetUserLikes 获取用户的点赞列表
func (s *LikeService) GetUserLikes(userID uint, page, pageSize int) ([]*models.Like, int64, error) {
	var likes []*models.Like
	var total int64

	db := database.DB.Model(&models.Like{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Preload("Article").Preload("Work").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&likes).Error

	return likes, total, err
}
