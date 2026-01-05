package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"

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
		database.DeleteCache(fmt.Sprintf("article:%d", articleID))
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

// GetFavoriteList 获取用户收藏列表（包含文章和作品）
// 数据安全：只返回用户自己的收藏，并且只返回公开的内容
func (s *FavoriteService) GetFavoriteList(userID uint, page, pageSize int) ([]*models.FavoriteResponse, int64, error) {
	var allFavorites []*models.FavoriteResponse

	// 1. 获取文章收藏
	articleDB := database.DB.Model(&models.ArticleFavorite{}).
		Joins("JOIN articles ON articles.id = article_favorites.article_id").
		Where("article_favorites.user_id = ?", userID).
		Where("(articles.status = ? OR (articles.status = ? AND articles.author_id = ?))", 1, 2, userID)

	var articleTotal int64
	if err := articleDB.Count(&articleTotal).Error; err != nil {
		return nil, 0, err
	}

	// 获取所有文章收藏（不分页，用于合并排序）
	var articleFavorites []*models.ArticleFavorite
	if err := articleDB.Select("article_favorites.*").
		Order("article_favorites.created_at DESC").
		Find(&articleFavorites).Error; err != nil {
		return nil, 0, err
	}

	// 预加载文章信息
	if len(articleFavorites) > 0 {
		articleIDs := make([]uint, len(articleFavorites))
		for i, f := range articleFavorites {
			articleIDs[i] = f.ArticleID
		}

		var articles []*models.Article
		if err := database.DB.Where("id IN ?", articleIDs).
			Where("(status = ? OR (status = ? AND author_id = ?))", 1, 2, userID).
			Preload("Category").Preload("Tags").Preload("Author").
			Find(&articles).Error; err != nil {
			return nil, 0, err
		}

		articleMap := make(map[uint]*models.Article)
		for _, article := range articles {
			articleMap[article.ID] = article
		}
		for i := range articleFavorites {
			if article, ok := articleMap[articleFavorites[i].ArticleID]; ok {
				articleFavorites[i].Article = article
			}
		}

		// 转换为响应格式，只添加有文章数据的收藏项
		for _, favorite := range articleFavorites {
			// 只处理有文章数据的收藏项
			if favorite.Article != nil {
				resp := favorite.ToResponse()
				resp.Type = "article"
				allFavorites = append(allFavorites, resp)
			}
		}
	}

	// 2. 获取作品收藏（只获取已发布的作品）
	workDB := database.DB.Model(&models.Favorite{}).
		Joins("JOIN works ON works.id = favorites.work_id").
		Where("favorites.user_id = ? AND favorites.work_id IS NOT NULL", userID).
		Where("works.status = ?", 1)

	var workTotal int64
	if err := workDB.Count(&workTotal).Error; err != nil {
		return nil, 0, err
	}

	// 获取所有作品收藏（不分页，用于合并排序）
	var workFavorites []*models.Favorite
	if err := workDB.Order("created_at DESC").
		Find(&workFavorites).Error; err != nil {
		return nil, 0, err
	}

	// 预加载作品信息
	if len(workFavorites) > 0 {
		workIDs := make([]uint, 0, len(workFavorites))
		for _, f := range workFavorites {
			if f.WorkID != nil {
				workIDs = append(workIDs, *f.WorkID)
			}
		}

		if len(workIDs) > 0 {
			var works []*models.Work
			if err := database.DB.Where("id IN ? AND status = ?", workIDs, 1).
				Preload("Author").
				Find(&works).Error; err != nil {
				return nil, 0, err
			}

			workMap := make(map[uint]*models.Work)
			for _, work := range works {
				workMap[work.ID] = work
			}
			for i := range workFavorites {
				if workFavorites[i].WorkID != nil {
					if work, ok := workMap[*workFavorites[i].WorkID]; ok {
						workFavorites[i].Work = work
					}
				}
			}

			// 转换为响应格式
			for _, favorite := range workFavorites {
				if favorite.WorkID != nil && favorite.Work != nil {
					resp := &models.FavoriteResponse{
						ID:        favorite.ID,
						UserID:    favorite.UserID,
						WorkID:    *favorite.WorkID,
						Type:      "work",
						CreatedAt: favorite.CreatedAt,
					}
					// 转换为WorkResponse
					workResp := favorite.Work.ToResponse()
					resp.Work = workResp
					allFavorites = append(allFavorites, resp)
				}
			}
		}
	}

	// 3. 按创建时间降序排序（使用简单的冒泡排序，数据量不大时性能可接受）
	for i := 0; i < len(allFavorites)-1; i++ {
		for j := i + 1; j < len(allFavorites); j++ {
			if allFavorites[i].CreatedAt.Before(allFavorites[j].CreatedAt) {
				allFavorites[i], allFavorites[j] = allFavorites[j], allFavorites[i]
			}
		}
	}

	// 4. 应用分页
	// 使用实际过滤后的数量作为总数，确保总数和实际返回的数据一致
	total := int64(len(allFavorites))
	offset := (page - 1) * pageSize
	start := offset
	end := offset + pageSize

	if start >= len(allFavorites) {
		allFavorites = []*models.FavoriteResponse{}
	} else if end > len(allFavorites) {
		allFavorites = allFavorites[start:]
	} else {
		allFavorites = allFavorites[start:end]
	}

	return allFavorites, total, nil
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
	// 检查作品是否存在并获取状态
	var work models.Work
	if err := database.DB.First(&work, workID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("作品不存在")
		}
		return err
	}

	// 只有已发布（status=1）的作品才能被收藏
	if work.Status != 1 {
		return errors.New("该作品尚未发布，无法收藏")
	}

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
	database.DeleteCache(fmt.Sprintf("work:%d", workID))

	// 发送收藏通知给作品作者
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
	// 检查作品是否存在并获取状态
	var work models.Work
	if err := database.DB.First(&work, workID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("作品不存在")
		}
		return err
	}

	// 只有已发布（status=1）的作品才能取消收藏
	// 这样可以防止对未发布作品进行任何收藏相关操作
	if work.Status != 1 {
		return errors.New("该作品尚未发布，无法操作")
	}

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
		database.DeleteCache(fmt.Sprintf("work:%d", workID))
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
