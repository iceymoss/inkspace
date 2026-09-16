package router

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestKnowledgeAdminRoutesRequireAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := SetupAdminRouter()
	for _, item := range []struct{ method, path string }{{"GET", "overview"}, {"GET", "workspaces"}, {"GET", "workspaces/1"}, {"PUT", "workspaces/1/audit"}, {"DELETE", "workspaces/1"}, {"GET", "docs"}, {"GET", "docs/1"}, {"PUT", "docs/1/audit"}, {"DELETE", "docs/1"}, {"GET", "shares"}, {"DELETE", "shares/1"}, {"GET", "usage"}, {"GET", "quota"}, {"PUT", "quota"}, {"GET", "audit-logs"}} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(item.method, "/api/admin/knowledge/"+item.path, nil))
		if w.Code != 401 {
			t.Errorf("%s %s status=%d", item.method, item.path, w.Code)
		}
	}
}
