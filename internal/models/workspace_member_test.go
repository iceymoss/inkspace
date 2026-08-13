package models

import (
	"reflect"
	"testing"
)

func TestWorkspaceRolesAreDistinct(t *testing.T) {
	roles := []string{WorkspaceRoleOwner, WorkspaceRoleAdmin, WorkspaceRoleEditor, WorkspaceRoleViewer}
	seen := make(map[string]bool, len(roles))
	for _, role := range roles {
		if role == "" || seen[role] {
			t.Fatalf("workspace roles must be non-empty and distinct: %v", roles)
		}
		seen[role] = true
	}
	if !reflect.DeepEqual(roles, []string{"owner", "admin", "editor", "viewer"}) {
		t.Fatalf("workspace roles = %v", roles)
	}
}
