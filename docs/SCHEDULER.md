# 定时任务调度器

## 概述

本项目使用独立的定时任务调度器来处理后台统计和数据更新任务，与主服务解耦，提高系统稳定性和可扩展性。

## 架构设计

```
┌─────────────────────────────────────┐
│  用户服务 (8081)                     │
│  管理服务 (8083)                     │
│  定时任务调度器 (独立进程)           │
├─────────────────────────────────────┤
│  MySQL (3306)                       │
│  Redis (6379)                       │
└─────────────────────────────────────┘
```

### 特点

- ✅ **抽象任务框架** - 易于扩展新任务
- ✅ **独立进程运行** - 不影响API服务性能
- ✅ **自动重试** - 任务失败后下个周期自动重试
- ✅ **首次立即执行** - 启动后立即执行一次
- ✅ **Redis缓存** - 计算结果存储在Redis，API直接读取

## 已实现的任务

### 1. 热门文章统计 (hot_articles)

**功能：** 根据多维度指标计算文章热度排名

**计算公式：**
```
热度得分 = log(浏览量+1) × 50% 
         + log(评论数+1) × 20%
         + log(点赞数+1) × 15%
         + log(收藏数+1) × 15%
```

**执行周期：** 每10分钟

**存储位置：** Redis `hot:articles`（JSON数组，存储前20篇热门文章ID）

**过期时间：** 20分钟

## 启动方式

### 开发模式
```bash
go run cmd/scheduler/main.go
```

### 编译并运行
```bash
# 编译
go build -o bin/scheduler cmd/scheduler/main.go

# 运行
./bin/scheduler
```

### Docker方式（未来支持）
```bash
docker-compose up -d scheduler
```

## 完整启动流程

```bash
# 1. 启动数据库
docker-compose up -d mysql redis

# 2. 数据库迁移（会在服务启动时自动执行）
# 初始化基础数据（可选）
mysql -h localhost -u inkspace -pinkspace123 inkspace < scripts/init.sql

# 3. 启动服务（4个终端）
go run cmd/server/main.go       # 终端1: 用户服务
go run cmd/admin/main.go        # 终端2: 管理服务
go run cmd/scheduler/main.go    # 终端3: 定时任务
cd web/blog && pnpm dev         # 终端4: 前端
```

## 如何添加新任务

### 1. 创建任务文件

在 `internal/scheduler/` 目录创建新任务：

```go
// internal/scheduler/your_task.go
package scheduler

import (
	"context"
	"log"
)

type YourTask struct{}

func NewYourTask() *YourTask {
	return &YourTask{}
}

func (t *YourTask) Name() string {
	return "你的任务名称"
}

func (t *YourTask) Run(ctx context.Context) error {
	log.Println("执行你的任务...")
	
	// 实现你的任务逻辑
	// ...
	
	log.Println("任务完成")
	return nil
}
```

### 2. 注册任务

在 `cmd/scheduler/main.go` 中注册：

```go
// 注册任务
sched.RegisterTask("your_task", scheduler.NewYourTask(), 30*time.Minute)
```

### 3. 重启调度器

```bash
# 停止旧的调度器（Ctrl+C）
# 启动新的调度器
go run cmd/scheduler/main.go
```

## 任务管理

### 查看任务状态

启动时会显示所有已注册的任务：

```
========================================
✅ 定时任务调度器启动成功
已注册的任务:
  - hot_articles (间隔: 10m0s)
  - daily_stats (间隔: 24h0m0s)
========================================
```

### 任务日志

每次任务执行都会记录日志：

```
🔄 执行定时任务: hot_articles
开始计算热门文章...
✅ 热门文章计算完成，共 25 篇文章，已存储前 20 篇到Redis
📊 Top 5 热门文章: [4 3 2 5 1]
✅ 任务 hot_articles 执行成功 (耗时: 125ms)
```

## Redis数据结构

### hot:articles

**类型：** String (JSON)

**格式：**
```json
[4, 3, 2, 5, 1, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]
```

**说明：** 按热度降序排列的文章ID数组

**过期时间：** 20分钟

**查看命令：**
```bash
redis-cli
> GET hot:articles
> TTL hot:articles  # 查看剩余时间
```

## API使用

### 获取热门文章

```http
GET /api/articles/hot?limit=6
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 4,
      "title": "热门文章1",
      "view_count": 150,
      "like_count": 25,
      ...
    }
  ]
}
```

### 降级策略

如果Redis中没有数据，API会自动降级返回最新文章。

## 性能优化

### 为什么使用定时任务？

1. **减少API延迟** - 计算在后台完成，API只需读取Redis
2. **降低数据库压力** - 避免每次API调用都查询统计数据
3. **保证数据一致性** - 所有用户看到相同的热门排名
4. **易于扩展** - 可以添加更多复杂的统计任务

### 性能对比

| 方案 | API响应时间 | 数据库查询 | 并发支持 |
|------|------------|----------|---------|
| 实时计算 | 200-500ms | 复杂JOIN | 低 |
| Redis缓存 | 5-10ms | 简单查询 | 高 |

## 监控和告警

### 健康检查

可以通过日志监控任务执行情况：

```bash
# 查看调度器日志
tail -f scheduler.log

# 或使用 grep 过滤
tail -f scheduler.log | grep "❌"  # 只看错误
tail -f scheduler.log | grep "✅"  # 只看成功
```

### 常见问题

**Q: 调度器启动失败？**
- 检查MySQL和Redis是否正常运行
- 检查配置文件 `config/config.yaml`
- 查看错误日志

**Q: 热门文章不更新？**
- 检查调度器是否在运行
- 查看Redis中是否有数据：`redis-cli GET hot:articles`
- 检查任务执行日志

**Q: 如何修改更新频率？**
- 修改 `cmd/scheduler/main.go` 中的时间间隔
- 重启调度器

## 未来扩展

可以添加的任务：

- **每日统计任务** - 统计网站访问量、新增用户等
- **缓存清理任务** - 定期清理过期缓存
- **邮件通知任务** - 发送订阅邮件
- **数据备份任务** - 定期备份数据库
- **热门标签统计** - 计算热门标签
- **用户活跃度统计** - 计算用户活跃度排名

## 最佳实践

1. **任务执行时间** - 避免在高峰期执行重量级任务
2. **错误处理** - 任务失败不应影响其他任务
3. **日志记录** - 记录关键步骤和错误信息
4. **幂等性** - 任务多次执行应该产生相同结果
5. **超时控制** - 设置合理的超时时间

---

**文档维护**: 添加新任务时需同步更新此文档

