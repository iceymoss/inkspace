package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWorkListStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		requestURL string
		want       *int
	}{
		{name: "public default", requestURL: "/api/works", want: intPtr(1)},
		{name: "public ignores draft filter", requestURL: "/api/works?status=0", want: intPtr(1)},
		{name: "public ignores pending filter", requestURL: "/api/works?status=2", want: intPtr(1)},
		{name: "public ignores rejected filter", requestURL: "/api/works?status=3", want: intPtr(1)},
		{name: "admin default", requestURL: "/api/admin/works", want: nil},
		{name: "admin status filter", requestURL: "/api/admin/works?status=2", want: intPtr(2)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _ := gin.CreateTestContext(httptest.NewRecorder())
			context.Request = httptest.NewRequest("GET", tt.requestURL, nil)

			got := workListStatus(context)
			if tt.want == nil {
				if got != nil {
					t.Fatalf("workListStatus() = %d, want nil", *got)
				}
				return
			}
			if got == nil || *got != *tt.want {
				if got == nil {
					t.Fatalf("workListStatus() = nil, want %d", *tt.want)
				}
				t.Fatalf("workListStatus() = %d, want %d", *got, *tt.want)
			}
		})
	}
}

func intPtr(value int) *int {
	return &value
}
