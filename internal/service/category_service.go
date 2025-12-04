package service

import (
	"errors"

	"mysite/internal/database"
	"mysite/internal/models"
)

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

func (s *CategoryService) Create(req *models.CategoryRequest) (*models.Category, error) {
	// Check if name exists
	var count int64
	if err := database.DB.Model(&models.Category{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("分类名称已存在")
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Sort:        req.Sort,
	}

	if err := database.DB.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Update(id uint, req *models.CategoryRequest) (*models.Category, error) {
	category, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Slug = req.Slug
	category.Description = req.Description
	category.Sort = req.Sort

	if err := database.DB.Save(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Delete(id uint) error {
	// Check if category has articles
	var count int64
	if err := database.DB.Model(&models.Article{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该分类下还有文章，无法删除")
	}

	return database.DB.Delete(&models.Category{}, id).Error
}

func (s *CategoryService) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) GetList() ([]*models.Category, error) {
	var categories []*models.Category
	if err := database.DB.Order("sort DESC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
