package service

import (
	"errors"
	"log"

	"mysite/internal/database"
	"mysite/internal/models"

	"gorm.io/gorm"
)

type FollowService struct {
	notificationService *NotificationService
}

func NewFollowService() *FollowService {
	return &FollowService{
		notificationService: NewNotificationService(),
	}
}

// Follow 关注用户
func (s *FollowService) Follow(followerID, followingID uint) error {
	// 不能关注自己
	if followerID == followingID {
		return errors.New("不能关注自己")
	}

	// 检查被关注用户是否存在
	var following models.User
	if err := database.DB.First(&following, followingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 检查是否已关注（排除软删除的记录）
	var count int64
	if err := database.DB.Model(&models.UserFollow{}).
		Where("follower_id = ? AND following_id = ? AND deleted_at IS NULL", followerID, followingID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("已经关注过该用户")
	}

	// 创建关注记录
	follow := &models.UserFollow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建关注记录
		if err := tx.Create(follow).Error; err != nil {
			return err
		}

		// 更新关注者的关注数
		if err := tx.Model(&models.User{}).
			Where("id = ?", followerID).
			UpdateColumn("following_count", gorm.Expr("following_count + ?", 1)).Error; err != nil {
			return err
		}

		// 更新被关注者的粉丝数
		if err := tx.Model(&models.User{}).
			Where("id = ?", followingID).
			UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	// 事务成功后，异步发送关注通知
	if err == nil {
		go func() {
			notifErr := s.notificationService.CreateFollowNotification(followerID, followingID)
			if notifErr != nil {
				log.Printf("❌ 创建关注通知失败: 用户%d -> 用户%d, 错误: %v", followerID, followingID, notifErr)
			} else {
				log.Printf("✅ 成功创建关注通知: 用户%d -> 用户%d", followerID, followingID)
			}
		}()
	}

	return err
}

// Unfollow 取消关注
func (s *FollowService) Unfollow(followerID, followingID uint) error {
	// 严格验证：只能取消自己的关注，不能取消别人的关注
	// 检查是否已关注（排除软删除的记录）
	var follow models.UserFollow
	if err := database.DB.Where("follower_id = ? AND following_id = ? AND deleted_at IS NULL", followerID, followingID).
		First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未关注该用户")
		}
		return err
	}

	// 再次验证：确保这条关注记录确实属于当前用户（双重验证，防止数据安全问题）
	if follow.FollowerID != followerID {
		return errors.New("无权取消该关注")
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 软删除关注记录（使用Where条件确保只能删除自己的关注）
		result := tx.Where("follower_id = ? AND following_id = ?", followerID, followingID).
			Delete(&models.UserFollow{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("取消关注失败")
		}

		// 更新关注者的关注数（严格使用followerID，确保只更新自己的数据）
		if err := tx.Model(&models.User{}).
			Where("id = ?", followerID).
			UpdateColumn("following_count", gorm.Expr("following_count - ?", 1)).Error; err != nil {
			return err
		}

		// 更新被关注者的粉丝数
		if err := tx.Model(&models.User{}).
			Where("id = ?", followingID).
			UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// IsFollowing 检查是否已关注（排除软删除的记录）
func (s *FollowService) IsFollowing(followerID, followingID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.UserFollow{}).
		Where("follower_id = ? AND following_id = ? AND deleted_at IS NULL", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

// GetFollowStats 获取关注统计
func (s *FollowService) GetFollowStats(userID, currentUserID uint) (*models.FollowStats, error) {
	stats := &models.FollowStats{}

	// 获取关注数和粉丝数
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	stats.FollowingCount = user.FollowingCount
	stats.FollowerCount = user.FollowerCount

	// 如果是当前登录用户查询，检查关注状态
	if currentUserID > 0 && currentUserID != userID {
		// 检查是否已关注
		isFollowing, err := s.IsFollowing(currentUserID, userID)
		if err != nil {
			return nil, err
		}
		stats.IsFollowing = isFollowing

		// 检查对方是否关注了自己（互关）
		isFollower, err := s.IsFollowing(userID, currentUserID)
		if err != nil {
			return nil, err
		}
		stats.IsFollower = isFollower
	}

	return stats, nil
}

// GetFollowingList 获取关注列表
func (s *FollowService) GetFollowingList(userID uint, page, pageSize int, currentUserID uint) ([]*models.FollowResponse, int64, error) {
	var follows []*models.UserFollow
	var total int64

	// 查询关注列表（排除软删除的记录）
	db := database.DB.Model(&models.UserFollow{}).
		Where("follower_id = ? AND deleted_at IS NULL", userID).
		Preload("Following")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&follows).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	responses := make([]*models.FollowResponse, len(follows))
	for i, follow := range follows {
		responses[i] = &models.FollowResponse{
			ID:          follow.ID,
			FollowerID:  follow.FollowerID,
			FollowingID: follow.FollowingID,
			CreatedAt:   follow.CreatedAt,
		}
		if follow.Following != nil {
			// 返回公开信息，不包含Email、Role等敏感信息
			responses[i].User = follow.Following.ToPublicResponse()
		}
		
		// 如果当前用户已登录，检查是否已关注列表中的用户
		if currentUserID > 0 && follow.FollowingID > 0 {
			isFollowing, _ := s.IsFollowing(currentUserID, follow.FollowingID)
			responses[i].IsFollowing = isFollowing
		}
	}

	return responses, total, nil
}

// GetFollowerList 获取粉丝列表
func (s *FollowService) GetFollowerList(userID uint, page, pageSize int, currentUserID uint) ([]*models.FollowResponse, int64, error) {
	var follows []*models.UserFollow
	var total int64

	// 查询粉丝列表（排除软删除的记录）
	db := database.DB.Model(&models.UserFollow{}).
		Where("following_id = ? AND deleted_at IS NULL", userID).
		Preload("Follower")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&follows).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	responses := make([]*models.FollowResponse, len(follows))
	for i, follow := range follows {
		responses[i] = &models.FollowResponse{
			ID:          follow.ID,
			FollowerID:  follow.FollowerID,
			FollowingID: follow.FollowingID,
			CreatedAt:   follow.CreatedAt,
		}
		if follow.Follower != nil {
			// 返回公开信息，不包含Email、Role等敏感信息
			responses[i].User = follow.Follower.ToPublicResponse()
		}
		
		// 如果当前用户已登录，检查是否已关注列表中的用户（粉丝）
		if currentUserID > 0 && follow.FollowerID > 0 {
			isFollowing, _ := s.IsFollowing(currentUserID, follow.FollowerID)
			responses[i].IsFollowing = isFollowing
		}
	}

	return responses, total, nil
}

