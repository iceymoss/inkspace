package service

import (
	"errors"

	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkspacePermission string

const (
	WorkspacePermissionView            WorkspacePermission = "view"
	WorkspacePermissionEdit            WorkspacePermission = "edit"
	WorkspacePermissionManageMembers   WorkspacePermission = "manage_members"
	WorkspacePermissionDeleteWorkspace WorkspacePermission = "delete_workspace"
)

func HasWorkspacePermission(role string, permission WorkspacePermission) bool {
	switch permission {
	case WorkspacePermissionView:
		return role == models.WorkspaceRoleOwner || role == models.WorkspaceRoleAdmin ||
			role == models.WorkspaceRoleEditor || role == models.WorkspaceRoleViewer
	case WorkspacePermissionEdit:
		return role == models.WorkspaceRoleOwner || role == models.WorkspaceRoleAdmin || role == models.WorkspaceRoleEditor
	case WorkspacePermissionManageMembers:
		return role == models.WorkspaceRoleOwner || role == models.WorkspaceRoleAdmin
	case WorkspacePermissionDeleteWorkspace:
		return role == models.WorkspaceRoleOwner
	default:
		return false
	}
}

func workspaceCapabilities(db *gorm.DB, workspace *models.Workspace, userID uint) (*models.WorkspaceCapabilities, error) {
	role := models.WorkspaceRoleOwner
	if workspace.OwnerID != userID {
		var member models.WorkspaceMember
		if err := db.Where("workspace_id = ? AND user_id = ?", workspace.ID, userID).First(&member).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrKnowledgeNotFound
			}
			return nil, err
		}
		role = member.Role
	}
	return &models.WorkspaceCapabilities{
		CanEdit:          HasWorkspacePermission(role, WorkspacePermissionEdit),
		CanManageMembers: HasWorkspacePermission(role, WorkspacePermissionManageMembers),
		CanDelete:        HasWorkspacePermission(role, WorkspacePermissionDeleteWorkspace),
	}, nil
}

func authorizeWorkspace(db *gorm.DB, workspaceID, userID uint, permission WorkspacePermission) (*models.Workspace, error) {
	return authorizeWorkspaceQuery(db, workspaceID, userID, permission, false)
}

func authorizeWorkspaceForUpdate(db *gorm.DB, workspaceID, userID uint, permission WorkspacePermission) (*models.Workspace, error) {
	return authorizeWorkspaceQuery(db, workspaceID, userID, permission, true)
}

func authorizeWorkspaceQuery(db *gorm.DB, workspaceID, userID uint, permission WorkspacePermission, lock bool) (*models.Workspace, error) {
	var workspace models.Workspace
	workspaceQuery := db.Where("id = ?", workspaceID)
	if lock {
		workspaceQuery = workspaceQuery.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := workspaceQuery.First(&workspace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}

	role := models.WorkspaceRoleOwner
	if workspace.OwnerID != userID {
		var member models.WorkspaceMember
		memberQuery := db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID)
		if lock {
			memberQuery = memberQuery.Clauses(clause.Locking{Strength: "UPDATE"})
		}
		if err := memberQuery.First(&member).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrKnowledgeNotFound
			}
			return nil, err
		}
		role = member.Role
	}
	if !HasWorkspacePermission(role, WorkspacePermissionView) {
		return nil, ErrKnowledgeNotFound
	}
	if !HasWorkspacePermission(role, permission) {
		return nil, ErrKnowledgeForbidden
	}
	return &workspace, nil
}

func workspaceForDoc(db *gorm.DB, docID, userID uint, permission WorkspacePermission) (*models.Doc, *models.Workspace, error) {
	var doc models.Doc
	if err := db.Where("id = ?", docID).First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrKnowledgeNotFound
		}
		return nil, nil, err
	}
	workspace, err := authorizeWorkspace(db, doc.WorkspaceID, userID, permission)
	if err != nil {
		return nil, nil, err
	}
	if err := requireResourceOwner(doc.OwnerID, workspace.OwnerID); err != nil {
		return nil, nil, err
	}
	return &doc, workspace, nil
}

func workspaceForCatalog(db *gorm.DB, catalogID, userID uint, permission WorkspacePermission) (*models.Catalog, *models.Workspace, error) {
	var catalog models.Catalog
	if err := db.Where("id = ?", catalogID).First(&catalog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrKnowledgeNotFound
		}
		return nil, nil, err
	}
	workspace, err := authorizeWorkspace(db, catalog.WorkspaceID, userID, permission)
	if err != nil {
		return nil, nil, err
	}
	if err := requireResourceOwner(catalog.OwnerID, workspace.OwnerID); err != nil {
		return nil, nil, err
	}
	return &catalog, workspace, nil
}

func requireResourceOwner(resourceOwnerID, workspaceOwnerID uint) error {
	if resourceOwnerID == 0 || resourceOwnerID != workspaceOwnerID {
		return ErrKnowledgeNotFound
	}
	return nil
}
