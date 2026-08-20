package service

import (
	"errors"
	"testing"

	"github.com/iceymoss/inkspace/internal/models"
)

func TestHasWorkspacePermission(t *testing.T) {
	roles := []string{
		models.WorkspaceRoleOwner,
		models.WorkspaceRoleAdmin,
		models.WorkspaceRoleEditor,
		models.WorkspaceRoleViewer,
		"invalid",
	}
	tests := []struct {
		permission WorkspacePermission
		allowed    map[string]bool
	}{
		{permission: WorkspacePermissionView, allowed: map[string]bool{
			models.WorkspaceRoleOwner: true, models.WorkspaceRoleAdmin: true,
			models.WorkspaceRoleEditor: true, models.WorkspaceRoleViewer: true,
		}},
		{permission: WorkspacePermissionEdit, allowed: map[string]bool{
			models.WorkspaceRoleOwner: true, models.WorkspaceRoleAdmin: true, models.WorkspaceRoleEditor: true,
		}},
		{permission: WorkspacePermissionManageMembers, allowed: map[string]bool{
			models.WorkspaceRoleOwner: true, models.WorkspaceRoleAdmin: true,
		}},
		{permission: WorkspacePermissionDeleteWorkspace, allowed: map[string]bool{
			models.WorkspaceRoleOwner: true,
		}},
		{permission: WorkspacePermission("unknown"), allowed: map[string]bool{}},
	}

	for _, test := range tests {
		for _, role := range roles {
			want := test.allowed[role]
			if got := HasWorkspacePermission(role, test.permission); got != want {
				t.Errorf("HasWorkspacePermission(%q, %q) = %v, want %v", role, test.permission, got, want)
			}
		}
	}
}

func TestValidMemberRole(t *testing.T) {
	tests := map[string]bool{
		models.WorkspaceRoleOwner:  false,
		models.WorkspaceRoleAdmin:  true,
		models.WorkspaceRoleEditor: true,
		models.WorkspaceRoleViewer: true,
		"":                         false,
		"invalid":                  false,
	}
	for role, want := range tests {
		if got := validMemberRole(role); got != want {
			t.Errorf("validMemberRole(%q) = %v, want %v", role, got, want)
		}
	}
}

func TestRequireResourceOwner(t *testing.T) {
	tests := []struct {
		name             string
		resourceOwnerID  uint
		workspaceOwnerID uint
		wantErr          error
	}{
		{name: "same owner", resourceOwnerID: 7, workspaceOwnerID: 7},
		{name: "different owner", resourceOwnerID: 7, workspaceOwnerID: 8, wantErr: ErrKnowledgeNotFound},
		{name: "missing resource owner", resourceOwnerID: 0, workspaceOwnerID: 7, wantErr: ErrKnowledgeNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := requireResourceOwner(test.resourceOwnerID, test.workspaceOwnerID)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("requireResourceOwner(%d, %d) error = %v, want %v", test.resourceOwnerID, test.workspaceOwnerID, err, test.wantErr)
			}
		})
	}
}
