package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"

	"gorm.io/gorm"
)

type AdService struct{}

func NewAdService() *AdService {
	return &AdService{}
}

// ========== 广告位置管理 ==========

// CreatePosition 创建广告位置
func (s *AdService) CreatePosition(req *models.AdPositionRequest) (*models.AdPosition, error) {
	maxCount := req.MaxCount
	if maxCount <= 0 {
		maxCount = 4 // 默认最大4个
	}

	position := &models.AdPosition{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		MaxCount:    maxCount,
		Status:      req.Status,
		Sort:        req.Sort,
	}

	if err := database.DB.Create(position).Error; err != nil {
		return nil, err
	}

	return position, nil
}

// GetPositionList 获取广告位置列表
func (s *AdService) GetPositionList(page, pageSize int) ([]*models.AdPosition, int64, error) {
	var positions []*models.AdPosition
	var total int64

	offset := (page - 1) * pageSize
	if err := database.DB.Model(&models.AdPosition{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Order("sort ASC, id DESC").
		Offset(offset).Limit(pageSize).
		Find(&positions).Error; err != nil {
		return nil, 0, err
	}

	return positions, total, nil
}

// GetPositionByID 根据ID获取广告位置
func (s *AdService) GetPositionByID(id uint) (*models.AdPosition, error) {
	var position models.AdPosition
	if err := database.DB.First(&position, id).Error; err != nil {
		return nil, err
	}
	return &position, nil
}

// UpdatePosition 更新广告位置
func (s *AdService) UpdatePosition(id uint, req *models.AdPositionRequest) (*models.AdPosition, error) {
	var position models.AdPosition
	if err := database.DB.First(&position, id).Error; err != nil {
		return nil, err
	}

	maxCount := req.MaxCount
	if maxCount <= 0 {
		maxCount = 4 // 默认最大4个
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"code":        req.Code,
		"description": req.Description,
		"max_count":   maxCount,
		"status":      req.Status,
		"sort":        req.Sort,
	}

	if err := database.DB.Model(&position).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &position, nil
}

// DeletePosition 删除广告位置
func (s *AdService) DeletePosition(id uint) error {
	return database.DB.Delete(&models.AdPosition{}, id).Error
}

// ========== 广告消息管理 ==========

// CreateAdvertisement 创建广告消息
func (s *AdService) CreateAdvertisement(req *models.AdvertisementRequest) (*models.Advertisement, error) {
	ad := &models.Advertisement{
		Title:       req.Title,
		Image:       req.Image,
		Link:        req.Link,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := database.DB.Create(ad).Error; err != nil {
		return nil, err
	}

	return ad, nil
}

// GetAdvertisementList 获取广告消息列表
func (s *AdService) GetAdvertisementList(page, pageSize int) ([]*models.Advertisement, int64, error) {
	var ads []*models.Advertisement
	var total int64

	offset := (page - 1) * pageSize
	if err := database.DB.Model(&models.Advertisement{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Order("id DESC").
		Offset(offset).Limit(pageSize).
		Find(&ads).Error; err != nil {
		return nil, 0, err
	}

	return ads, total, nil
}

// GetAdvertisementByID 根据ID获取广告消息
func (s *AdService) GetAdvertisementByID(id uint) (*models.Advertisement, error) {
	var ad models.Advertisement
	if err := database.DB.First(&ad, id).Error; err != nil {
		return nil, err
	}
	return &ad, nil
}

// UpdateAdvertisement 更新广告消息
func (s *AdService) UpdateAdvertisement(id uint, req *models.AdvertisementRequest) (*models.Advertisement, error) {
	var ad models.Advertisement
	if err := database.DB.First(&ad, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"title":       req.Title,
		"image":       req.Image,
		"link":        req.Link,
		"description": req.Description,
		"status":      req.Status,
	}

	if err := database.DB.Model(&ad).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &ad, nil
}

// DeleteAdvertisement 删除广告消息
func (s *AdService) DeleteAdvertisement(id uint) error {
	return database.DB.Delete(&models.Advertisement{}, id).Error
}

// ========== 广告投放管理 ==========

// CreatePlacement 创建广告投放
func (s *AdService) CreatePlacement(req *models.AdPlacementRequest) (*models.AdPlacement, error) {
	// 检查位置是否存在并获取最大投放数量
	var position models.AdPosition
	if err := database.DB.First(&position, req.PositionID).Error; err != nil {
		return nil, errors.New("广告位置不存在")
	}

	// 检查当前投放数量是否超过最大限制
	var currentCount int64
	if err := database.DB.Model(&models.AdPlacement{}).
		Where("position_id = ? AND status = ?", req.PositionID, 1).
		Count(&currentCount).Error; err != nil {
		return nil, err
	}

	maxCount := position.MaxCount
	if maxCount <= 0 {
		maxCount = 4 // 默认最大4个
	}

	if int(currentCount) >= maxCount {
		return nil, fmt.Errorf("该广告位最多只能投放 %d 个广告，当前已有 %d 个", maxCount, currentCount)
	}

	var startTime *time.Time
	var endTime *time.Time

	if req.StartTime != nil && !req.StartTime.Time.IsZero() {
		startTime = &req.StartTime.Time
	}
	if req.EndTime != nil && !req.EndTime.Time.IsZero() {
		endTime = &req.EndTime.Time
	}

	placement := &models.AdPlacement{
		PositionID:      req.PositionID,
		AdvertisementID: req.AdvertisementID,
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          req.Status,
		Sort:            req.Sort,
	}

	if err := database.DB.Create(placement).Error; err != nil {
		return nil, err
	}

	return placement, nil
}

// GetPlacementList 获取广告投放列表
func (s *AdService) GetPlacementList(page, pageSize int, positionID *uint) ([]*models.AdPlacement, int64, error) {
	var placements []*models.AdPlacement
	var total int64

	db := database.DB.Model(&models.AdPlacement{})
	if positionID != nil {
		db = db.Where("position_id = ?", *positionID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Position").Preload("Advertisement").
		Order("sort ASC, id DESC").
		Offset(offset).Limit(pageSize).
		Find(&placements).Error; err != nil {
		return nil, 0, err
	}

	return placements, total, nil
}

// GetPlacementByID 根据ID获取广告投放
func (s *AdService) GetPlacementByID(id uint) (*models.AdPlacement, error) {
	var placement models.AdPlacement
	if err := database.DB.Preload("Position").Preload("Advertisement").
		First(&placement, id).Error; err != nil {
		return nil, err
	}
	return &placement, nil
}

// UpdatePlacement 更新广告投放
func (s *AdService) UpdatePlacement(id uint, req *models.AdPlacementRequest) (*models.AdPlacement, error) {
	var placement models.AdPlacement
	if err := database.DB.First(&placement, id).Error; err != nil {
		return nil, err
	}

	// 如果修改了位置，需要检查新位置的最大投放数量
	if req.PositionID != placement.PositionID {
		var position models.AdPosition
		if err := database.DB.First(&position, req.PositionID).Error; err != nil {
			return nil, errors.New("广告位置不存在")
		}

		// 检查新位置的当前投放数量
		var currentCount int64
		if err := database.DB.Model(&models.AdPlacement{}).
			Where("position_id = ? AND status = ? AND id != ?", req.PositionID, 1, id).
			Count(&currentCount).Error; err != nil {
			return nil, err
		}

		maxCount := position.MaxCount
		if maxCount <= 0 {
			maxCount = 4 // 默认最大4个
		}

		if int(currentCount) >= maxCount {
			return nil, fmt.Errorf("该广告位最多只能投放 %d 个广告，当前已有 %d 个", maxCount, currentCount)
		}
	}

	var startTime *time.Time
	var endTime *time.Time

	if req.StartTime != nil && !req.StartTime.Time.IsZero() {
		startTime = &req.StartTime.Time
	}
	if req.EndTime != nil && !req.EndTime.Time.IsZero() {
		endTime = &req.EndTime.Time
	}

	updates := map[string]interface{}{
		"position_id":      req.PositionID,
		"advertisement_id": req.AdvertisementID,
		"start_time":       startTime,
		"end_time":         endTime,
		"status":           req.Status,
		"sort":             req.Sort,
	}

	if err := database.DB.Model(&placement).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &placement, nil
}

// DeletePlacement 删除广告投放
func (s *AdService) DeletePlacement(id uint) error {
	return database.DB.Delete(&models.AdPlacement{}, id).Error
}

// ========== 前端获取广告 ==========

// GetAdsByPositionCode 根据位置代码获取当前有效的广告列表
func (s *AdService) GetAdsByPositionCode(code string) ([]*models.AdDisplayResponse, error) {
	now := time.Now()

	// 先查找位置
	var position models.AdPosition
	if err := database.DB.Where("code = ? AND status = ?", code, 1).First(&position).Error; err != nil {
		return []*models.AdDisplayResponse{}, nil // 位置不存在或已禁用，返回空数组
	}

	// 查询该位置下的所有投放
	var placements []*models.AdPlacement
	query := database.DB.Where("position_id = ?", position.ID).
		Where("status = ?", 1)

	// 时间过滤：开始时间必须 <= 当前时间（如果设置了）
	query = query.Where("(start_time IS NULL OR start_time <= ?)", now)
	// 时间过滤：结束时间必须 >= 当前时间（如果设置了）
	query = query.Where("(end_time IS NULL OR end_time >= ?)", now)

	if err := query.
		Order("sort ASC, id DESC").
		Find(&placements).Error; err != nil {
		return []*models.AdDisplayResponse{}, err
	}

	// 确保返回空数组而不是nil
	result := make([]*models.AdDisplayResponse, 0)
	for _, p := range placements {
		// 查询对应的广告信息
		var ad models.Advertisement
		if err := database.DB.Where("id = ? AND status = ?", p.AdvertisementID, 1).First(&ad).Error; err == nil {
			result = append(result, &models.AdDisplayResponse{
				ID:          p.ID,
				Title:       ad.Title,
				Image:       ad.Image,
				Link:        ad.Link,
				Description: ad.Description,
			})
		}
	}

	return result, nil
}

// RecordAdClick 记录广告点击
func (s *AdService) RecordAdClick(placementID uint) error {
	return database.DB.Model(&models.AdPlacement{}).
		Where("id = ?", placementID).
		UpdateColumn("click_count", gorm.Expr("click_count + ?", 1)).Error
}

// RecordAdView 记录广告展示
func (s *AdService) RecordAdView(placementID uint) error {
	return database.DB.Model(&models.AdPlacement{}).
		Where("id = ?", placementID).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}
