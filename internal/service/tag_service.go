package service

import (
	"errors"

	"mysite/internal/database"
	"mysite/internal/models"
)

type TagService struct{}

func NewTagService() *TagService {
	return &TagService{}
}

func (s *TagService) Create(req *models.TagRequest) (*models.Tag, error) {
	// Check if name exists
	var count int64
	if err := database.DB.Model(&models.Tag{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("标签名称已存在")
	}

	tag := &models.Tag{
		Name:  req.Name,
		Slug:  req.Slug,
		Color: req.Color,
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

	tag.Name = req.Name
	tag.Slug = req.Slug
	tag.Color = req.Color

	if err := database.DB.Save(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
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

