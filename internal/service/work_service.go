package service

import (
	"encoding/json"

	"mysite/internal/database"
	"mysite/internal/models"

	"gorm.io/gorm"
)

type WorkService struct{}

func NewWorkService() *WorkService {
	return &WorkService{}
}

func (s *WorkService) Create(req *models.WorkRequest) (*models.Work, error) {
	imagesJSON, _ := json.Marshal(req.Images)

	work := &models.Work{
		Title:       req.Title,
		Description: req.Description,
		Cover:       req.Cover,
		Images:      string(imagesJSON),
		Link:        req.Link,
		Sort:        req.Sort,
		Status:      req.Status,
	}

	if err := database.DB.Create(work).Error; err != nil {
		return nil, err
	}

	return work, nil
}

func (s *WorkService) Update(id uint, req *models.WorkRequest) (*models.Work, error) {
	work, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	imagesJSON, _ := json.Marshal(req.Images)

	work.Title = req.Title
	work.Description = req.Description
	work.Cover = req.Cover
	work.Images = string(imagesJSON)
	work.Link = req.Link
	work.Sort = req.Sort
	work.Status = req.Status

	if err := database.DB.Save(work).Error; err != nil {
		return nil, err
	}

	return work, nil
}

func (s *WorkService) Delete(id uint) error {
	return database.DB.Delete(&models.Work{}, id).Error
}

func (s *WorkService) GetByID(id uint) (*models.Work, error) {
	var work models.Work
	if err := database.DB.First(&work, id).Error; err != nil {
		return nil, err
	}
	return &work, nil
}

func (s *WorkService) GetList(page, pageSize int, status *int) ([]*models.Work, int64, error) {
	var works []*models.Work
	var total int64

	db := database.DB.Model(&models.Work{})

	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("sort DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&works).Error; err != nil {
		return nil, 0, err
	}

	return works, total, nil
}

func (s *WorkService) IncrementViewCount(id uint) error {
	return database.DB.Model(&models.Work{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

