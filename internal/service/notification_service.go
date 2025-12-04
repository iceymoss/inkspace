package service

import (
	"strconv"

	"mysite/internal/database"
	"mysite/internal/models"

	"gorm.io/gorm"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// Create 创建通知
func (s *NotificationService) Create(req *models.NotificationRequest) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:     req.UserID,
		FromUserID: req.FromUserID,
		Type:       req.Type,
		Title:      req.Title,
		Content:    req.Content,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Link:       req.Link,
		IsRead:     false,
	}

	if err := database.DB.Create(notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

// GetList 获取通知列表
func (s *NotificationService) GetList(userID uint, query *models.NotificationListQuery) ([]*models.Notification, int64, error) {
	var notifications []*models.Notification
	var total int64

	db := database.DB.Model(&models.Notification{}).Where("user_id = ?", userID).Preload("FromUser")

	// 按类型筛选
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	// 按已读状态筛选
	if query.IsRead != nil {
		db = db.Where("is_read = ?", *query.IsRead)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// MarkAsRead 标记为已读
func (s *NotificationService) MarkAsRead(id, userID uint) error {
	return database.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}

// MarkAllAsRead 全部标记为已读
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}

// GetUnreadCount 获取未读数量
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// Delete 删除通知
func (s *NotificationService) Delete(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Notification{}).Error
}

// Helper functions to create specific notification types

// CreateCommentNotification 创建评论通知（评论文章时通知文章作者）
func (s *NotificationService) CreateCommentNotification(fromUserID, toUserID, articleID uint, commentContent string) error {
	if fromUserID == toUserID || toUserID == 0 {
		return nil // 不给自己发通知，游客评论不通知
	}

	// 截取评论内容前50个字符作为预览
	content := commentContent
	if len([]rune(content)) > 50 {
		content = string([]rune(content)[:50]) + "..."
	}

	notification := &models.NotificationRequest{
		UserID:     toUserID,
		FromUserID: fromUserID,
		Type:       "comment",
		Title:      "收到新评论",
		Content:    content,
		TargetType: "article",
		TargetID:   articleID,
		Link:       "/blog/" + strconv.Itoa(int(articleID)),
	}

	_, err := s.Create(notification)
	return err
}

// CreateFollowNotification 创建关注通知
func (s *NotificationService) CreateFollowNotification(fromUserID, toUserID uint) error {
	if fromUserID == toUserID {
		return nil
	}

	notification := &models.NotificationRequest{
		UserID:     toUserID,
		FromUserID: fromUserID,
		Type:       "system",
		Title:      "新增粉丝",
		Content:    "关注了你",
		TargetType: "user",
		TargetID:   fromUserID,
		Link:       "/users/" + strconv.Itoa(int(fromUserID)),
	}

	_, err := s.Create(notification)
	return err
}

// CreateLikeNotification 创建点赞文章通知
func (s *NotificationService) CreateLikeNotification(fromUserID, toUserID, articleID uint) error {
	if fromUserID == toUserID || toUserID == 0 {
		return nil // 不给自己发通知，游客点赞不通知
	}

	notification := &models.NotificationRequest{
		UserID:     toUserID,
		FromUserID: fromUserID,
		Type:       "like",
		Title:      "收到点赞",
		Content:    "赞了你的文章",
		TargetType: "article",
		TargetID:   articleID,
		Link:       "/blog/" + strconv.Itoa(int(articleID)),
	}

	_, err := s.Create(notification)
	return err
}

// CreateReplyNotification 创建回复评论通知（回复评论时通知被回复的评论作者）
func (s *NotificationService) CreateReplyNotification(fromUserID, toUserID, articleID, commentID uint, replyContent string) error {
	if fromUserID == toUserID || toUserID == 0 {
		return nil // 不给自己发通知，游客回复不通知
	}

	// 截取回复内容前50个字符作为预览
	content := replyContent
	if len([]rune(content)) > 50 {
		content = string([]rune(content)[:50]) + "..."
	}

	notification := &models.NotificationRequest{
		UserID:     toUserID,
		FromUserID: fromUserID,
		Type:       "reply",
		Title:      "收到回复",
		Content:    content,
		TargetType: "comment",
		TargetID:   commentID,
		Link:       "/blog/" + strconv.Itoa(int(articleID)),
	}

	_, err := s.Create(notification)
	return err
}

// CreateCommentLikeNotification 创建点赞评论通知（点赞评论时通知评论作者）
func (s *NotificationService) CreateCommentLikeNotification(fromUserID, toUserID, articleID, commentID uint) error {
	if fromUserID == toUserID || toUserID == 0 {
		return nil // 不给自己发通知，游客点赞不通知
	}

	notification := &models.NotificationRequest{
		UserID:     toUserID,
		FromUserID: fromUserID,
		Type:       "like",
		Title:      "收到点赞",
		Content:    "赞了你的评论",
		TargetType: "comment",
		TargetID:   commentID,
		Link:       "/blog/" + strconv.Itoa(int(articleID)),
	}

	_, err := s.Create(notification)
	return err
}

