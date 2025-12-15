package service

import (
	"errors"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
)

type TagService struct{}

func NewTagService() *TagService {
	return &TagService{}
}

func (s *TagService) Create(req *models.TagRequest) (*models.Tag, error) {
	// Check if name exists for this user (or globally if no user)
	var count int64
	query := database.DB.Model(&models.Tag{}).Where("name = ?", req.Name)
	if req.UserID != nil {
		query = query.Where("(user_id = ? OR user_id IS NULL)", *req.UserID)
	} else {
		query = query.Where("user_id IS NULL")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("标签名称已存在")
	}

	tag := &models.Tag{
		Name:   req.Name,
		Slug:   req.Slug,
		Color:  req.Color,
		UserID: req.UserID,
	}

	if err := database.DB.Create(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *TagService) Update(id uint, req *models.TagRequest) (*models.Tag, error) {
	tag, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updateData := map[string]interface{}{
		"name":  req.Name,
		"slug":  req.Slug,
		"color": req.Color,
	}
	if req.UserID != nil {
		updateData["user_id"] = *req.UserID
	}

	if err := database.DB.Model(&models.Tag{}).
		Where("id = ?", tag.ID).
		Updates(updateData).Error; err != nil {
		return nil, err
	}

	// 重新加载以返回最新数据
	return s.GetByID(tag.ID)
}

func (s *TagService) Delete(id uint) error {
	return database.DB.Delete(&models.Tag{}, id).Error
}

func (s *TagService) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *TagService) GetList() ([]*models.Tag, error) {
	var tags []*models.Tag
	if err := database.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
