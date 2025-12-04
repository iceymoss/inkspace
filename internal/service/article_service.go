package service

import (
	"errors"
	"fmt"
	"time"

	"mysite/internal/config"
	"mysite/internal/database"
	"mysite/internal/models"
	"mysite/internal/utils"

	"gorm.io/gorm"
)

type ArticleService struct{}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

func (s *ArticleService) Create(req *models.ArticleRequest, authorID uint) (*models.Article, error) {
	article := &models.Article{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		Cover:       req.Cover,
		CategoryID:  req.CategoryID,
		AuthorID:    authorID,
		Status:      req.Status,
		IsTop:       req.IsTop,
		IsRecommend: req.IsRecommend,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Create article
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		// Associate tags
		if len(req.TagIDs) > 0 {
			var tags []models.Tag
			if err := tx.Where("id IN ?", req.TagIDs).Find(&tags).Error; err != nil {
				return err
			}
			if err := tx.Model(article).Association("Tags").Append(tags); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Clear cache
	utils.DeleteCachePattern("article:*")

	return article, nil
}

func (s *ArticleService) Update(id uint, req *models.ArticleRequest, userID uint, role string) (*models.Article, error) {
	article, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check permission
	if role != "admin" && article.AuthorID != userID {
		return nil, errors.New("无权限修改")
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Update article
		article.Title = req.Title
		article.Content = req.Content
		article.Summary = req.Summary
		article.Cover = req.Cover
		article.CategoryID = req.CategoryID
		article.Status = req.Status
		article.IsTop = req.IsTop
		article.IsRecommend = req.IsRecommend

		if err := tx.Save(article).Error; err != nil {
			return err
		}

		// Update tags
		if err := tx.Model(article).Association("Tags").Clear(); err != nil {
			return err
		}
		if len(req.TagIDs) > 0 {
			var tags []models.Tag
			if err := tx.Where("id IN ?", req.TagIDs).Find(&tags).Error; err != nil {
				return err
			}
			if err := tx.Model(article).Association("Tags").Append(tags); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Clear cache
	utils.DeleteCache(fmt.Sprintf("article:%d", id))
	utils.DeleteCachePattern("article:list:*")

	return article, nil
}

func (s *ArticleService) Delete(id uint, userID uint, role string) error {
	article, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// Check permission
	if role != "admin" && article.AuthorID != userID {
		return errors.New("无权限删除")
	}

	if err := database.DB.Delete(article).Error; err != nil {
		return err
	}

	// Clear cache
	utils.DeleteCache(fmt.Sprintf("article:%d", id))
	utils.DeleteCachePattern("article:list:*")

	return nil
}

func (s *ArticleService) GetByID(id uint) (*models.Article, error) {
	var article models.Article

	// Try to get from cache
	cacheKey := fmt.Sprintf("article:%d", id)
	if err := utils.GetCache(cacheKey, &article); err == nil {
		return &article, nil
	}

	// Get from database
	if err := database.DB.Preload("Category").Preload("Tags").Preload("Author").First(&article, id).Error; err != nil {
		return nil, err
	}

	// Save to cache
	expiration := time.Duration(config.AppConfig.Cache.ArticleExpire) * time.Second
	utils.SetCache(cacheKey, &article, expiration)

	return &article, nil
}

func (s *ArticleService) GetList(query *models.ArticleListQuery) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	db := database.DB.Model(&models.Article{})

	// Filter by status
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	} else {
		// Default: only show published
		db = db.Where("status = ?", 1)
	}

	// Filter by category
	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}

	// Filter by tag
	if query.TagID > 0 {
		db = db.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", query.TagID)
	}

	// Search by keyword
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", keyword, keyword)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get list
	offset := (query.Page - 1) * query.PageSize
	db = db.Preload("Category").Preload("Tags").Preload("Author").
		Order("is_top DESC, created_at DESC").
		Offset(offset).Limit(query.PageSize)

	if err := db.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (s *ArticleService) IncrementViewCount(id uint) error {
	return database.DB.Model(&models.Article{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (s *ArticleService) IncrementLikeCount(id uint) error {
	if err := database.DB.Model(&models.Article{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		return err
	}

	// Clear cache
	utils.DeleteCache(fmt.Sprintf("article:%d", id))

	return nil
}

