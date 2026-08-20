package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDocWriteRequiresRevision(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name        string
		route       string
		requestPath string
		body        string
		handle      gin.HandlerFunc
	}{
		{name: "manual save", route: "/docs/:id", requestPath: "/docs/1", body: `{"title":"title","content":"body"}`, handle: NewDocHandler().Save},
		{name: "autosave", route: "/docs/:id/autosave", requestPath: "/docs/1/autosave", body: `{"content":"body"}`, handle: NewDocHandler().Autosave},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("user_id", uint(7))
				c.Next()
			})
			router.PUT(test.route, test.handle)
			request := httptest.NewRequest(http.MethodPut, test.requestPath, strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"code":400`) {
				t.Fatalf("response = %d %s, want API code 400", recorder.Code, recorder.Body.String())
			}
		})
	}
}
