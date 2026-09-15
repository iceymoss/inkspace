package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestKnowledgeAdminAuditValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, body := range []string{`{"audit_status":1}`, `{"audit_status":1,"reason":" "}`, `{"reason":"有效理由"}`, `{"audit_status":2,"reason":"有效理由"}`} {
		r := gin.New()
		r.PUT("/:type/:id/audit", NewKnowledgeAdminHandler().Audit)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/docs/1/audit", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		var response struct{ Code int }
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if response.Code != 400 {
			t.Fatalf("body=%s response=%s", body, w.Body.String())
		}
	}
}
