package router

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// writeDiscovery serves a fully generated representation. Callers must finish
// their fresh public reads and serialization before calling it: a conditional
// request must never hide a database failure or an article visibility change.
func writeDiscovery(c *gin.Context, contentType string, body []byte) {
	etag := fmt.Sprintf(`"%x"`, sha256.Sum256(body))
	c.Header("Cache-Control", "public, max-age=300")
	c.Header("ETag", etag)
	c.Header("Content-Type", contentType)
	if discoveryETagMatches(strings.Join(c.Request.Header.Values("If-None-Match"), ","), etag) {
		c.Status(http.StatusNotModified)
		return
	}
	c.Header("Content-Length", strconv.Itoa(len(body)))
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}
	c.Data(http.StatusOK, contentType, body)
}

// discoveryError deliberately has no representation or validator and is never
// publicly cacheable. Empty errors also preserve HEAD semantics without exposing
// query or serialization details.
func discoveryError(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("ETag", "")
	c.Header("Last-Modified", "")
	c.Status(http.StatusInternalServerError)
}

// discoveryETagMatches implements weak comparison for GET/HEAD, including lists
// and wildcard. Commas inside opaque tags are data, not list delimiters. Invalid
// syntax is ignored rather than turning a malformed condition into a cache hit.
func discoveryETagMatches(value, etag string) bool {
	value = strings.Trim(value, " \t")
	if value == "*" {
		return true
	}
	matched := false
	for value != "" {
		value = strings.TrimLeft(value, " \t")
		value = strings.TrimPrefix(value, "W/")
		if len(value) < 2 || value[0] != '"' {
			return false
		}
		end := 1
		for end < len(value) && value[end] != '"' {
			if value[end] < 0x21 || value[end] == 0x7f {
				return false
			}
			end++
		}
		if end == len(value) {
			return false
		}
		matched = matched || value[:end+1] == etag
		value = strings.TrimLeft(value[end+1:], " \t")
		if value == "" {
			return matched
		}
		if value[0] != ',' {
			return false
		}
		value = value[1:]
	}
	return false
}
