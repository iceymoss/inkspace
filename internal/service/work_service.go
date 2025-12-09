package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"mysite/internal/database"
	"mysite/internal/models"

	"gorm.io/gorm"
)

type WorkService struct{}

func NewWorkService() *WorkService {
	return &WorkService{}
}

func (s *WorkService) Create(req *models.WorkRequest, authorID uint, role string) (*models.Work, error) {
	// 验证照片数量限制
	if req.Type == "photography" {
		maxPhotos := s.GetPhotoLimit(role)
		if len(req.Images) > maxPhotos {
			return nil, fmt.Errorf("照片数量超过限制（最多%d张）", maxPhotos)
		}
		if len(req.Images) == 0 {
			return nil, errors.New("摄影作品至少需要1张照片")
		}
		
		// 检查每日配额（摄影作品）
		canCreate, err := s.CheckDailyQuota(authorID, "photography")
		if err != nil {
			return nil, err
		}
		if !canCreate {
			return nil, errors.New("今日摄影作品发布数量已达上限（3个相册/天）")
		}
		
		// 更新相册元数据中的照片数量
		if req.Metadata == nil {
			req.Metadata = make(map[string]interface{})
		}
		req.Metadata["photo_count"] = len(req.Images)
	}

	imagesJSON, _ := json.Marshal(req.Images)
	metadataJSON, _ := json.Marshal(req.Metadata)

	// 摄影作品自动设置 daily_quota
	dailyQuota := req.Type == "photography"

	work := &models.Work{
		Title:       req.Title,
		Type:        req.Type,
		Metadata:    string(metadataJSON),
		DailyQuota:  dailyQuota,
		Description: req.Description,
		Cover:       req.Cover,
		Images:      string(imagesJSON),
		Link:        req.Link,
		GithubURL:   req.GithubURL,
		DemoURL:     req.DemoURL,
		TechStack:   req.TechStack,
		AuthorID:    authorID,
		Sort:        req.Sort,
		Status:      req.Status,
		IsRecommend: req.IsRecommend,
	}

	if err := database.DB.Create(work).Error; err != nil {
		return nil, err
	}

	return work, nil
}

func (s *WorkService) Update(id uint, req *models.WorkRequest, userID uint, role string) (*models.Work, error) {
	// 先查询作品，但使用WHERE条件确保权限（非管理员只能查询自己的作品）
	var work models.Work
	query := database.DB.Where("id = ?", id)
	
	// 非管理员只能更新自己的作品
	if role != "admin" {
		query = query.Where("author_id = ?", userID)
	}
	
	if err := query.First(&work).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("作品不存在或无权限修改")
		}
		return nil, err
	}

	// 验证照片数量限制
	if req.Type == "photography" {
		maxPhotos := s.GetPhotoLimit(role)
		if len(req.Images) > maxPhotos {
			return nil, fmt.Errorf("照片数量超过限制（最多%d张）", maxPhotos)
		}
		if len(req.Images) == 0 {
			return nil, errors.New("摄影作品至少需要1张照片")
		}
		
		// 更新相册元数据中的照片数量
		if req.Metadata == nil {
			req.Metadata = make(map[string]interface{})
		}
		req.Metadata["photo_count"] = len(req.Images)
	}

	imagesJSON, _ := json.Marshal(req.Images)
	metadataJSON, _ := json.Marshal(req.Metadata)

	// 使用WHERE条件更新，确保权限（非管理员只能更新自己的作品）
	updateData := map[string]interface{}{
		"title":        req.Title,
		"type":         req.Type,
		"metadata":     string(metadataJSON),
		"daily_quota":  req.Type == "photography",
		"description": req.Description,
		"cover":        req.Cover,
		"images":       string(imagesJSON),
		"link":         req.Link,
		"github_url":   req.GithubURL,
		"demo_url":     req.DemoURL,
		"tech_stack":   req.TechStack,
		"sort":         req.Sort,
		"status":       req.Status,
		"is_recommend": req.IsRecommend,
	}
	
	updateQuery := database.DB.Model(&models.Work{}).Where("id = ?", id)
	// 非管理员只能更新自己的作品
	if role != "admin" {
		updateQuery = updateQuery.Where("author_id = ?", userID)
	}
	
	if err := updateQuery.Updates(updateData).Error; err != nil {
		return nil, err
	}
	
	// 重新加载作品以获取最新数据
	if err := database.DB.Preload("Author").First(&work, id).Error; err != nil {
		return nil, err
	}

	return &work, nil
}

func (s *WorkService) Delete(id uint, userID uint, role string) error {
	// 使用WHERE条件删除，确保权限（非管理员只能删除自己的作品）
	deleteQuery := database.DB.Where("id = ?", id)
	
	// 非管理员只能删除自己的作品
	if role != "admin" {
		deleteQuery = deleteQuery.Where("author_id = ?", userID)
	}
	
	result := deleteQuery.Delete(&models.Work{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("作品不存在或无权限删除")
	}
	
	return nil
}

func (s *WorkService) GetByID(id uint) (*models.Work, error) {
	var work models.Work
	if err := database.DB.Preload("Author").First(&work, id).Error; err != nil {
		return nil, err
	}
	return &work, nil
}

func (s *WorkService) GetList(page, pageSize int, workType string, status *int, sortBy string) ([]*models.Work, int64, error) {
	var works []*models.Work
	var total int64

	db := database.DB.Model(&models.Work{})

	// 按类型筛选
	if workType != "" {
		db = db.Where("type = ?", workType)
	}

	// 按状态筛选
	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	
	// 根据排序参数决定排序方式
	orderBy := "sort DESC, created_at DESC" // 默认排序
	switch sortBy {
	case "hot":
		// 热度排序：综合浏览量、评论数、点赞数、收藏数
		orderBy = "(view_count * 0.4 + comment_count * 0.25 + like_count * 0.2 + favorite_count * 0.1) DESC, created_at DESC"
	case "time":
		// 时间排序：最新优先
		orderBy = "created_at DESC"
	case "view":
		// 浏览量排序
		orderBy = "view_count DESC, created_at DESC"
	case "like":
		// 点赞数排序
		orderBy = "like_count DESC, created_at DESC"
	default:
		// 默认：推荐优先，然后按时间
		orderBy = "sort DESC, created_at DESC"
	}
	
	if err := db.Preload("Author").Order(orderBy).Offset(offset).Limit(pageSize).Find(&works).Error; err != nil {
		return nil, 0, err
	}

	return works, total, nil
}

// GetMyWorks 获取用户自己的作品列表
func (s *WorkService) GetMyWorks(authorID uint, page, pageSize int, workType string, status *int) ([]*models.Work, int64, error) {
	var works []*models.Work
	var total int64

	db := database.DB.Model(&models.Work{}).Where("author_id = ?", authorID)

	// 按类型筛选
	if workType != "" {
		db = db.Where("type = ?", workType)
	}

	// 按状态筛选
	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Author").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&works).Error; err != nil {
		return nil, 0, err
	}

	return works, total, nil
}

// GetUserWorks 获取指定用户的作品列表（公开访问）
func (s *WorkService) GetUserWorks(authorID uint, page, pageSize int, workType string, sortBy string) ([]*models.Work, int64, error) {
	var works []*models.Work
	var total int64

	// 只显示已发布的作品
	db := database.DB.Model(&models.Work{}).Where("author_id = ? AND status = ?", authorID, 1)

	// 按类型筛选
	if workType != "" && workType != "all" {
		db = db.Where("type = ?", workType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	// 根据排序参数决定排序方式
	orderBy := "created_at DESC" // 默认最新排序
	switch sortBy {
	case "hot":
		// 热度排序：综合浏览量、评论数、点赞数、收藏数
		orderBy = "(view_count * 0.4 + comment_count * 0.25 + like_count * 0.2 + favorite_count * 0.1) DESC, created_at DESC"
	case "latest", "time":
		// 时间排序：最新优先
		orderBy = "created_at DESC"
	default:
		// 默认：最新排序
		orderBy = "created_at DESC"
	}

	if err := db.Preload("Author").Order(orderBy).Offset(offset).Limit(pageSize).Find(&works).Error; err != nil {
		return nil, 0, err
	}

	return works, total, nil
}

func (s *WorkService) IncrementViewCount(id uint) error {
	return database.DB.Model(&models.Work{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GetRecommended 获取推荐作品
func (s *WorkService) GetRecommended(limit int) ([]*models.Work, error) {
	if limit <= 0 {
		limit = 3
	}

	var works []*models.Work
	err := database.DB.Where("status = ? AND is_recommend = ?", 1, true).
		Order("sort DESC, created_at DESC").
		Limit(limit).
		Find(&works).Error

	return works, err
}

// SetRecommend 设置作品推荐状态
func (s *WorkService) SetRecommend(id uint, isRecommend bool) error {
	return database.DB.Model(&models.Work{}).
		Where("id = ?", id).
		Update("is_recommend", isRecommend).Error
}

// GetHotWorks 获取热门作品（从Redis ZSET读取，支持分页）
func (s *WorkService) GetHotWorks(page, pageSize int) ([]*models.Work, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	ctx := database.Ctx
	key := "hot:works:zset"

	// 从Redis ZSET获取总数
	total, err := database.RDB.ZCard(ctx, key).Result()
	if err != nil || total == 0 {
		// 如果Redis中没有数据，返回最新作品作为降级方案
		log.Printf("Redis中没有热门作品数据，使用最新作品作为降级: %v", err)
		var works []*models.Work
		var count int64
		offset := (page - 1) * pageSize
		err := database.DB.Model(&models.Work{}).Where("status = ?", 1).Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
		err = database.DB.Where("status = ?", 1).
			Order("created_at DESC").
			Offset(offset).
			Limit(pageSize).
			Preload("Author").
			Find(&works).Error
		return works, count, err
	}

	// 计算分页
	offset := (page - 1) * pageSize
	start := int64(offset)
	end := start + int64(pageSize) - 1

	// 从ZSET获取作品ID（按分数降序，即从高到低）
	workIDs, err := database.RDB.ZRevRange(ctx, key, start, end).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("从Redis ZSET获取作品ID失败: %w", err)
	}

	if len(workIDs) == 0 {
		return []*models.Work{}, total, nil
	}

	// 将字符串ID转换为uint
	ids := make([]uint, 0, len(workIDs))
	for _, idStr := range workIDs {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			ids = append(ids, uint(id))
		}
	}

	if len(ids) == 0 {
		return []*models.Work{}, total, nil
	}

	// 从数据库查询作品详情
	var works []*models.Work
	if err := database.DB.Where("id IN ?", ids).Preload("Author").Find(&works).Error; err != nil {
		return nil, 0, err
	}

	// 按ID顺序排序（保持热度排序）
	workMap := make(map[uint]*models.Work)
	for _, work := range works {
		workMap[work.ID] = work
	}

	sortedWorks := make([]*models.Work, 0, len(ids))
	for _, id := range ids {
		if work, ok := workMap[id]; ok {
			sortedWorks = append(sortedWorks, work)
		}
	}

	return sortedWorks, total, nil
}

// CheckDailyQuota 检查用户当日是否还能发布摄影作品
func (s *WorkService) CheckDailyQuota(userID uint, workType string) (bool, error) {
	// 只有摄影作品需要检查配额
	if workType != "photography" {
		return true, nil
	}

	// 获取今天的日期（使用 DATE 函数）
	var count int64
	err := database.DB.Model(&models.Work{}).
		Where("author_id = ? AND type = ? AND DATE(created_at) = CURDATE()", userID, workType).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	// 每天最多3张
	return count < 3, nil
}

// GetTodayQuotaUsage 获取今日配额使用情况
func (s *WorkService) GetTodayQuotaUsage(userID uint) (int, error) {
	var count int64
	err := database.DB.Model(&models.Work{}).
		Where("author_id = ? AND type = ? AND DATE(created_at) = CURDATE()", userID, "photography").
		Count(&count).Error

	return int(count), err
}

// GetPhotoLimit 获取用户的照片数量限制
func (s *WorkService) GetPhotoLimit(role string) int {
	if role == "admin" {
		return 50 // 管理员50张
	}
	return 10 // 普通用户10张
}
