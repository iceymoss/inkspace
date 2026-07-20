package service

import (
	"errors"
	"sort"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/gorm"
)

type CatalogService struct{}

func NewCatalogService() *CatalogService { return &CatalogService{} }

func (s *CatalogService) Create(workspaceID, ownerID uint, req *models.CatalogCreateRequest) (*models.Catalog, error) {
	if err := validateWorkspaceAndParent(database.DB, workspaceID, ownerID, req.ParentID); err != nil {
		return nil, err
	}
	catalog := &models.Catalog{
		WorkspaceID: workspaceID, ParentID: req.ParentID, OwnerID: ownerID, Name: req.Name, Sort: req.Sort,
	}
	if err := database.DB.Create(catalog).Error; err != nil {
		return nil, err
	}
	return catalog, nil
}

func (s *CatalogService) Tree(workspaceID, ownerID uint) ([]*models.CatalogResponse, error) {
	if err := requireWorkspace(database.DB, workspaceID, ownerID); err != nil {
		return nil, err
	}
	var catalogs []models.Catalog
	if err := database.DB.Where("workspace_id = ? AND owner_id = ?", workspaceID, ownerID).Find(&catalogs).Error; err != nil {
		return nil, err
	}
	return buildCatalogTree(catalogs), nil
}

func (s *CatalogService) Update(id, ownerID uint, req *models.CatalogUpdateRequest) (*models.Catalog, error) {
	if _, err := s.get(id, ownerID); err != nil {
		return nil, err
	}
	result := database.DB.Model(&models.Catalog{}).Where("id = ? AND owner_id = ?", id, ownerID).Update("name", req.Name)
	if result.Error != nil {
		return nil, result.Error
	}
	return s.get(id, ownerID)
}

func (s *CatalogService) Move(id, ownerID uint, req *models.CatalogMoveRequest) (*models.Catalog, error) {
	catalog, err := s.get(id, ownerID)
	if err != nil {
		return nil, err
	}
	var catalogs []models.Catalog
	if err := database.DB.Where("workspace_id = ? AND owner_id = ?", catalog.WorkspaceID, ownerID).Find(&catalogs).Error; err != nil {
		return nil, err
	}
	if wouldCreateCatalogCycle(id, req.ParentID, catalogs) {
		return nil, ErrCatalogCycle
	}
	if req.ParentID != nil {
		found := false
		for _, item := range catalogs {
			if item.ID == *req.ParentID {
				found = true
				break
			}
		}
		if !found {
			return nil, ErrKnowledgeNotFound
		}
	}
	if err := database.DB.Model(&models.Catalog{}).Where("id = ? AND owner_id = ?", id, ownerID).
		Updates(map[string]interface{}{"parent_id": req.ParentID, "sort": req.Sort}).Error; err != nil {
		return nil, err
	}
	return s.get(id, ownerID)
}

func (s *CatalogService) Delete(id, ownerID uint) error {
	var workspaceID uint
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var catalog models.Catalog
		if err := tx.Where("id = ? AND owner_id = ?", id, ownerID).First(&catalog).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		workspaceID = catalog.WorkspaceID
		var catalogs []models.Catalog
		if err := tx.Where("workspace_id = ? AND owner_id = ?", catalog.WorkspaceID, ownerID).Find(&catalogs).Error; err != nil {
			return err
		}
		catalogIDs := catalogDescendantIDs(id, catalogs)
		var docIDs []uint
		if err := tx.Model(&models.Doc{}).Where("catalog_id IN ? AND owner_id = ?", catalogIDs, ownerID).Pluck("id", &docIDs).Error; err != nil {
			return err
		}
		if len(docIDs) > 0 {
			if err := tx.Where("doc_id IN ? AND owner_id = ?", docIDs, ownerID).Delete(&models.ShareLink{}).Error; err != nil {
				return err
			}
			if err := tx.Where("doc_id IN ? AND owner_id = ?", docIDs, ownerID).Delete(&models.DocVersion{}).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN ? AND owner_id = ?", docIDs, ownerID).Delete(&models.Doc{}).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.Workspace{}).Where("id = ? AND owner_id = ?", catalog.WorkspaceID, ownerID).
				UpdateColumn("doc_count", gorm.Expr("GREATEST(doc_count - ?, 0)", len(docIDs))).Error; err != nil {
				return err
			}
		}
		return tx.Where("id IN ? AND owner_id = ?", catalogIDs, ownerID).Delete(&models.Catalog{}).Error
	})
	if err == nil {
		deleteWorkspaceCache(workspaceID)
	}
	return err
}

func (s *CatalogService) get(id, ownerID uint) (*models.Catalog, error) {
	var catalog models.Catalog
	if err := database.DB.Where("id = ? AND owner_id = ?", id, ownerID).First(&catalog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	return &catalog, nil
}

func requireWorkspace(db *gorm.DB, workspaceID, ownerID uint) error {
	var count int64
	if err := db.Model(&models.Workspace{}).Where("id = ? AND owner_id = ?", workspaceID, ownerID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrKnowledgeNotFound
	}
	return nil
}

func validateWorkspaceAndParent(db *gorm.DB, workspaceID, ownerID uint, parentID *uint) error {
	if err := requireWorkspace(db, workspaceID, ownerID); err != nil {
		return err
	}
	if parentID == nil {
		return nil
	}
	var count int64
	if err := db.Model(&models.Catalog{}).Where("id = ? AND workspace_id = ? AND owner_id = ?", *parentID, workspaceID, ownerID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrKnowledgeNotFound
	}
	return nil
}

func buildCatalogTree(catalogs []models.Catalog) []*models.CatalogResponse {
	sort.SliceStable(catalogs, func(i, j int) bool {
		if catalogs[i].Sort == catalogs[j].Sort {
			return catalogs[i].ID < catalogs[j].ID
		}
		return catalogs[i].Sort < catalogs[j].Sort
	})
	nodes := make(map[uint]*models.CatalogResponse, len(catalogs))
	for i := range catalogs {
		nodes[catalogs[i].ID] = catalogs[i].ToResponse()
	}
	roots := make([]*models.CatalogResponse, 0)
	for i := range catalogs {
		node := nodes[catalogs[i].ID]
		if catalogs[i].ParentID != nil {
			if parent, ok := nodes[*catalogs[i].ParentID]; ok {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	return roots
}

func catalogDescendantIDs(rootID uint, catalogs []models.Catalog) []uint {
	children := make(map[uint][]uint)
	for _, catalog := range catalogs {
		if catalog.ParentID != nil {
			children[*catalog.ParentID] = append(children[*catalog.ParentID], catalog.ID)
		}
	}
	ids := make([]uint, 0)
	seen := make(map[uint]bool)
	stack := []uint{rootID}
	for len(stack) > 0 {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
		stack = append(stack, children[id]...)
	}
	return ids
}

func wouldCreateCatalogCycle(catalogID uint, parentID *uint, catalogs []models.Catalog) bool {
	if parentID == nil {
		return false
	}
	for _, id := range catalogDescendantIDs(catalogID, catalogs) {
		if id == *parentID {
			return true
		}
	}
	return false
}
