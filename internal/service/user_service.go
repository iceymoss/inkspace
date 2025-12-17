package service

import (
	"errors"
	"log"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/utils"

	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(req *models.UserRegisterRequest) (*models.User, error) {
	// 检查注册功能是否开放
	settingService := NewSettingService()
	registerSetting, err := settingService.Get(models.SettingRegisterEnabled)
	if err != nil {
		// 如果配置不存在，默认允许注册（向后兼容）
		log.Printf("警告: 无法获取注册配置，默认允许注册: %v", err)
	} else {
		if registerSetting.Value != "1" && registerSetting.Value != "true" {
			return nil, errors.New("注册功能已关闭，可向管理员申请账号")
		}
	}

	// Check if username exists
	var count int64
	if err := database.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// Check if email exists
	if err := database.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("邮箱已存在")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Nickname: nickname,
		Role:     "user",
		Status:   1,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req *models.UserLoginRequest) (string, *models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("用户名或密码错误")
		}
		return "", nil, err
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return "", nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return "", nil, errors.New("账号已被禁用")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	// 同步用户统计数据，避免各页面展示不一致
	if err := refreshUserStats(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id uint, req *models.UserUpdateRequest) (*models.User, error) {
	// 使用WHERE条件更新，确保只能更新自己的信息
	updateData := make(map[string]interface{})

	if req.Nickname != "" {
		updateData["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.Bio != "" {
		updateData["bio"] = req.Bio
	}
	if req.Avatar != "" {
		updateData["avatar"] = req.Avatar
	}

	if len(updateData) == 0 {
		// 没有需要更新的字段，直接返回用户信息
		return s.GetUserByID(id)
	}

	// 使用WHERE条件确保只能更新自己的信息
	result := database.DB.Model(&models.User{}).
		Where("id = ?", id).
		Updates(updateData)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("用户不存在")
	}

	// 重新加载用户信息
	return s.GetUserByID(id)
}

func (s *UserService) GetUserList(query *models.UserListQuery) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	db := database.DB.Model(&models.User{})

	// 关键词过滤：用户名或昵称
	if query.Username != "" {
		like := "%" + query.Username + "%"
		db = db.Where("username LIKE ? OR nickname LIKE ?", like, like)
	}

	// 邮箱模糊搜索
	if query.Email != "" {
		like := "%" + query.Email + "%"
		db = db.Where("email LIKE ?", like)
	}

	// 角色过滤
	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}

	// 状态过滤
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("id DESC").Offset(offset).Limit(query.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 同步每个用户的统计信息
	for _, u := range users {
		if err := refreshUserStats(u); err != nil {
			return nil, 0, err
		}
	}

	return users, total, nil
}

// SearchUsers 根据关键字搜索用户（用户名或昵称），只返回正常状态的用户
func (s *UserService) SearchUsers(keyword string, limit int) ([]*models.User, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	var users []*models.User
	like := "%" + keyword + "%"
	if err := database.DB.
		Where("status = ?", 1).
		Where("username LIKE ? OR nickname LIKE ?", like, like).
		Order("id DESC").
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, err
	}

	// 为搜索结果同步统计数据，确保文章数、作品数、粉丝数等与个人主页一致
	for _, u := range users {
		if err := refreshUserStats(u); err != nil {
			return nil, err
		}
	}

	return users, nil
}

// refreshUserStats 同步用户相关统计字段，避免不同页面展示不一致
func refreshUserStats(user *models.User) error {
	var count int64

	// 文章数
	if err := database.DB.Model(&models.Article{}).
		Where("author_id = ? AND deleted_at IS NULL", user.ID).
		Count(&count).Error; err != nil {
		return err
	}
	user.ArticleCount = int(count)

	// 作品数
	if err := database.DB.Model(&models.Work{}).
		Where("author_id = ? AND deleted_at IS NULL", user.ID).
		Count(&count).Error; err != nil {
		return err
	}
	user.WorkCount = int(count)

	// 评论数
	if err := database.DB.Model(&models.Comment{}).
		Where("user_id = ? AND deleted_at IS NULL", user.ID).
		Count(&count).Error; err != nil {
		return err
	}
	user.CommentCount = int(count)

	// 关注数
	if err := database.DB.Model(&models.UserFollow{}).
		Where("follower_id = ? AND deleted_at IS NULL", user.ID).
		Count(&count).Error; err != nil {
		return err
	}
	user.FollowingCount = int(count)

	// 粉丝数
	if err := database.DB.Model(&models.UserFollow{}).
		Where("following_id = ? AND deleted_at IS NULL", user.ID).
		Count(&count).Error; err != nil {
		return err
	}
	user.FollowerCount = int(count)

	// 收藏数
	if err := database.DB.Model(&models.Favorite{}).
		Where("user_id = ? AND deleted_at IS NULL", user.ID).
		Count(&count).Error; err != nil {
		return err
	}
	user.FavoriteCount = int(count)

	return nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("当前密码错误")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	if err := database.DB.Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("password", hashedPassword).Error; err != nil {
		return err
	}

	return nil
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(userID uint, status int) error {
	if err := database.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

// UpdateUserRole 更新用户角色
func (s *UserService) UpdateUserRole(userID uint, role string) error {
	if err := database.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("role", role).Error; err != nil {
		return err
	}

	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID uint) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 不允许删除管理员
	if user.Role == "admin" {
		return errors.New("不能删除管理员账号")
	}

	if err := database.DB.Delete(user).Error; err != nil {
		return err
	}

	return nil
}
