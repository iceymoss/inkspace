package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"mysite/internal/database"

	"github.com/go-redis/redis/v8"
)

// ArticleRankTask 文章榜单生成任务（周榜、月榜、年榜）
type ArticleRankTask struct{}

// NewArticleRankTask 创建文章榜单任务
func NewArticleRankTask() *ArticleRankTask {
	return &ArticleRankTask{}
}

// Name 返回任务名称
func (t *ArticleRankTask) Name() string {
	return "文章榜单生成"
}

// Run 执行任务
func (t *ArticleRankTask) Run(ctx context.Context) error {
	log.Println("开始生成文章榜单...")

	now := time.Now()

	// 检查是否是周日（一周结束）
	if now.Weekday() == time.Sunday {
		if err := t.generateWeekRank(ctx, now); err != nil {
			log.Printf("❌ 生成周榜失败: %v", err)
		} else {
			log.Println("✅ 周榜生成成功")
		}
	}

	// 检查是否是月初第一天
	if now.Day() == 1 {
		if err := t.generateMonthRank(ctx, now); err != nil {
			log.Printf("❌ 生成月榜失败: %v", err)
		} else {
			log.Println("✅ 月榜生成成功")
		}
	}

	// 检查是否是年初第一天
	if now.Month() == time.January && now.Day() == 1 {
		if err := t.generateYearRank(ctx, now); err != nil {
			log.Printf("❌ 生成年榜失败: %v", err)
		} else {
			log.Println("✅ 年榜生成成功")
		}
	}

	return nil
}

// generateWeekRank 生成周榜（从热门文章ZSet获取前50条）
func (t *ArticleRankTask) generateWeekRank(ctx context.Context, now time.Time) error {
	hotKey := "hot:articles:zset"
	weekKey := "rank:articles:week"

	// 从热门文章ZSet获取前50条（按分值降序）
	articleIDStrs, err := database.RDB.ZRevRange(ctx, hotKey, 0, 49).Result()
	if err != nil {
		return fmt.Errorf("从热门文章ZSet获取数据失败: %w", err)
	}

	if len(articleIDStrs) == 0 {
		log.Println("热门文章ZSet为空，跳过周榜生成")
		return nil
	}

	// 获取对应的分值
	zsetMembers := make([]*redis.Z, 0, len(articleIDStrs))
	for _, idStr := range articleIDStrs {
		score, err := database.RDB.ZScore(ctx, hotKey, idStr).Result()
		if err != nil {
			continue
		}
		zsetMembers = append(zsetMembers, &redis.Z{
			Score:  score,
			Member: idStr,
		})
	}

	// 清空旧的周榜
	if err := database.RDB.Del(ctx, weekKey).Err(); err != nil {
		log.Printf("⚠️ 清空旧周榜失败: %v", err)
	}

	// 存储到周榜ZSet
	if len(zsetMembers) > 0 {
		if err := database.RDB.ZAdd(ctx, weekKey, zsetMembers...).Err(); err != nil {
			return fmt.Errorf("存储周榜到Redis失败: %w", err)
		}
		// 设置过期时间为30天
		database.RDB.Expire(ctx, weekKey, 30*24*time.Hour)
	}

	log.Printf("✅ 周榜生成完成，共 %d 篇文章", len(zsetMembers))
	return nil
}

// generateMonthRank 生成月榜（从本月所有周榜中获取最热门的50条）
func (t *ArticleRankTask) generateMonthRank(ctx context.Context, now time.Time) error {
	monthKey := "rank:articles:month"

	// 获取本月的所有周榜数据
	// 这里简化处理：直接使用当前周榜作为月榜
	// 实际应该合并本月所有周榜的数据
	weekKey := "rank:articles:week"

	// 从周榜获取前50条
	articleIDStrs, err := database.RDB.ZRevRange(ctx, weekKey, 0, 49).Result()
	if err != nil {
		// 如果周榜不存在，从热门文章获取
		log.Println("周榜不存在，从热门文章生成月榜")
		return t.generateMonthRankFromHot(ctx)
	}

	if len(articleIDStrs) == 0 {
		return t.generateMonthRankFromHot(ctx)
	}

	// 获取对应的分值
	zsetMembers := make([]*redis.Z, 0, len(articleIDStrs))
	for _, idStr := range articleIDStrs {
		score, err := database.RDB.ZScore(ctx, weekKey, idStr).Result()
		if err != nil {
			continue
		}
		zsetMembers = append(zsetMembers, &redis.Z{
			Score:  score,
			Member: idStr,
		})
	}

	// 清空旧的月榜
	if err := database.RDB.Del(ctx, monthKey).Err(); err != nil {
		log.Printf("⚠️ 清空旧月榜失败: %v", err)
	}

	// 存储到月榜ZSet
	if len(zsetMembers) > 0 {
		if err := database.RDB.ZAdd(ctx, monthKey, zsetMembers...).Err(); err != nil {
			return fmt.Errorf("存储月榜到Redis失败: %w", err)
		}
		// 设置过期时间为365天
		database.RDB.Expire(ctx, monthKey, 365*24*time.Hour)
	}

	log.Printf("✅ 月榜生成完成，共 %d 篇文章", len(zsetMembers))
	return nil
}

// generateMonthRankFromHot 从热门文章生成月榜（降级方案）
func (t *ArticleRankTask) generateMonthRankFromHot(ctx context.Context) error {
	hotKey := "hot:articles:zset"
	monthKey := "rank:articles:month"

	// 从热门文章ZSet获取前50条
	articleIDStrs, err := database.RDB.ZRevRange(ctx, hotKey, 0, 49).Result()
	if err != nil {
		return fmt.Errorf("从热门文章ZSet获取数据失败: %w", err)
	}

	if len(articleIDStrs) == 0 {
		return nil
	}

	// 获取对应的分值
	zsetMembers := make([]*redis.Z, 0, len(articleIDStrs))
	for _, idStr := range articleIDStrs {
		score, err := database.RDB.ZScore(ctx, hotKey, idStr).Result()
		if err != nil {
			continue
		}
		zsetMembers = append(zsetMembers, &redis.Z{
			Score:  score,
			Member: idStr,
		})
	}

	// 清空旧的月榜
	database.RDB.Del(ctx, monthKey)

	// 存储到月榜ZSet
	if len(zsetMembers) > 0 {
		if err := database.RDB.ZAdd(ctx, monthKey, zsetMembers...).Err(); err != nil {
			return fmt.Errorf("存储月榜到Redis失败: %w", err)
		}
		database.RDB.Expire(ctx, monthKey, 365*24*time.Hour)
	}

	return nil
}

// generateYearRank 生成年榜（从月榜中获取最热门的50条）
func (t *ArticleRankTask) generateYearRank(ctx context.Context, now time.Time) error {
	yearKey := "rank:articles:year"
	monthKey := "rank:articles:month"

	// 从月榜获取前50条
	articleIDStrs, err := database.RDB.ZRevRange(ctx, monthKey, 0, 49).Result()
	if err != nil {
		// 如果月榜不存在，从热门文章获取
		log.Println("月榜不存在，从热门文章生成年榜")
		return t.generateYearRankFromHot(ctx)
	}

	if len(articleIDStrs) == 0 {
		return t.generateYearRankFromHot(ctx)
	}

	// 获取对应的分值
	zsetMembers := make([]*redis.Z, 0, len(articleIDStrs))
	for _, idStr := range articleIDStrs {
		score, err := database.RDB.ZScore(ctx, monthKey, idStr).Result()
		if err != nil {
			continue
		}
		zsetMembers = append(zsetMembers, &redis.Z{
			Score:  score,
			Member: idStr,
		})
	}

	// 清空旧的年榜
	if err := database.RDB.Del(ctx, yearKey).Err(); err != nil {
		log.Printf("⚠️ 清空旧年榜失败: %v", err)
	}

	// 存储到年榜ZSet
	if len(zsetMembers) > 0 {
		if err := database.RDB.ZAdd(ctx, yearKey, zsetMembers...).Err(); err != nil {
			return fmt.Errorf("存储年榜到Redis失败: %w", err)
		}
		// 设置过期时间为永久（年榜保留）
		// 不设置过期时间，永久保存
	}

	log.Printf("✅ 年榜生成完成，共 %d 篇文章", len(zsetMembers))
	return nil
}

// generateYearRankFromHot 从热门文章生成年榜（降级方案）
func (t *ArticleRankTask) generateYearRankFromHot(ctx context.Context) error {
	hotKey := "hot:articles:zset"
	yearKey := "rank:articles:year"

	// 从热门文章ZSet获取前50条
	articleIDStrs, err := database.RDB.ZRevRange(ctx, hotKey, 0, 49).Result()
	if err != nil {
		return fmt.Errorf("从热门文章ZSet获取数据失败: %w", err)
	}

	if len(articleIDStrs) == 0 {
		return nil
	}

	// 获取对应的分值
	zsetMembers := make([]*redis.Z, 0, len(articleIDStrs))
	for _, idStr := range articleIDStrs {
		score, err := database.RDB.ZScore(ctx, hotKey, idStr).Result()
		if err != nil {
			continue
		}
		zsetMembers = append(zsetMembers, &redis.Z{
			Score:  score,
			Member: idStr,
		})
	}

	// 清空旧的年榜
	database.RDB.Del(ctx, yearKey)

	// 存储到年榜ZSet
	if len(zsetMembers) > 0 {
		if err := database.RDB.ZAdd(ctx, yearKey, zsetMembers...).Err(); err != nil {
			return fmt.Errorf("存储年榜到Redis失败: %w", err)
		}
		// 年榜永久保存，不设置过期时间
	}

	return nil
}
