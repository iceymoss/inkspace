package service

import (
	"errors"
	"fmt"
	"log"

	"mysite/internal/database"
	"mysite/internal/models"
	"mysite/internal/utils"

	"gorm.io/gorm"
)

type FavoriteService struct {
	notificationService *NotificationService
}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{
		notificationService: NewNotificationService(),
	}
}

// AddFavorite 收藏文章
func (s *FavoriteService) AddFavorite(userID, articleID uint) error {
	// 检查文章是否存在
	var article models.Article
	if err := database.DB.First(&article, articleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	// 检查是否已收藏
	var count int64
	if err := database.DB.Model(&models.ArticleFavorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("已经收藏过该文章")
	}

	// 创建收藏记录
	favorite := &models.ArticleFavorite{
		UserID:    userID,
		ArticleID: articleID,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建收藏记录
		if err := tx.Create(favorite).Error; err != nil {
			return err
		}

		// 更新文章收藏数
		if err := tx.Model(&models.Article{}).
			Where("id = ?", articleID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}

		// 更新用户收藏数
		if err := tx.Model(&models.User{}).
			Where("id = ?", userID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}

		// 发送收藏通知给文章作者
		if article.AuthorID != userID {
			go s.notificationService.CreateFavoriteNotification(userID, article.AuthorID, &articleID, nil)
		}

		return nil
	})

	if err == nil {
		// 清除文章缓存
		utils.DeleteCache(fmt.Sprintf("article:%d", articleID))
	}

	return err
}

// RemoveFavorite 取消收藏
func (s *FavoriteService) RemoveFavorite(userID, articleID uint) error {
	// 检查是否已收藏
	var favorite models.ArticleFavorite
	if err := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).
		First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未收藏该文章")
		}
		return err
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除收藏记录
		if err := tx.Delete(&favorite).Error; err != nil {
			return err
		}

		// 更新文章收藏数
		if err := tx.Model(&models.Article{}).
			Where("id = ?", articleID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			return err
		}

		// 更新用户收藏数
		if err := tx.Model(&models.User{}).
			Where("id = ?", userID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// IsFavorited 检查是否已收藏
func (s *FavoriteService) IsFavorited(userID, articleID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.ArticleFavorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error
	return count > 0, err
}

// GetFavoriteList 获取用户收藏列表
func (s *FavoriteService) GetFavoriteList(userID uint, page, pageSize int) ([]*models.FavoriteResponse, int64, error) {
	var favorites []*models.ArticleFavorite
	var total int64

	// 查询收藏列表
	db := database.DB.Model(&models.ArticleFavorite{}).
		Where("user_id = ?", userID).
		Preload("Article").
		Preload("Article.Category").
		Preload("Article.Tags").
		Preload("Article.Author")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&favorites).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	responses := make([]*models.FavoriteResponse, len(favorites))
	for i, favorite := range favorites {
		responses[i] = favorite.ToResponse()
	}

	return responses, total, nil
}

// GetFavoriteCount 获取文章的收藏数
func (s *FavoriteService) GetFavoriteCount(articleID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.ArticleFavorite{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}

// AddWorkFavorite 收藏作品
func (s *FavoriteService) AddWorkFavorite(userID, workID uint) error {
	// 检查是否已收藏
	var count int64
	if err := database.DB.Model(&models.Favorite{}).
		Where("user_id = ? AND work_id = ?", userID, workID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("已经收藏过该作品")
	}

	// 创建收藏记录
	favorite := &models.Favorite{
		UserID: userID,
		WorkID: &workID,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建收藏记录
		if err := tx.Create(favorite).Error; err != nil {
			return err
		}

		// 增加作品收藏数
		if err := tx.Model(&models.Work{}).
			Where("id = ?", workID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}

		// 增加用户收藏数
		if err := tx.Model(&models.User{}).
			Where("id = ?", userID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 清除作品缓存
	utils.DeleteCache(fmt.Sprintf("work:%d", workID))

	// 发送收藏通知给作品作者
	var work models.Work
	if err := database.DB.First(&work, workID).Error; err != nil {
		log.Printf("❌ 获取作品信息失败 (ID: %d): %v", workID, err)
		return nil
	}

	if work.AuthorID != userID {
		go func() {
			err := s.notificationService.CreateFavoriteNotification(userID, work.AuthorID, nil, &workID)
			if err != nil {
				log.Printf("❌ 创建作品收藏通知失败: 用户%d -> 用户%d, 作品%d, 错误: %v", userID, work.AuthorID, workID, err)
			} else {
				log.Printf("✅ 成功创建作品收藏通知: 用户%d -> 用户%d, 作品%d", userID, work.AuthorID, workID)
			}
		}()
	} else {
		log.Printf("ℹ️ 用户收藏自己的作品，不发送通知 (用户ID: %d, 作品ID: %d)", userID, workID)
	}

	return nil
}

// RemoveWorkFavorite 取消收藏作品
func (s *FavoriteService) RemoveWorkFavorite(userID, workID uint) error {
	var favorite models.Favorite
	if err := database.DB.Where("user_id = ? AND work_id = ?", userID, workID).
		First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未收藏该作品")
		}
		return err
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除收藏记录
		if err := tx.Delete(&favorite).Error; err != nil {
			return err
		}

		// 减少作品收藏数
		if err := tx.Model(&models.Work{}).
			Where("id = ?", workID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			return err
		}

		// 减少用户收藏数
		if err := tx.Model(&models.User{}).
			Where("id = ?", userID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	if err == nil {
		// 清除作品缓存
		utils.DeleteCache(fmt.Sprintf("work:%d", workID))
	}

	return err
}

// CheckWorkFavorited 检查是否已收藏作品
func (s *FavoriteService) CheckWorkFavorited(userID, workID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Favorite{}).
		Where("user_id = ? AND work_id = ?", userID, workID).
		Count(&count).Error

	return count > 0, err
}
