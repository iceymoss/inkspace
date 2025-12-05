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
		Select("id, view_count, comment_count").
		Find(&works).Error; err != nil {
		return fmt.Errorf("查询作品失败: %w", err)
	}

	if len(works) == 0 {
		log.Println("没有已发布的作品")
		return nil
	}

	// 2. 计算每个作品的得分
	scores := make([]WorkScore, 0, len(works))
	for _, work := range works {
		score := t.calculateScore(work)
		scores = append(scores, WorkScore{
			ID:    work.ID,
			Score: score,
		})
	}

	// 3. 按得分排序
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].Score > scores[i].Score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// 4. 取前10个热门作品ID
	topCount := 10
	if len(scores) < topCount {
		topCount = len(scores)
	}

	hotWorkIDs := make([]uint, topCount)
	for i := 0; i < topCount; i++ {
		hotWorkIDs[i] = scores[i].ID
	}

	// 5. 存储到Redis
	key := "hot:works"
	data, err := json.Marshal(hotWorkIDs)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	// 设置过期时间为6分钟（略长于任务间隔，防止缓存失效）
	if err := database.RDB.Set(ctx, key, data, 6*time.Minute).Err(); err != nil {
		return fmt.Errorf("存储到Redis失败: %w", err)
	}

	log.Printf("✅ 热门作品计算完成，共 %d 个作品，已存储前 %d 个到Redis", len(works), topCount)
	log.Printf("📊 Top 3 热门作品: %v", hotWorkIDs[:min(3, len(hotWorkIDs))])

	return nil
}

// calculateScore 计算作品得分
// 权重：浏览量 60%、评论数 40%
func (t *HotWorksTask) calculateScore(work models.Work) float64 {
	// 归一化处理：使用对数函数降低极端值的影响
	viewScore := math.Log1p(float64(work.ViewCount)) * 0.6     // 60%
	commentScore := math.Log1p(float64(work.CommentCount)) * 0.4 // 40%

	totalScore := viewScore + commentScore

	return totalScore
}

