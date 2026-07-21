package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func TestServeSPA(t *testing.T) {
	gin.SetMode(gin.TestMode)
	assets := fstest.MapFS{
		"index.html":        {Data: []byte("<html>app</html>")},
		"assets/app-123.js": {Data: []byte("console.log('app')")},
	}
	r := gin.New()
	r.GET("/api/known", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	serveSPA(r, assets)

	tests := []struct {
		name         string
		requestPath  string
		accept       string
		status       int
		cacheControl string
	}{
		{name: "asset", requestPath: "/assets/app-123.js", status: http.StatusOK, cacheControl: "public, max-age=31536000, immutable"},
		{name: "history fallback", requestPath: "/dashboard/workspaces/1", accept: "text/html", status: http.StatusOK, cacheControl: "no-cache"},
		{name: "known API", requestPath: "/api/known", status: http.StatusNoContent},
		{name: "unknown API", requestPath: "/api/missing", accept: "text/html", status: http.StatusNotFound},
		{name: "unknown static file", requestPath: "/assets/missing.js", accept: "text/html", status: http.StatusNotFound},
		{name: "non HTML request", requestPath: "/dashboard", accept: "application/json", status: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)
			if res.Code != tt.status {
				t.Fatalf("status = %d, want %d", res.Code, tt.status)
			}
			if got := res.Header().Get("Cache-Control"); got != tt.cacheControl {
				t.Fatalf("Cache-Control = %q, want %q", got, tt.cacheControl)
			}
		})
	}
}
