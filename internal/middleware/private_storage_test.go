package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBlockPrivateKnowledgeStorage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(BlockPrivateKnowledgeStorage())
	router.GET("/*path", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	tests := []struct {
		path string
		want int
	}{
		{path: "/uploads/knowledge/1/file.pdf", want: http.StatusNotFound},
		{path: "/uploads/images/public.png", want: http.StatusNoContent},
	}
	for _, test := range tests {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, test.path, nil)
		router.ServeHTTP(recorder, request)
		if recorder.Code != test.want {
			t.Errorf("GET %s status = %d, want %d", test.path, recorder.Code, test.want)
		}
	}
}
