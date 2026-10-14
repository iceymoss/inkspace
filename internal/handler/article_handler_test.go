package handler

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/utils"
)

func TestArticleDetailResponseReadingTime(t *testing.T) {
	content := strings.Repeat("文", 301)
	tests := []struct {
		name   string
		stored int
		want   int
	}{
		{name: "legacy article estimates without mutation", stored: 0, want: 2},
		{name: "persisted estimate takes priority", stored: 3, want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article := &models.Article{Content: content, ReadingTime: tt.stored}
			response := articleDetailResponse(article)
			if response.ReadingTime != tt.want {
				t.Fatalf("response.ReadingTime = %d, want %d", response.ReadingTime, tt.want)
			}
			if article.ReadingTime != tt.stored {
				t.Fatalf("article.ReadingTime = %d, want stored %d", article.ReadingTime, tt.stored)
			}
		})
	}
}

func TestArticleDetailResponseDoesNotChangeSource(t *testing.T) {
	gin.SetMode(gin.TestMode)
	content := strings.Repeat("文", 301)
	// Both cache hits and database reads reach the detail response as an Article.
	// Encode the source before and after as the Redis cache does to catch mutations.
	for _, source := range []string{"database", "cache"} {
		t.Run(source, func(t *testing.T) {
			article := &models.Article{Title: "正文外标题", Content: content, ReadingTime: 0}
			original, err := json.Marshal(article)
			if err != nil {
				t.Fatal(err)
			}
			if source == "cache" {
				article = &models.Article{}
				if err := json.Unmarshal(original, article); err != nil {
					t.Fatal(err)
				}
			}

			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)
			context.Request = httptest.NewRequest("GET", "/api/articles/1", nil)
			utils.Success(context, articleDetailResponse(article))
			var result struct {
				Code int `json:"code"`
				Data struct {
					ReadingTime int `json:"reading_time"`
				} `json:"data"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if result.Code != 0 || result.Data.ReadingTime != 2 {
				t.Fatalf("detail response = %+v, want code 0 and reading_time 2", result)
			}
			after, err := json.Marshal(article)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != string(original) {
				t.Fatalf("source changed: before %s, after %s", original, after)
			}
			if got := article.ToResponse().ReadingTime; got != 0 {
				t.Fatalf("other response reading time = %d, want 0", got)
			}
		})
	}
}
