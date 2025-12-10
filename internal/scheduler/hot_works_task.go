package scheduler

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"mysite/internal/database"
	"mysite/internal/models"

	"github.com/go-redis/redis/v8"
)

// HotWorksTask 热门作品统计任务
type HotWorksTask struct{}

// NewHotWorksTask 创建热门作品任务
func NewHotWorksTask() *HotWorksTask {
	return &HotWorksTask{}
}

// Name 返回任务名称
func (t *HotWorksTask) Name() string {
	return "热门作品统计"
}

// WorkScore 作品评分结构
type WorkScore struct {
	ID    uint    `json:"id"`
	Score float64 `json:"score"`
}

// Run 执行任务
func (t *HotWorksTask) Run(ctx context.Context) error {
	log.Println("开始计算热门作品...")

	// 1. 获取所有已发布的作品
	var works []models.Work
	if err := database.DB.Where("status = ?", 1).
		Select("id, view_count, comment_count, like_count, favorite_count, created_at").
		Find(&works).Error; err != nil {
		return fmt.Errorf("查询作品失败: %w", err)
	}

	if len(works) == 0 {
		log.Println("没有已发布的作品")
		return nil
	}

	// 2. 计算每个作品的得分并存储到ZSET
	key := "hot:works:zset"
	
	// 先清空旧的ZSET
	if err := database.RDB.Del(ctx, key).Err(); err != nil {
		log.Printf("警告: 清空旧ZSET失败: %v", err)
	}

	// 批量添加作品到ZSET（最多500个）
	maxCount := 500
	addedCount := 0
	
	// 先计算所有作品的得分并排序
	scores := make([]WorkScore, 0, len(works))
	for _, work := range works {
		score := t.calculateScore(work)
		scores = append(scores, WorkScore{
			ID:    work.ID,
			Score: score,
		})
	}

	// 按得分排序（降序）
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].Score > scores[i].Score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// 取前500个作品，添加到ZSET
	topCount := maxCount
	if len(scores) < topCount {
		topCount = len(scores)
	}

	// 使用Pipeline批量添加
	pipe := database.RDB.Pipeline()
	for i := 0; i < topCount; i++ {
		pipe.ZAdd(ctx, key, &redis.Z{
			Score:  scores[i].Score,
			Member: scores[i].ID,
		})
		addedCount++
	}
	
	// 设置过期时间为1小时
	pipe.Expire(ctx, key, time.Hour)
	
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("存储到Redis ZSET失败: %w", err)
	}

	log.Printf("✅ 热门作品计算完成，共 %d 个作品，已存储前 %d 个到Redis ZSET", len(works), addedCount)
	if addedCount >= 3 {
		log.Printf("📊 Top 3 热门作品ID: %v (得分: %.2f, %.2f, %.2f)", 
			[]uint{scores[0].ID, scores[1].ID, scores[2].ID},
			scores[0].Score, scores[1].Score, scores[2].Score)
	}

	return nil
}

// calculateScore 计算作品得分
// 权重：浏览量 40%、评论数 25%、点赞数 20%、收藏数 10%、时间衰减 5%
func (t *HotWorksTask) calculateScore(work models.Work) float64 {
	// 归一化处理：使用对数函数降低极端值的影响
	viewScore := math.Log1p(float64(work.ViewCount)) * 0.4        // 40%
	commentScore := math.Log1p(float64(work.CommentCount)) * 0.25  // 25%
	likeScore := math.Log1p(float64(work.LikeCount)) * 0.2         // 20%
	favoriteScore := math.Log1p(float64(work.FavoriteCount)) * 0.1 // 10%

	// 时间衰减：新作品有加分，使用指数衰减
	daysSinceCreated := time.Since(work.CreatedAt).Hours() / 24
	timeBonus := math.Exp(-daysSinceCreated/30) * 0.05 // 30天衰减周期，5%权重

	totalScore := viewScore + commentScore + likeScore + favoriteScore + timeBonus

	return totalScore
}

