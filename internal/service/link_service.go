package service

import (
	"errors"

	"mysite/internal/database"
	"mysite/internal/models"
)

type LinkService struct{}

func NewLinkService() *LinkService {
	return &LinkService{}
}

func (s *LinkService) Create(req *models.LinkRequest) (*models.Link, error) {
	// 检查URL是否已存在
	var count int64
	if err := database.DB.Model(&models.Link{}).Where("url = ?", req.URL).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该链接已存在")
	}

	link := &models.Link{
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Description: req.Description,
		Email:       req.Email,
		Sort:        req.Sort,
		Status:      req.Status,
	}

	if err := database.DB.Create(link).Error; err != nil {
		return nil, err
	}

	return link, nil
}

func (s *LinkService) Update(id uint, req *models.LinkRequest) (*models.Link, error) {
	link, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	link.Name = req.Name
	link.URL = req.URL
	link.Logo = req.Logo
	link.Description = req.Description
	link.Email = req.Email
	link.Sort = req.Sort
	link.Status = req.Status

	if err := database.DB.Save(link).Error; err != nil {
		return nil, err
	}

	return link, nil
}

func (s *LinkService) Delete(id uint) error {
	return database.DB.Delete(&models.Link{}, id).Error
}

func (s *LinkService) GetByID(id uint) (*models.Link, error) {
	var link models.Link
	if err := database.DB.First(&link, id).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (s *LinkService) GetList(status *int) ([]*models.Link, error) {
	var links []*models.Link
	db := database.DB.Model(&models.Link{})

	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Order("sort DESC, id ASC").Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}
