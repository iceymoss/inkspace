package service

import (
	"errors"
	"fmt"
	"strings"

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
	if err := database.DB.Order("article_count DESC, id ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetAdminList 获取管理后台用的标签列表，支持分页、筛选和排序
func (s *TagService) GetAdminList(page, pageSize int, keyword, scope string, hasArticles *int, sort string) ([]*models.Tag, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	db := database.DB.Model(&models.Tag{})

	// 关键词搜索（名称或别名）
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("name LIKE ? OR slug LIKE ?", like, like)
	}

	// 标签归属：system=系统标签（user_id IS NULL），user=用户标签（user_id IS NOT NULL）
	switch scope {
	case "system":
		db = db.Where("user_id IS NULL")
	case "user":
		db = db.Where("user_id IS NOT NULL")
	}

	// 是否有文章关联
	if hasArticles != nil {
		if *hasArticles == 1 {
			db = db.Where("article_count > 0")
		} else if *hasArticles == 0 {
			db = db.Where("article_count = 0")
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderBy := "id DESC"

	// 详细排序：支持 id/name/article_count/created_at
	if sort != "" {
		parts := strings.Split(sort, "_")
		if len(parts) == 2 {
			field := parts[0]
			dir := strings.ToUpper(parts[1])
			if dir != "ASC" && dir != "DESC" {
				dir = "DESC"
			}

			switch field {
			case "id", "name", "article_count", "created_at":
				orderBy = fmt.Sprintf("%s %s", field, dir)
			}
		}
	}

	var tags []*models.Tag
	if err := db.
		Order(orderBy).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}
