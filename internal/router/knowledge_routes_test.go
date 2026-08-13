package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/handler"
	"github.com/iceymoss/inkspace/internal/middleware"
)

func TestWorkspaceMemberRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	registerWorkspaceMemberRoutes(router.Group("/api"), handler.NewWorkspaceMemberHandler())
	want := map[string]bool{
		http.MethodGet + " /api/workspaces/:id/members":            false,
		http.MethodPost + " /api/workspaces/:id/members":           false,
		http.MethodPut + " /api/workspaces/:id/members/:userId":    false,
		http.MethodDelete + " /api/workspaces/:id/members/:userId": false,
	}
	for _, route := range router.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for route, found := range want {
		if !found {
			t.Errorf("route %s is not registered", route)
		}
	}
}

func TestKnowledgeRoutesRequireAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/workspaces", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/api/workspaces", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/workspaces status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}
