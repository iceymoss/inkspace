package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

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
	work, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 权限检查：只有作者本人或管理员可以修改
	if role != "admin" && work.AuthorID != userID {
		return nil, errors.New("无权限修改此作品")
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

	work.Title = req.Title
	work.Type = req.Type
	work.Metadata = string(metadataJSON)
	work.DailyQuota = req.Type == "photography"
	work.Description = req.Description
	work.Cover = req.Cover
	work.Images = string(imagesJSON)
	work.Link = req.Link
	work.GithubURL = req.GithubURL
	work.DemoURL = req.DemoURL
	work.TechStack = req.TechStack
	work.Sort = req.Sort
	work.Status = req.Status
	work.IsRecommend = req.IsRecommend

	if err := database.DB.Save(work).Error; err != nil {
		return nil, err
	}

	return work, nil
}

func (s *WorkService) Delete(id uint, userID uint, role string) error {
	work, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// 权限检查：只有作者本人或管理员可以删除
	if role != "admin" && work.AuthorID != userID {
		return errors.New("无权限删除此作品")
	}

	return database.DB.Delete(&models.Work{}, id).Error
}

func (s *WorkService) GetByID(id uint) (*models.Work, error) {
	var work models.Work
	if err := database.DB.Preload("Author").First(&work, id).Error; err != nil {
		return nil, err
	}
	return &work, nil
}

func (s *WorkService) GetList(page, pageSize int, workType string, status *int) ([]*models.Work, int64, error) {
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
	if err := db.Preload("Author").Order("sort DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&works).Error; err != nil {
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

// GetHotWorks 获取热门作品（从Redis读取）
func (s *WorkService) GetHotWorks(limit int) ([]*models.Work, error) {
	if limit <= 0 {
		limit = 4
	}

	// 从Redis获取热门作品ID列表
	ctx := database.Ctx
	data, err := database.RDB.Get(ctx, "hot:works").Result()
	if err != nil {
		// 如果Redis中没有数据，返回最新作品作为降级方案
		log.Printf("Redis中没有热门作品数据，使用最新作品作为降级: %v", err)
		var works []*models.Work
		err := database.DB.Where("status = ?", 1).
			Order("created_at DESC").
			Limit(limit).
			Find(&works).Error
		return works, err
	}

	// 解析作品ID列表
	var workIDs []uint
	if err := json.Unmarshal([]byte(data), &workIDs); err != nil {
		return nil, fmt.Errorf("解析热门作品ID失败: %w", err)
	}

	// 限制数量
	if len(workIDs) > limit {
		workIDs = workIDs[:limit]
	}

	if len(workIDs) == 0 {
		return []*models.Work{}, nil
	}

	// 从数据库查询作品详情
	var works []*models.Work
	if err := database.DB.Where("id IN ?", workIDs).Preload("Author").Find(&works).Error; err != nil {
		return nil, err
	}

	// 按ID顺序排序（保持热度排序）
	workMap := make(map[uint]*models.Work)
	for _, work := range works {
		workMap[work.ID] = work
	}

	sortedWorks := make([]*models.Work, 0, len(workIDs))
	for _, id := range workIDs {
		if work, ok := workMap[id]; ok {
			sortedWorks = append(sortedWorks, work)
		}
	}

	return sortedWorks, nil
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
