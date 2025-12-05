package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"mysite/internal/database"
	"mysite/internal/models"
)

// HotArticlesTask 热门文章统计任务
type HotArticlesTask struct{}

// NewHotArticlesTask 创建热门文章任务
func NewHotArticlesTask() *HotArticlesTask {
	return &HotArticlesTask{}
}

// Name 返回任务名称
func (t *HotArticlesTask) Name() string {
	return "热门文章统计"
}

// ArticleScore 文章评分结构
type ArticleScore struct {
	ID    uint    `json:"id"`
	Score float64 `json:"score"`
}

// Run 执行任务
func (t *HotArticlesTask) Run(ctx context.Context) error {
	log.Println("开始计算热门文章...")

	// 1. 获取所有已发布的文章
	var articles []models.Article
	if err := database.DB.Where("status = ?", 1).
		Select("id, view_count, comment_count, like_count, favorite_count").
		Find(&articles).Error; err != nil {
		return fmt.Errorf("查询文章失败: %w", err)
	}

	if len(articles) == 0 {
		log.Println("没有已发布的文章")
		return nil
	}

	// 2. 计算每篇文章的得分
	scores := make([]ArticleScore, 0, len(articles))
	for _, article := range articles {
		score := t.calculateScore(article)
		scores = append(scores, ArticleScore{
			ID:    article.ID,
			Score: score,
		})
	}

	// 3. 按得分排序（冒泡排序，简单实现）
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].Score > scores[i].Score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// 4. 取前20篇热门文章ID
	topCount := 20
	if len(scores) < topCount {
		topCount = len(scores)
	}

	hotArticleIDs := make([]uint, topCount)
	for i := 0; i < topCount; i++ {
		hotArticleIDs[i] = scores[i].ID
	}

	// 5. 存储到Redis
	key := "hot:articles"
	data, err := json.Marshal(hotArticleIDs)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	// 设置过期时间为20分钟（略长于任务间隔，防止缓存失效）
	if err := database.RDB.Set(ctx, key, data, 20*time.Minute).Err(); err != nil {
		return fmt.Errorf("存储到Redis失败: %w", err)
	}

	log.Printf("✅ 热门文章计算完成，共 %d 篇文章，已存储前 %d 篇到Redis", len(articles), topCount)
	log.Printf("📊 Top 5 热门文章: %v", hotArticleIDs[:min(5, len(hotArticleIDs))])

	return nil
}

// calculateScore 计算文章得分
// 权重：浏览量 50%、评论数 20%、点赞数 15%、收藏数 15%
func (t *HotArticlesTask) calculateScore(article models.Article) float64 {
	// 归一化处理：使用对数函数降低极端值的影响
	viewScore := math.Log1p(float64(article.ViewCount)) * 0.5    // 50%
	commentScore := math.Log1p(float64(article.CommentCount)) * 0.2 // 20%
	likeScore := math.Log1p(float64(article.LikeCount)) * 0.15     // 15%
	favoriteScore := math.Log1p(float64(article.FavoriteCount)) * 0.15 // 15%

	totalScore := viewScore + commentScore + likeScore + favoriteScore

	return totalScore
}

// min 返回两个整数的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

