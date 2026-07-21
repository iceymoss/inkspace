package router

import (
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func serveSPA(r *gin.Engine, assets fs.FS) {
	fileServer := http.FileServer(http.FS(assets))

	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusNotFound)
			return
		}

		requestPath := strings.TrimPrefix(path.Clean(c.Request.URL.Path), "/")
		if requestPath == "." {
			requestPath = ""
		}
		if requestPath != "" {
			if info, err := fs.Stat(assets, requestPath); err == nil && !info.IsDir() {
				setStaticCacheHeaders(c, requestPath)
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		if !acceptsHTML(c.Request) || path.Ext(requestPath) != "" || strings.HasPrefix(requestPath, "assets/") ||
			strings.HasPrefix(requestPath, "api/") || strings.HasPrefix(requestPath, "uploads/") {
			c.Status(http.StatusNotFound)
			return
		}

		index, err := fs.ReadFile(assets, "index.html")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
}

func acceptsHTML(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return accept == "" || strings.Contains(accept, "text/html") || strings.Contains(accept, "*/*")
}

func setStaticCacheHeaders(c *gin.Context, name string) {
	if contentType := mime.TypeByExtension(path.Ext(name)); contentType != "" {
		c.Header("Content-Type", contentType)
	}
	if name == "index.html" {
		c.Header("Cache-Control", "no-cache")
		return
	}
	if strings.HasPrefix(name, "assets/") {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	}
}
