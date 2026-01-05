package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"

	"gorm.io/gorm"
)

type ArticleService struct{}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

func (s *ArticleService) Create(req *models.ArticleRequest, authorID uint) (*models.Article, error) {
	// 确保 status 正确设置：0=草稿, 1=已发布, 2=私有
	status := req.Status
	if status != 0 && status != 1 && status != 2 {
		status = 1 // 默认已发布
	}

	article := &models.Article{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		Cover:       req.Cover,
		CategoryID:  req.CategoryID,
		AuthorID:    authorID,
		Status:      status,
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
	err = database.DeleteCachePattern("article:*")
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (s *ArticleService) Update(id uint, req *models.ArticleRequest, userID uint, role string) (*models.Article, error) {
	// 先查询文章，但使用WHERE条件确保权限（非管理员只能查询自己的文章）
	var article models.Article
	query := database.DB.Where("id = ?", id)

	// 非管理员只能更新自己的文章
	if role != "admin" {
		query = query.Where("author_id = ?", userID)
	}

	if err := query.Preload("Category").Preload("Tags").First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在或无权限修改")
		}
		return nil, err
	}

	// 获取旧的标签IDs和分类ID
	oldTagIDs := make([]uint, len(article.Tags))
	for i, tag := range article.Tags {
		oldTagIDs[i] = tag.ID
	}
	oldCategoryID := article.CategoryID

	err := database.DB.Transaction(func(tx *gorm.DB) error {
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

		// 注意：标签计数的更新移到后面，在更新标签关联时根据差集来处理
		// 这样可以避免重复计算（如果标签没有变化）

		// Update article - 使用WHERE条件确保权限（非管理员只能更新自己的文章）
		updateData := map[string]interface{}{
			"title":        req.Title,
			"content":      req.Content,
			"summary":      req.Summary,
			"cover":        req.Cover,
			"category_id":  req.CategoryID,
			"is_top":       req.IsTop,
			"is_recommend": req.IsRecommend,
		}

		// 确保 status 正确设置：0=草稿, 1=已发布, 2=私有
		status := req.Status
		if status == 0 || status == 1 || status == 2 {
			updateData["status"] = status
		}

		updateQuery := tx.Model(&models.Article{}).Where("id = ?", id)
		// 非管理员只能更新自己的文章
		if role != "admin" {
			updateQuery = updateQuery.Where("author_id = ?", userID)
		}

		if err := updateQuery.Updates(updateData).Error; err != nil {
			return err
		}

		// 重新加载文章以获取最新数据（用于 Association 操作）
		var updatedArticle models.Article
		if err := tx.First(&updatedArticle, id).Error; err != nil {
			return err
		}

		// 计算需要更新计数的标签ID（旧标签和新标签的差集）
		oldTagIDSet := make(map[uint]bool)
		for _, tagID := range oldTagIDs {
			oldTagIDSet[tagID] = true
		}
		newTagIDSet := make(map[uint]bool)
		for _, tagID := range req.TagIDs {
			newTagIDSet[tagID] = true
		}

		// 找出需要减少计数的标签（在旧标签中但不在新标签中）
		tagsToDecrease := make([]uint, 0)
		for tagID := range oldTagIDSet {
			if !newTagIDSet[tagID] {
				tagsToDecrease = append(tagsToDecrease, tagID)
			}
		}

		// 找出需要增加计数的标签（在新标签中但不在旧标签中）
		tagsToIncrease := make([]uint, 0)
		for tagID := range newTagIDSet {
			if !oldTagIDSet[tagID] {
				tagsToIncrease = append(tagsToIncrease, tagID)
			}
		}

		// 更新标签关联
		if err := tx.Model(&updatedArticle).Association("Tags").Clear(); err != nil {
			return err
		}

		// 如果有新标签，添加关联
		if len(req.TagIDs) > 0 {
			var tags []models.Tag
			if err := tx.Where("id IN ?", req.TagIDs).Find(&tags).Error; err != nil {
				return err
			}
			if len(tags) > 0 {
				if err := tx.Model(&updatedArticle).Association("Tags").Append(tags); err != nil {
					return err
				}
			}
		}

		// 更新标签文章数（减少旧标签的计数）
		if len(tagsToDecrease) > 0 {
			if err := tx.Model(&models.Tag{}).
				Where("id IN ?", tagsToDecrease).
				UpdateColumn("article_count", gorm.Expr("article_count - ?", 1)).Error; err != nil {
				return err
			}
		}

		// 更新标签文章数（增加新标签的计数）
		if len(tagsToIncrease) > 0 {
			if err := tx.Model(&models.Tag{}).
				Where("id IN ?", tagsToIncrease).
				UpdateColumn("article_count", gorm.Expr("article_count + ?", 1)).Error; err != nil {
				return err
			}
		}

		// 重新加载文章以获取最新数据（包括标签）
		if err := tx.Preload("Category").Preload("Tags").First(&article, id).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Clear cache
	err = database.DeleteCache(fmt.Sprintf("article:%d", id))
	if err != nil {
		return nil, err
	}

	err = database.DeleteCachePattern("article:list:*")
	if err != nil {

		return nil, err
	}

	return &article, nil
}

func (s *ArticleService) Delete(id uint, userID uint, role string) error {
	// 先查询文章以获取相关信息（用于更新计数）
	var article models.Article
	query := database.DB.Where("id = ?", id)

	// 非管理员只能查询自己的文章
	if role != "admin" {
		query = query.Where("author_id = ?", userID)
	}

	if err := query.Preload("Tags").First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在或无权限删除")
		}
		return err
	}

	// 获取标签IDs
	tagIDs := make([]uint, len(article.Tags))
	for i, tag := range article.Tags {
		tagIDs[i] = tag.ID
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除文章 - 使用WHERE条件确保权限（非管理员只能删除自己的文章）
		deleteQuery := tx.Where("id = ?", id)
		if role != "admin" {
			deleteQuery = deleteQuery.Where("author_id = ?", userID)
		}

		result := deleteQuery.Delete(&models.Article{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("文章不存在或无权限删除")
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
	database.DeleteCache(fmt.Sprintf("article:%d", id))
	database.DeleteCachePattern("article:list:*")

	return nil
}

func (s *ArticleService) GetByID(id uint) (*models.Article, error) {
	var article models.Article

	// Try to get from cache
	cacheKey := fmt.Sprintf("article:%d", id)
	if err := database.GetCache(cacheKey, &article); err == nil {
		return &article, nil
	}

	// Get from database
	if err := database.DB.Preload("Category").Preload("Tags").Preload("Author").First(&article, id).Error; err != nil {
		return nil, err
	}

	// Save to cache
	expiration := time.Duration(config.AppConfig.Cache.ArticleExpire) * time.Second
	database.SetCache(cacheKey, &article, expiration)

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
		// Default: only show published (status=1)
		// But if author_id is specified (viewing own articles) or ShowAll is true (admin), show all statuses
		if query.AuthorID == 0 && !query.ShowAll {
			// 查看所有用户的文章列表，只显示公开的（status=1）
			db = db.Where("status = ?", 1)
		}
		// If author_id > 0 (viewing own articles) or ShowAll is true (admin), don't filter by status
		// 这样可以显示所有状态：草稿(0)、公开(1)、私有(2)
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
		// 热门排序使用综合得分
		// 阅读量50% + 收藏量15% + 点赞量20% + 评论量15%
		sortField = "(view_count * 0.5 + favorite_count * 0.15 + like_count * 0.2 + comment_count * 0.15)"
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
	dbQuery := database.DB.Where("id IN ?", articleIDs).Where("status = ?", 1)

	// 应用筛选条件
	if query.CategoryID > 0 {
		dbQuery = dbQuery.Where("category_id = ?", query.CategoryID)
	}
	if query.TagID > 0 {
		dbQuery = dbQuery.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", query.TagID)
	}
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		dbQuery = dbQuery.Where("title LIKE ? OR content LIKE ?", keyword, keyword)
	}

	// 如果指定了排序方式（非热门排序），应用排序
	if query.SortBy != "" && query.SortBy != "hot" {
		orderBy := s.buildOrderBy(query)
		dbQuery = dbQuery.Order(orderBy)
	}

	if err := dbQuery.Preload("Category").Preload("Tags").Preload("Author").
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
	defer database.DeleteCachePattern("article:*")

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
			Order("is_top DESC, created_at DESC").
			Limit(limit).
			Find(&articles).Error
		if err != nil {
			return nil, err
		}
		// 降级方案也需要处理置顶逻辑
		topArticles := make([]*models.Article, 0)
		normalArticles := make([]*models.Article, 0)
		for _, article := range articles {
			if article.IsTop {
				topArticles = append(topArticles, article)
			} else {
				normalArticles = append(normalArticles, article)
			}
		}
		result := make([]*models.Article, 0, len(articles))
		result = append(result, topArticles...)
		result = append(result, normalArticles...)
		return result, nil
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

	// 从数据库查询文章详情 - 只查询公开的文章（status=1），不包含草稿和私有文章
	var articles []*models.Article
	if err := database.DB.Where("id IN ? AND status = ?", articleIDs, 1).
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

	// 将置顶文章放在最前面，然后按照热度排序
	topArticles := make([]*models.Article, 0)
	normalArticles := make([]*models.Article, 0)

	for _, article := range sortedArticles {
		if article.IsTop {
			topArticles = append(topArticles, article)
		} else {
			normalArticles = append(normalArticles, article)
		}
	}

	// 置顶文章在前，非置顶文章在后（都保持原有的热度排序）
	result := make([]*models.Article, 0, len(sortedArticles))
	result = append(result, topArticles...)
	result = append(result, normalArticles...)

	return result, nil
}
