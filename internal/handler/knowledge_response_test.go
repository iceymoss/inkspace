package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/service"
)

func TestKnowledgeErrorStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "not found", err: service.ErrKnowledgeNotFound, want: http.StatusNotFound},
		{name: "disabled share", err: service.ErrShareDisabled, want: http.StatusForbidden},
		{name: "expired share", err: service.ErrShareExpired, want: http.StatusForbidden},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)
			knowledgeError(context, test.err)
			if recorder.Code != test.want {
				t.Fatalf("knowledgeError() status = %d, want %d", recorder.Code, test.want)
			}
		})
	}
}
