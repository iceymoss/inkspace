package service

import (
	"errors"

	"mysite/internal/database"
	"mysite/internal/models"
	"mysite/internal/utils"

	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Register(req *models.UserRegisterRequest) (*models.User, error) {
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

func (s *UserService) GetUserList(page, pageSize int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	db := database.DB.Model(&models.User{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
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

