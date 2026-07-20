package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/gorm"
)

type WorkspaceService struct{}

func NewWorkspaceService() *WorkspaceService { return &WorkspaceService{} }

func (s *WorkspaceService) Create(req *models.WorkspaceRequest, ownerID uint) (*models.Workspace, error) {
	workspace := &models.Workspace{
		OwnerID: ownerID, Name: req.Name, Description: req.Description, Icon: req.Icon, Sort: req.Sort,
	}
	if err := database.DB.Create(workspace).Error; err != nil {
		return nil, err
	}
	s.setCache(workspace)
	return workspace, nil
}

func (s *WorkspaceService) List(ownerID uint) ([]*models.Workspace, error) {
	var workspaces []*models.Workspace
	err := database.DB.Where("owner_id = ?", ownerID).Order("sort DESC, id DESC").Find(&workspaces).Error
	return workspaces, err
}

func (s *WorkspaceService) Get(id, ownerID uint) (*models.Workspace, error) {
	key := fmt.Sprintf("workspace:%d", id)
	var workspace models.Workspace
	if database.RDB != nil && database.GetCache(key, &workspace) == nil && workspace.ID == id && workspace.OwnerID == ownerID {
		return &workspace, nil
	}
	if err := database.DB.Where("id = ? AND owner_id = ?", id, ownerID).First(&workspace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	s.setCache(&workspace)
	return &workspace, nil
}

func (s *WorkspaceService) Update(id, ownerID uint, req *models.WorkspaceRequest) (*models.Workspace, error) {
	result := database.DB.Model(&models.Workspace{}).Where("id = ? AND owner_id = ?", id, ownerID).Updates(map[string]interface{}{
		"name": req.Name, "description": req.Description, "icon": req.Icon, "sort": req.Sort,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		if _, err := s.Get(id, ownerID); err != nil {
			return nil, err
		}
	}
	s.deleteCache(id)
	return s.Get(id, ownerID)
}

func (s *WorkspaceService) Delete(id, ownerID uint) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var workspace models.Workspace
		if err := tx.Where("id = ? AND owner_id = ?", id, ownerID).First(&workspace).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		var docIDs []uint
		if err := tx.Model(&models.Doc{}).Where("workspace_id = ? AND owner_id = ?", id, ownerID).Pluck("id", &docIDs).Error; err != nil {
			return err
		}
		if len(docIDs) > 0 {
			if err := tx.Where("doc_id IN ? AND owner_id = ?", docIDs, ownerID).Delete(&models.ShareLink{}).Error; err != nil {
				return err
			}
			if err := tx.Where("doc_id IN ? AND owner_id = ?", docIDs, ownerID).Delete(&models.DocVersion{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("workspace_id = ? AND owner_id = ?", id, ownerID).Delete(&models.Doc{}).Error; err != nil {
			return err
		}
		if err := tx.Where("workspace_id = ? AND owner_id = ?", id, ownerID).Delete(&models.Catalog{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&models.Workspace{}).Error
	})
	if err == nil {
		s.deleteCache(id)
	}
	return err
}

func (s *WorkspaceService) setCache(workspace *models.Workspace) {
	if database.RDB != nil {
		_ = database.SetCache(fmt.Sprintf("workspace:%d", workspace.ID), workspace, 10*time.Minute)
	}
}

func (s *WorkspaceService) deleteCache(id uint) {
	deleteWorkspaceCache(id)
}

func deleteWorkspaceCache(id uint) {
	if database.RDB != nil {
		_ = database.DeleteCache(fmt.Sprintf("workspace:%d", id))
	}
}
