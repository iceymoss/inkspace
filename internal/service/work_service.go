package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"

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

	// 检查作品审核配置，决定作品状态
	// Status: 0=draft, 1=published, 2=pending(待审核), 3=rejected(审核不通过)
	workStatus := req.Status
	if workStatus != 0 && workStatus != 1 {
		workStatus = 1 // 默认已发布
	}

	settingService := NewSettingService()
	workAuditSetting, err := settingService.Get(models.SettingWorkAudit)
	if err != nil {
		// 如果配置不存在，默认不审核（向后兼容）
		log.Printf("警告: 无法获取作品审核配置，默认不审核: %v", err)
	} else {
		// 如果开启了审核，设置为待审核状态（status=2）
		// 数据库存储的是字符串 '1' 或 '0'
		auditEnabled := workAuditSetting.Value == "1" ||
			workAuditSetting.Value == "true" ||
			workAuditSetting.Value == "True" ||
			workAuditSetting.Value == "TRUE"
		if auditEnabled && workStatus == 1 {
			// 只有发布状态的作品才需要审核，草稿不需要
			workStatus = 2 // 待审核
			log.Printf("作品审核已开启，作品将设置为待审核状态 (status=2)")
		} else {
			log.Printf("作品审核未开启或为草稿，作品状态: %d", workStatus)
		}
	}

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
		Status:      workStatus, // 根据审核配置设置状态
		IsRecommend: req.IsRecommend,
	}

	// 创建作品
	if err := database.DB.Create(work).Error; err != nil {
		return nil, err
	}

	// 如果审核已开启（workStatus=2），立即在同一个事务中更新 status 为 2
	// 这样可以覆盖 GORM 的 default:1 标签和数据库的默认值
	if workStatus == 2 {
		if err := database.DB.Model(work).Update("status", 2).Error; err != nil {
			log.Printf("更新作品状态为待审核失败: %v", err)
		} else {
			work.Status = 2 // 同步更新内存中的值
			log.Printf("作品创建后，已将 Status 更新为 2（待审核），ID: %d", work.ID)
		}
	}

	// 只有已发布的作品（status=1）才更新用户作品数
	// 待审核的作品（status=2）在审核通过时再更新作品数
	if workStatus == 1 {
		if err := database.DB.Model(&models.User{}).
			Where("id = ?", authorID).
			UpdateColumn("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
			log.Printf("更新用户作品数失败: %v", err)
		}
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

	// 处理作品状态：检查作品审核配置
	// Status: 0=draft, 1=published, 2=pending(待审核), 3=rejected(审核不通过)
	workStatus := req.Status
	oldStatus := work.Status

	// 如果用户尝试将作品状态设置为已发布（status=1），需要检查审核配置
	// 管理员可以绕过审核，直接设置为已发布
	if workStatus == 1 && role != "admin" {
		settingService := NewSettingService()
		workAuditSetting, err := settingService.Get(models.SettingWorkAudit)
		if err != nil {
			log.Printf("警告: 无法获取作品审核配置，默认不审核: %v", err)
		} else {
			auditEnabled := workAuditSetting.Value == "1" ||
				workAuditSetting.Value == "true" ||
				workAuditSetting.Value == "True" ||
				workAuditSetting.Value == "TRUE"

			if auditEnabled {
				// 如果审核已开启，无论作品原本是什么状态，只要尝试设置为已发布，都需要重新审核
				// 这样可以防止：先发布正常内容通过审核，然后修改为违规内容绕过审核
				workStatus = 2 // 待审核
				if oldStatus == 1 {
					log.Printf("作品审核已开启，已发布的作品更新内容后需要重新审核 (status=2)")
				} else {
					log.Printf("作品审核已开启，将作品状态从 %d 改为待审核状态 (status=2)", oldStatus)
				}
			} else {
				log.Printf("作品审核未开启，作品状态: %d", workStatus)
			}
		}
	}

	// 使用WHERE条件更新，确保权限（非管理员只能更新自己的作品）
	updateData := map[string]interface{}{
		"title":        req.Title,
		"type":         req.Type,
		"metadata":     string(metadataJSON),
		"daily_quota":  req.Type == "photography",
		"description":  req.Description,
		"cover":        req.Cover,
		"images":       string(imagesJSON),
		"link":         req.Link,
		"github_url":   req.GithubURL,
		"demo_url":     req.DemoURL,
		"tech_stack":   req.TechStack,
		"sort":         req.Sort,
		"status":       workStatus, // 使用处理后的状态
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

	// 如果状态被设置为待审核（status=2），需要显式更新以确保覆盖GORM的default:1
	if workStatus == 2 {
		if err := database.DB.Model(&models.Work{}).Where("id = ?", id).Update("status", 2).Error; err != nil {
			log.Printf("更新作品状态为待审核失败: %v", err)
		} else {
			log.Printf("作品更新后，已将 Status 更新为 2（待审核），ID: %d", id)
		}
	}

	// 处理用户作品数的更新
	// 从待审核/拒绝变为通过：增加作品数
	if (oldStatus == 2 || oldStatus == 3) && workStatus == 1 {
		if err := database.DB.Model(&models.User{}).
			Where("id = ?", work.AuthorID).
			UpdateColumn("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
			log.Printf("更新用户作品数失败: %v", err)
		}
	}
	// 从通过变为待审核/拒绝：减少作品数
	if oldStatus == 1 && (workStatus == 2 || workStatus == 3) {
		if err := database.DB.Model(&models.User{}).
			Where("id = ?", work.AuthorID).
			UpdateColumn("work_count", gorm.Expr("work_count - ?", 1)).Error; err != nil {
			log.Printf("更新用户作品数失败: %v", err)
		}
	}

	// 重新加载作品以获取最新数据
	if err := database.DB.Preload("Author").First(&work, id).Error; err != nil {
		return nil, err
	}

	return &work, nil
}

func (s *WorkService) Delete(id uint, userID uint, role string) error {
	// 先查询作品以获取状态和作者ID
	var work models.Work
	query := database.DB.Where("id = ?", id)

	// 非管理员只能删除自己的作品
	if role != "admin" {
		query = query.Where("author_id = ?", userID)
	}

	if err := query.First(&work).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("作品不存在或无权限删除")
		}
		return err
	}

	// 删除作品
	result := database.DB.Delete(&work)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("作品不存在或无权限删除")
	}

	// 只有已发布的作品（status=1）才减少用户作品数
	// 待审核（status=2）和审核不通过（status=3）的作品不计入作品数
	if work.Status == 1 {
		if err := database.DB.Model(&models.User{}).
			Where("id = ?", work.AuthorID).
			UpdateColumn("work_count", gorm.Expr("work_count - ?", 1)).Error; err != nil {
			log.Printf("更新用户作品数失败: %v", err)
		}
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
	// 注意：如果 status 为 nil，不进行状态筛选，显示所有状态的作品
	// 这样管理后台可以选择"全部"来查看所有状态的作品
	// 用户端应该始终传递 status=1 来只显示已发布的作品

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
// 只显示已发布的作品（status=1），不显示待审核（status=2）和审核不通过（status=3）的作品
func (s *WorkService) GetUserWorks(authorID uint, page, pageSize int, workType string, sortBy string) ([]*models.Work, int64, error) {
	var works []*models.Work
	var total int64

	// 只显示已发布的作品（status=1）
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
// 只返回已发布（status=1）且被推荐（is_recommend=true）的作品
func (s *WorkService) GetRecommended(limit int) ([]*models.Work, error) {
	if limit <= 0 {
		limit = 3
	}

	var works []*models.Work
	err := database.DB.Where("status = ? AND is_recommend = ?", 1, true).
		Preload("Author").
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

	// 从数据库查询作品详情，只返回已发布（status=1）的作品
	var works []*models.Work
	if err := database.DB.Where("id IN ? AND status = ?", ids, 1).Preload("Author").Find(&works).Error; err != nil {
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

// UpdateWorkStatus 更新作品审核状态
// status: 1=通过, 3=拒绝
// auditMessage: 审核消息（可选，用于记录审核通过或拒绝的原因）
func (s *WorkService) UpdateWorkStatus(id uint, status int, auditMessage string) error {
	// 先查询作品以获取当前状态和相关信息
	var work models.Work
	if err := database.DB.First(&work, id).Error; err != nil {
		return err
	}

	oldStatus := work.Status

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 更新作品状态和审核消息
		updateData := map[string]interface{}{
			"status":        status,
			"audit_message": auditMessage, // 保存审核消息
		}
		if err := tx.Model(&models.Work{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
			return err
		}

		// 如果状态从待审核(2)变为通过(1)，需要增加用户作品数
		// 如果状态从通过(1)变为拒绝(3)，需要减少用户作品数
		if oldStatus != status {
			// 从待审核变为通过：增加用户作品数
			if oldStatus == 2 && status == 1 {
				if err := tx.Model(&models.User{}).
					Where("id = ?", work.AuthorID).
					UpdateColumn("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
					return err
				}
			}

			// 从通过变为拒绝：减少用户作品数
			if oldStatus == 1 && status == 3 {
				if err := tx.Model(&models.User{}).
					Where("id = ?", work.AuthorID).
					UpdateColumn("work_count", gorm.Expr("work_count - ?", 1)).Error; err != nil {
					return err
				}
			}

			// 从拒绝变为通过：增加用户作品数
			if oldStatus == 3 && status == 1 {
				if err := tx.Model(&models.User{}).
					Where("id = ?", work.AuthorID).
					UpdateColumn("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 发送审核通知（异步，不阻塞主流程）
	// 注意：只有在状态真正改变时才发送通知
	if oldStatus != status && (status == 1 || status == 3) {
		go func() {
			notificationService := NewNotificationService()
			// 如果审核消息为空，使用默认消息
			rejectReason := auditMessage
			if rejectReason == "" && status == 3 {
				rejectReason = "审核不通过"
			}
			if err := notificationService.CreateWorkAuditNotification(id, status, rejectReason); err != nil {
				log.Printf("❌ 创建作品审核通知失败: 作品ID=%d, 状态=%d, 错误=%v", id, status, err)
			} else {
				log.Printf("✅ 成功创建作品审核通知: 作品ID=%d, 状态=%d, 作者ID=%d", id, status, work.AuthorID)
			}
		}()
	}

	return nil
}
