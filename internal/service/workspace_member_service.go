package service

import (
	"errors"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkspaceMemberService struct{}

func NewWorkspaceMemberService() *WorkspaceMemberService { return &WorkspaceMemberService{} }

func (s *WorkspaceMemberService) List(workspaceID, userID uint) ([]*models.WorkspaceMemberResponse, error) {
	workspace, err := authorizeWorkspace(database.DB, workspaceID, userID, WorkspacePermissionView)
	if err != nil {
		return nil, err
	}

	var owner models.User
	if err := database.DB.First(&owner, workspace.OwnerID).Error; err != nil {
		return nil, err
	}
	result := []*models.WorkspaceMemberResponse{{
		WorkspaceID: workspace.ID, UserID: owner.ID, Username: owner.Username, Nickname: owner.Nickname,
		Avatar: owner.Avatar, Role: models.WorkspaceRoleOwner, CreatedAt: workspace.CreatedAt, UpdatedAt: workspace.UpdatedAt,
	}}
	var members []struct {
		models.WorkspaceMember
		Username string
		Nickname string
		Avatar   string
	}
	if err := database.DB.Model(&models.WorkspaceMember{}).
		Select("workspace_members.*, users.username, users.nickname, users.avatar").
		Joins("JOIN users ON users.id = workspace_members.user_id AND users.deleted_at IS NULL").
		Where("workspace_members.workspace_id = ?", workspaceID).Order("workspace_members.id ASC").Scan(&members).Error; err != nil {
		return nil, err
	}
	for i := range members {
		result = append(result, &models.WorkspaceMemberResponse{
			WorkspaceID: workspaceID, UserID: members[i].UserID, Username: members[i].Username,
			Nickname: members[i].Nickname, Avatar: members[i].Avatar, Role: members[i].Role,
			CreatedAt: members[i].CreatedAt, UpdatedAt: members[i].UpdatedAt,
		})
	}
	return result, nil
}

func (s *WorkspaceMemberService) Add(workspaceID, actorID uint, req *models.WorkspaceMemberCreateRequest) (*models.WorkspaceMemberResponse, error) {
	if !validMemberRole(req.Role) {
		return nil, ErrKnowledgeInvalid
	}
	var user models.User
	var member models.WorkspaceMember
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		workspace, err := authorizeWorkspaceForUpdate(tx, workspaceID, actorID, WorkspacePermissionManageMembers)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("BINARY username = ? AND status = ?", req.Username, 1).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		if user.ID == workspace.OwnerID {
			return ErrWorkspaceOwner
		}
		result := tx.Unscoped().Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("workspace_id = ? AND user_id = ?", workspaceID, user.ID).First(&member)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			member = models.WorkspaceMember{WorkspaceID: workspaceID, UserID: user.ID, Role: req.Role}
			return tx.Create(&member).Error
		}
		if result.Error != nil {
			return result.Error
		}
		if !member.DeletedAt.Valid {
			return ErrWorkspaceMemberExists
		}
		member.Role = req.Role
		member.DeletedAt = gorm.DeletedAt{}
		return tx.Unscoped().Save(&member).Error
	})
	if err != nil {
		return nil, err
	}
	return memberResponse(&member, &user), nil
}

func (s *WorkspaceMemberService) Update(workspaceID, targetUserID, actorID uint, req *models.WorkspaceMemberUpdateRequest) (*models.WorkspaceMemberResponse, error) {
	if !validMemberRole(req.Role) {
		return nil, ErrKnowledgeInvalid
	}
	var member models.WorkspaceMember
	var user models.User
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		workspace, err := authorizeWorkspaceForUpdate(tx, workspaceID, actorID, WorkspacePermissionManageMembers)
		if err != nil {
			return err
		}
		if targetUserID == workspace.OwnerID {
			return ErrWorkspaceOwner
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("workspace_id = ? AND user_id = ?", workspaceID, targetUserID).First(&member).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND status = ?", targetUserID, 1).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		if member.Role != req.Role {
			result := tx.Model(&models.WorkspaceMember{}).Where("id = ?", member.ID).Update("role", req.Role)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return ErrKnowledgeNotFound
			}
		}
		return tx.First(&member, member.ID).Error
	})
	if err != nil {
		return nil, err
	}
	return memberResponse(&member, &user), nil
}

func (s *WorkspaceMemberService) Delete(workspaceID, targetUserID, actorID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		workspace, err := authorizeWorkspaceForUpdate(tx, workspaceID, actorID, WorkspacePermissionManageMembers)
		if err != nil {
			return err
		}
		if targetUserID == workspace.OwnerID {
			return ErrWorkspaceOwner
		}
		var member models.WorkspaceMember
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("workspace_id = ? AND user_id = ?", workspaceID, targetUserID).First(&member).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrKnowledgeNotFound
			}
			return err
		}
		result := tx.Where("id = ?", member.ID).Delete(&models.WorkspaceMember{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrKnowledgeNotFound
		}
		return nil
	})
}

func validMemberRole(role string) bool {
	return role == models.WorkspaceRoleAdmin || role == models.WorkspaceRoleEditor || role == models.WorkspaceRoleViewer
}

func memberResponse(member *models.WorkspaceMember, user *models.User) *models.WorkspaceMemberResponse {
	return &models.WorkspaceMemberResponse{
		WorkspaceID: member.WorkspaceID, UserID: user.ID, Username: user.Username, Nickname: user.Nickname,
		Avatar: user.Avatar, Role: member.Role, CreatedAt: member.CreatedAt, UpdatedAt: member.UpdatedAt,
	}
}
