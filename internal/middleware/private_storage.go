package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BlockPrivateKnowledgeStorage prevents private objects from bypassing document authorization.
func BlockPrivateKnowledgeStorage() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/uploads/knowledge" || strings.HasPrefix(path, "/uploads/knowledge/") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	}
}
