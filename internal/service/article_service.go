package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
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

		// 更新用户文章数
		if err := tx.Model(&models.User{}).
			Where("id = ?", authorID).
			UpdateColumn("article_count", gorm.Expr("article_count + ?", 1)).Error; err != nil {
			return err
		}

		// 更新分类文章数
		if req.CategoryID > 0 {
			if err := tx.Model(&models.Category{}).
				Where("id = ?", req.CategoryID).
				UpdateColumn("article_count", gorm.Expr("article_count + ?", 1)).Error; err != nil {
				return err
			}
		}

		// 更新标签文章数
		if len(req.TagIDs) > 0 {
			if err := tx.Model(&models.Tag{}).
				Where("id IN ?", req.TagIDs).
				UpdateColumn("article_count", gorm.Expr("article_count + ?", 1)).Error; err != nil {
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

	// 获取旧的标签IDs和分类ID
	oldTagIDs := make([]uint, len(article.Tags))
	for i, tag := range article.Tags {
		oldTagIDs[i] = tag.ID
	}
	oldCategoryID := article.CategoryID

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 如果分类改变，更新旧分类的文章数
		if oldCategoryID > 0 && oldCategoryID != req.CategoryID {
			if err := tx.Model(&models.Category{}).
				Where("id = ?", oldCategoryID).
				UpdateColumn("article_count", gorm.Expr("article_count - ?", 1)).Error; err != nil {
				return err
			}

			// 更新新分类的文章数
			if req.CategoryID > 0 {
				if err := tx.Model(&models.Category{}).
					Where("id = ?", req.CategoryID).
					UpdateColumn("article_count", gorm.Expr("article_count + ?", 1)).Error; err != nil {
					return err
				}
			}
		}

		// 更新旧标签的文章数（减少）
		if len(oldTagIDs) > 0 {
			if err := tx.Model(&models.Tag{}).
				Where("id IN ?", oldTagIDs).
				UpdateColumn("article_count", gorm.Expr("article_count - ?", 1)).Error; err != nil {
				return err
			}
		}

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

			// 更新新标签的文章数（增加）
			if err := tx.Model(&models.Tag{}).
				Where("id IN ?", req.TagIDs).
				UpdateColumn("article_count", gorm.Expr("article_count + ?", 1)).Error; err != nil {
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

	// 获取标签IDs
	tagIDs := make([]uint, len(article.Tags))
	for i, tag := range article.Tags {
		tagIDs[i] = tag.ID
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除文章
		if err := tx.Delete(article).Error; err != nil {
			return err
		}

		// 更新用户文章数
		if err := tx.Model(&models.User{}).
			Where("id = ?", article.AuthorID).
			UpdateColumn("article_count", gorm.Expr("article_count - ?", 1)).Error; err != nil {
			return err
		}

		// 更新分类文章数
		if article.CategoryID > 0 {
			if err := tx.Model(&models.Category{}).
				Where("id = ?", article.CategoryID).
				UpdateColumn("article_count", gorm.Expr("article_count - ?", 1)).Error; err != nil {
				return err
			}
		}

		// 更新标签文章数
		if len(tagIDs) > 0 {
			if err := tx.Model(&models.Tag{}).
				Where("id IN ?", tagIDs).
				UpdateColumn("article_count", gorm.Expr("article_count - ?", 1)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
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
	// 如果指定了榜单类型（hot, week, month, year），从Redis ZSet获取
	// 但如果同时指定了排序方式（非热门排序），则使用数据库查询并应用排序
	if query.RankType != "" && query.RankType != "none" && (query.SortBy == "" || query.SortBy == "hot") {
		return s.getListFromRank(query)
	}

	var articles []*models.Article
	var total int64

	db := database.DB.Model(&models.Article{})

	// Filter by status
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	} else {
		// Default: only show published
		// But if author_id is specified (viewing own articles), show all statuses
		if query.AuthorID == 0 {
			db = db.Where("status = ?", 1)
		}
		// If author_id > 0, don't filter by status (show all: draft and published)
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

	// Filter by author
	if query.AuthorID > 0 {
		db = db.Where("author_id = ?", query.AuthorID)
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
	db = db.Preload("Category").Preload("Tags").Preload("Author")

	// 排序逻辑
	orderBy := s.buildOrderBy(query)
	db = db.Order(orderBy).Offset(offset).Limit(query.PageSize)

	if err := db.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// buildOrderBy 构建排序SQL
func (s *ArticleService) buildOrderBy(query *models.ArticleListQuery) string {
	// 默认排序：置顶优先，然后按创建时间
	if query.SortBy == "" {
		return "is_top DESC, created_at DESC"
	}

	// 置顶优先
	order := "is_top DESC, "

	// 排序字段
	sortField := "created_at"
	switch query.SortBy {
	case "time":
		sortField = "created_at"
	case "view_count":
		sortField = "view_count"
	case "like_count":
		sortField = "like_count"
	case "comment_count":
		sortField = "comment_count"
	case "hot":
		// 热门排序需要特殊处理，使用综合得分
		// 这里简化处理，使用 view_count + like_count + comment_count 作为热门度
		sortField = "(view_count * 0.5 + like_count * 0.15 + comment_count * 0.2 + favorite_count * 0.15)"
	}

	// 排序方向
	sortOrder := "DESC"
	if query.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	return order + sortField + " " + sortOrder
}

// getListFromRank 从榜单（Redis ZSet）获取文章列表
func (s *ArticleService) getListFromRank(query *models.ArticleListQuery) ([]*models.Article, int64, error) {
	ctx := database.Ctx
	var zsetKey string

	// 根据榜单类型确定Redis Key
	switch query.RankType {
	case "hot":
		zsetKey = "hot:articles:zset"
	case "week":
		zsetKey = "rank:articles:week"
	case "month":
		zsetKey = "rank:articles:month"
	case "year":
		zsetKey = "rank:articles:year"
	default:
		// 默认使用热门
		zsetKey = "hot:articles:zset"
	}

	// 从Redis ZSet获取文章ID（按分值降序）
	// 如果使用热门排序，需要分页；否则获取所有文章ID以便在数据库层面排序
	var articleIDStrs []string
	var err error
	if query.SortBy == "" || query.SortBy == "hot" {
		// 热门排序：使用分页
		start := int64((query.Page - 1) * query.PageSize)
		end := start + int64(query.PageSize) - 1
		articleIDStrs, err = database.RDB.ZRevRange(ctx, zsetKey, start, end).Result()
	} else {
		// 其他排序：获取所有文章ID（最多500条）
		articleIDStrs, err = database.RDB.ZRevRange(ctx, zsetKey, 0, 499).Result()
	}
	if err != nil {
		log.Printf("从Redis ZSet获取文章ID失败: %v，降级到数据库查询", err)
		// 降级：使用数据库查询
		return s.getListFromDB(query)
	}

	if len(articleIDStrs) == 0 {
		return []*models.Article{}, 0, nil
	}

	// 转换ID为uint（Redis ZSet返回的是字符串）
	articleIDs := make([]uint, 0, len(articleIDStrs))
	for _, idStr := range articleIDStrs {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			articleIDs = append(articleIDs, uint(id))
		}
	}

	if len(articleIDs) == 0 {
		return []*models.Article{}, 0, nil
	}

	// 从数据库查询文章详情
	var articles []*models.Article
	db := database.DB.Where("id IN ?", articleIDs).Where("status = ?", 1)

	// 应用筛选条件
	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}
	if query.TagID > 0 {
		db = db.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", query.TagID)
	}
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", keyword, keyword)
	}

	// 如果指定了排序方式（非热门排序），应用排序
	if query.SortBy != "" && query.SortBy != "hot" {
		orderBy := s.buildOrderBy(query)
		db = db.Order(orderBy)
	}

	if err := db.Preload("Category").Preload("Tags").Preload("Author").
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	// 如果使用热门排序或未指定排序，按ID顺序排序（保持ZSet中的热度排序）
	if query.SortBy == "" || query.SortBy == "hot" {
		articleMap := make(map[uint]*models.Article)
		for _, article := range articles {
			articleMap[article.ID] = article
		}

		sortedArticles := make([]*models.Article, 0, len(articleIDs))
		for _, id := range articleIDs {
			if article, ok := articleMap[id]; ok {
				sortedArticles = append(sortedArticles, article)
			}
		}

		// 应用分页
		start := (query.Page - 1) * query.PageSize
		end := start + query.PageSize
		if end > len(sortedArticles) {
			end = len(sortedArticles)
		}
		if start < len(sortedArticles) {
			sortedArticles = sortedArticles[start:end]
		} else {
			sortedArticles = []*models.Article{}
		}

		// 获取总数（从ZSet获取）
		total, err := database.RDB.ZCard(ctx, zsetKey).Result()
		if err != nil {
			total = int64(len(sortedArticles))
		}

		return sortedArticles, total, nil
	}

	// 如果使用了其他排序方式，需要应用分页
	start := (query.Page - 1) * query.PageSize
	end := start + query.PageSize
	if end > len(articles) {
		end = len(articles)
	}
	if start < len(articles) {
		articles = articles[start:end]
	} else {
		articles = []*models.Article{}
	}

	// 获取总数（需要重新查询以应用筛选条件）
	var totalCount int64
	countDB := database.DB.Model(&models.Article{}).Where("id IN ?", articleIDs).Where("status = ?", 1)
	if query.CategoryID > 0 {
		countDB = countDB.Where("category_id = ?", query.CategoryID)
	}
	if query.TagID > 0 {
		countDB = countDB.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", query.TagID)
	}
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		countDB = countDB.Where("title LIKE ? OR content LIKE ?", keyword, keyword)
	}
	if err := countDB.Count(&totalCount).Error; err != nil {
		totalCount = int64(len(articles))
	}

	return articles, totalCount, nil
}

// getListFromDB 从数据库获取文章列表（降级方案）
func (s *ArticleService) getListFromDB(query *models.ArticleListQuery) ([]*models.Article, int64, error) {
	// 设置rank_type为空，使用普通查询
	query.RankType = ""
	return s.GetList(query)
}

func (s *ArticleService) IncrementViewCount(id uint) error {
	return database.DB.Model(&models.Article{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GetRecommended 获取推荐文章
func (s *ArticleService) GetRecommended(limit int) ([]*models.Article, error) {
	if limit <= 0 {
		limit = 3
	}

	var articles []*models.Article
	err := database.DB.Where("status = ? AND is_recommend = ?", 1, true).
		Preload("Category").Preload("Tags").Preload("Author").
		Order("created_at DESC").
		Limit(limit).
		Find(&articles).Error

	return articles, err
}

// SetRecommend 设置文章推荐状态
func (s *ArticleService) SetRecommend(id uint, isRecommend bool) error {
	// Clear cache
	defer utils.DeleteCachePattern("article:*")

	return database.DB.Model(&models.Article{}).
		Where("id = ?", id).
		Update("is_recommend", isRecommend).Error
}

// GetHotArticles 获取热门文章（从Redis读取）
func (s *ArticleService) GetHotArticles(limit int) ([]*models.Article, error) {
	if limit <= 0 {
		limit = 6
	}

	// 从Redis获取热门文章ID列表
	ctx := database.Ctx
	data, err := database.RDB.Get(ctx, "hot:articles").Result()
	if err != nil {
		// 如果Redis中没有数据，返回最新文章作为降级方案
		log.Printf("Redis中没有热门文章数据，使用最新文章作为降级: %v", err)
		var articles []*models.Article
		err := database.DB.Where("status = ?", 1).
			Preload("Category").Preload("Tags").Preload("Author").
			Order("created_at DESC").
			Limit(limit).
			Find(&articles).Error
		return articles, err
	}

	// 解析文章ID列表
	var articleIDs []uint
	if err := json.Unmarshal([]byte(data), &articleIDs); err != nil {
		return nil, fmt.Errorf("解析热门文章ID失败: %w", err)
	}

	// 限制数量
	if len(articleIDs) > limit {
		articleIDs = articleIDs[:limit]
	}

	if len(articleIDs) == 0 {
		return []*models.Article{}, nil
	}

	// 从数据库查询文章详情
	var articles []*models.Article
	if err := database.DB.Where("id IN ?", articleIDs).
		Preload("Category").Preload("Tags").Preload("Author").
		Find(&articles).Error; err != nil {
		return nil, err
	}

	// 按ID顺序排序（保持热度排序）
	articleMap := make(map[uint]*models.Article)
	for _, article := range articles {
		articleMap[article.ID] = article
	}

	sortedArticles := make([]*models.Article, 0, len(articleIDs))
	for _, id := range articleIDs {
		if article, ok := articleMap[id]; ok {
			sortedArticles = append(sortedArticles, article)
		}
	}

	return sortedArticles, nil
}
