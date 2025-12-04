package service

import (
	"errors"

	"mysite/internal/database"
	"mysite/internal/models"

	"gorm.io/gorm"
)

type FavoriteService struct{}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{}
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

		// TODO: 发送收藏通知给文章作者
		// if article.AuthorID != userID {
		//     notificationService.CreateFavoriteNotification(userID, article.AuthorID, articleID)
		// }

		return nil
	})

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

