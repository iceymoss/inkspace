package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
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

	// 处理 slug：
	// - 如果前端没传或为空，则根据名称自动生成 slug
	// - 如果前端传了 slug，则校验唯一性，避免数据库抛出 1062 错误
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		baseSlug := generateSlugFromName(req.Name)

		// 确保自动生成的 slug 在数据库中唯一
		var i int64 = 0
		for {
			tmp := baseSlug
			if i > 0 {
				tmp = fmt.Sprintf("%s-%d", baseSlug, i)
			}

			var slugCount int64
			if err := database.DB.Model(&models.Category{}).
				Where("slug = ?", tmp).
				Count(&slugCount).Error; err != nil {
				return nil, err
			}

			if slugCount == 0 {
				slug = tmp
				break
			}
			i++
		}
	} else {
		// 用户自定义 slug，校验唯一性
		var slugCount int64
		if err := database.DB.Model(&models.Category{}).
			Where("slug = ?", slug).
			Count(&slugCount).Error; err != nil {
			return nil, err
		}
		if slugCount > 0 {
			return nil, errors.New("分类别名已存在")
		}
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Logo:        req.Logo,
		Sort:        req.Sort,
	}

	if err := database.DB.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Update(id uint, req *models.CategoryRequest) (*models.Category, error) {
	// 先获取当前分类数据，用于对比 slug 是否变化
	category, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updateData := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"logo":        req.Logo,
		"sort":        req.Sort,
	}

	// 仅当传入了非空 slug 且与当前 slug 不同时，才尝试更新 slug，并做唯一性校验
	slug := strings.TrimSpace(req.Slug)
	if slug != "" && slug != category.Slug {
		var count int64
		if err := database.DB.Model(&models.Category{}).
			Where("slug = ? AND id != ?", slug, id).
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("分类别名已存在")
		}
		updateData["slug"] = slug
	}

	if err := database.DB.Model(&models.Category{}).
		Where("id = ?", id).
		Updates(updateData).Error; err != nil {
		return nil, err
	}

	// 重新加载以返回最新数据
	return s.GetByID(id)
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

// GetListPaged 获取带分页的分类列表（主要用于管理后台）
func (s *CategoryService) GetListPaged(page, pageSize int) ([]*models.Category, int64, error) {
	var (
		categories []*models.Category
		total      int64
	)

	// 先统计总数
	if err := database.DB.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查询当前页数据
	if err := database.DB.
		Order("sort DESC, id ASC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// generateSlugFromName 根据分类名称生成一个基础 slug（仅做简单字符清洗和格式化）
func generateSlugFromName(name string) string {
	s := strings.TrimSpace(strings.ToLower(name))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	var builder strings.Builder
	lastHyphen := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastHyphen = false
		} else if r == '-' {
			// 合并连续的 -
			if !lastHyphen {
				builder.WriteRune(r)
				lastHyphen = true
			}
		}
	}

	result := builder.String()
	if result == "" {
		result = "category"
	}

	return result
}