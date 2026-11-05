package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDiscoveryResponseHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	respond := func(c *gin.Context) { writeDiscovery(c, "application/xml; charset=utf-8", []byte("<xml>中文</xml>")) }
	r.GET("/discovery.xml", respond)
	r.HEAD("/discovery.xml", respond)
	get := httptest.NewRecorder()
	r.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/discovery.xml", nil))
	etag := get.Header().Get("ETag")
	if get.Code != 200 || etag == "" || get.Header().Get("Cache-Control") != "public, max-age=300" {
		t.Fatalf("status=%d headers=%v", get.Code, get.Header())
	}
	for _, method := range []string{http.MethodGet, http.MethodHead} {
		for _, tt := range []struct {
			name      string
			condition []string
			status    int
		}{
			{"unconditional", nil, 200},
			{"strong", []string{etag}, 304},
			{"weak", []string{"W/" + etag}, 304},
			{"list", []string{`"other", W/` + etag + `, "last"`}, 304},
			{"multiple fields", []string{`"other"`, "W/" + etag}, 304},
			{"wildcard", []string{" * "}, 304},
			{"quoted comma", []string{`"other,tag", ` + etag}, 304},
			{"empty list members", []string{`, , "other", , W/` + etag + `, ,`}, 304},
			{"different", []string{`"other"`}, 200},
			{"unquoted", []string{etag[1 : len(etag)-1]}, 200},
			{"invalid prefix", []string{"junk" + etag}, 200},
			{"weak case sensitive", []string{"w/" + etag}, 200},
		} {
			t.Run(method+"/"+tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req := httptest.NewRequest(method, "/discovery.xml", nil)
				for _, value := range tt.condition {
					req.Header.Add("If-None-Match", value)
				}
				r.ServeHTTP(w, req)
				if w.Code != tt.status || w.Header().Get("ETag") != etag || w.Header().Get("Cache-Control") != get.Header().Get("Cache-Control") {
					t.Fatalf("status=%d headers=%v", w.Code, w.Header())
				}
				if method == http.MethodHead || tt.status == 304 {
					if w.Body.Len() != 0 {
						t.Fatalf("unexpected body %q", w.Body.String())
					}
				} else if w.Body.String() != get.Body.String() {
					t.Fatal(w.Body.String())
				}
				if tt.status == 200 && (w.Header().Get("Content-Type") != get.Header().Get("Content-Type") || w.Header().Get("Content-Length") != get.Header().Get("Content-Length")) {
					t.Fatal("GET/HEAD representation headers differ")
				}
			})
		}
	}
}

func TestDiscoveryErrorHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/broken.xml", discoveryError)
	r.HEAD("/broken.xml", discoveryError)
	for _, method := range []string{http.MethodGet, http.MethodHead} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/broken.xml", nil)
		req.Header.Set("If-None-Match", "*")
		r.ServeHTTP(w, req)
		if w.Code != 500 || w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("ETag") != "" {
			t.Fatalf("status=%d headers=%v", w.Code, w.Header())
		}
		if method == http.MethodHead && w.Body.Len() != 0 {
			t.Fatal("HEAD error leaked body")
		}
	}
}
