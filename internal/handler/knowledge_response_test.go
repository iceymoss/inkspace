package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/service"
)

func TestKnowledgeErrorStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   int
	}{
		{name: "not found", err: service.ErrKnowledgeNotFound, wantStatus: http.StatusNotFound, wantCode: http.StatusNotFound},
		{name: "disabled share", err: service.ErrShareDisabled, wantStatus: http.StatusForbidden, wantCode: http.StatusForbidden},
		{name: "expired share", err: service.ErrShareExpired, wantStatus: http.StatusForbidden, wantCode: http.StatusForbidden},
		{name: "permission denied", err: service.ErrKnowledgeForbidden, wantStatus: http.StatusForbidden, wantCode: http.StatusForbidden},
		{name: "owner protected", err: service.ErrWorkspaceOwner, wantStatus: http.StatusForbidden, wantCode: http.StatusForbidden},
		{name: "member exists", err: service.ErrWorkspaceMemberExists, wantStatus: http.StatusConflict, wantCode: http.StatusConflict},
		{name: "invalid member role", err: service.ErrKnowledgeInvalid, wantStatus: http.StatusOK, wantCode: http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)
			knowledgeError(context, test.err)
			if recorder.Code != test.wantStatus {
				t.Fatalf("knowledgeError() status = %d, want %d", recorder.Code, test.wantStatus)
			}
			var body struct {
				Code int `json:"code"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if body.Code != test.wantCode {
				t.Fatalf("knowledgeError() body code = %d, want %d", body.Code, test.wantCode)
			}
		})
	}
}
