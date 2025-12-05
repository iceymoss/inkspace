package service

import (
	"fmt"
	"mysite/internal/database"
	"mysite/internal/models"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// CreateCommentNotification 创建评论通知
func (s *NotificationService) CreateCommentNotification(fromUserID, toUserID uint, articleID *uint, workID *uint, commentID uint) error {
	if fromUserID == toUserID {
		return nil // 不给自己发通知
	}

	// 查询评论内容
	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		// 如果查询失败，仍然创建通知，但不包含评论内容
		var content string
		if articleID != nil {
			content = "评论了你的文章"
		} else if workID != nil {
			content = "评论了你的作品"
		}

		notification := &models.Notification{
			UserID:     toUserID,
			FromUserID: fromUserID,
			Type:       "comment",
			Content:    content,
			ArticleID:  articleID,
			WorkID:     workID,
			CommentID:  &commentID,
			IsRead:     false,
		}
		return database.DB.Create(notification).Error
	}

	// 构建通知内容：包含评论内容
	var prefix string
	if articleID != nil {
		prefix = "评论了你的文章："
	} else if workID != nil {
		prefix = "评论了你的作品："
	} else {
		prefix = "评论："
	}

	// 限制评论内容长度（避免通知内容过长）
	commentContent := comment.Content
	if len(commentContent) > 100 {
		commentContent = commentContent[:100] + "..."
	}
	content := prefix + commentContent

	notification := &models.Notification{
		UserID:     toUserID,
		FromUserID: fromUserID,
		Type:       "comment",
		Content:    content,
		ArticleID:  articleID,
		WorkID:     workID,
		CommentID:  &commentID,
		IsRead:     false,
	}

	return database.DB.Create(notification).Error
}

// CreateLikeNotification 创建点赞通知
func (s *NotificationService) CreateLikeNotification(fromUserID, toUserID uint, articleID *uint, workID *uint) error {
	if fromUserID == toUserID {
		return nil // 不给自己发通知
	}

	var content string
	if articleID != nil {
		content = "点赞了你的文章"
	} else if workID != nil {
		content = "点赞了你的作品"
	}

	notification := &models.Notification{
		UserID:     toUserID,
		FromUserID: fromUserID,
		Type:       "like",
		Content:    content,
		ArticleID:  articleID,
		WorkID:     workID,
		IsRead:     false,
	}

	return database.DB.Create(notification).Error
}

// CreateFavoriteNotification 创建收藏通知
func (s *NotificationService) CreateFavoriteNotification(fromUserID, toUserID uint, articleID *uint, workID *uint) error {
	if fromUserID == toUserID {
		return nil // 不给自己发通知
	}

	var content string
	if articleID != nil {
		content = "收藏了你的文章"
	} else if workID != nil {
		content = "收藏了你的作品"
	}

	notification := &models.Notification{
		UserID:     toUserID,
		FromUserID: fromUserID,
		Type:       "favorite",
		Content:    content,
		ArticleID:  articleID,
		WorkID:     workID,
		IsRead:     false,
	}

	return database.DB.Create(notification).Error
}

// CreateFollowNotification 创建关注通知
func (s *NotificationService) CreateFollowNotification(fromUserID, toUserID uint) error {
	if fromUserID == toUserID {
		return nil
	}

	notification := &models.Notification{
		UserID:     toUserID,
		FromUserID: fromUserID,
		Type:       "follow",
		Content:    "关注了你",
		IsRead:     false,
	}

	return database.DB.Create(notification).Error
}

// GetNotifications 获取通知列表
func (s *NotificationService) GetNotifications(userID uint, page, pageSize int, onlyUnread bool) ([]*models.Notification, int64, error) {
	var notifications []*models.Notification
	var total int64

	db := database.DB.Model(&models.Notification{}).Where("user_id = ?", userID)

	if onlyUnread {
		db = db.Where("is_read = ?", false)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Preload("FromUser").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&notifications).Error

	return notifications, total, err
}

// GetUnreadCount 获取未读通知数量
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// MarkAsRead 标记通知为已读
func (s *NotificationService) MarkAsRead(notificationID, userID uint) error {
	return database.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead 标记所有通知为已读
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

// DeleteNotification 删除通知
func (s *NotificationService) DeleteNotification(notificationID, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&models.Notification{}).Error
}

// DeleteAllRead 删除所有已读通知
func (s *NotificationService) DeleteAllRead(userID uint) error {
	return database.DB.Where("user_id = ? AND is_read = ?", userID, true).
		Delete(&models.Notification{}).Error
}

// CreateReplyNotification 创建回复通知
func (s *NotificationService) CreateReplyNotification(fromUserID, toUserID uint, articleID *uint, workID *uint, commentID uint) error {
	if fromUserID == toUserID {
		return nil // 不给自己发通知
	}

	// 查询评论内容
	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		// 如果查询失败，仍然创建通知，但不包含评论内容
		var content string
		if articleID != nil {
			content = "回复了你的评论"
		} else if workID != nil {
			content = "回复了你的评论"
		}

		notification := &models.Notification{
			UserID:     toUserID,
			FromUserID: fromUserID,
			Type:       "reply",
			Content:    content,
			ArticleID:  articleID,
			WorkID:     workID,
			CommentID:  &commentID,
			IsRead:     false,
		}
		return database.DB.Create(notification).Error
	}

	// 构建通知内容：包含回复内容
	var prefix string
	if articleID != nil {
		prefix = "回复了你的评论："
	} else if workID != nil {
		prefix = "回复了你的评论："
	} else {
		prefix = "回复："
	}

	// 限制评论内容长度（避免通知内容过长）
	commentContent := comment.Content
	if len(commentContent) > 100 {
		commentContent = commentContent[:100] + "..."
	}
	content := prefix + commentContent

	notification := &models.Notification{
		UserID:     toUserID,
		FromUserID: fromUserID,
		Type:       "reply",
		Content:    content,
		ArticleID:  articleID,
		WorkID:     workID,
		CommentID:  &commentID,
		IsRead:     false,
	}

	return database.DB.Create(notification).Error
}

// GetNotificationMessage 获取通知消息内容
func (s *NotificationService) GetNotificationMessage(notification *models.Notification) string {
	var fromUserName string
	if notification.FromUser != nil {
		fromUserName = notification.FromUser.Nickname
		if fromUserName == "" {
			fromUserName = notification.FromUser.Username
		}
	}

	message := fmt.Sprintf("%s %s", fromUserName, notification.Content)
	return message
}
