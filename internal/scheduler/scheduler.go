package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

// Task 定时任务接口
type Task interface {
	// Run 执行任务
	Run(ctx context.Context) error
	// Name 任务名称
	Name() string
}

// TaskInfo 任务信息
type TaskInfo struct {
	Task     Task
	Interval time.Duration
	ticker   *time.Ticker
	cancel   context.CancelFunc
}

// Scheduler 调度器
type Scheduler struct {
	tasks map[string]*TaskInfo
	mu    sync.RWMutex
}

// NewScheduler 创建新的调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make(map[string]*TaskInfo),
	}
}

// RegisterTask 注册任务
func (s *Scheduler) RegisterTask(name string, task Task, interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[name] = &TaskInfo{
		Task:     task,
		Interval: interval,
	}

	log.Printf("✅ 注册任务: %s (间隔: %s)", name, interval)
}

// Start 启动所有任务
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, info := range s.tasks {
		// 立即执行一次
		go func(taskName string, taskInfo *TaskInfo) {
			log.Printf("🚀 首次执行任务: %s", taskName)
			ctx := context.Background()
			if err := taskInfo.Task.Run(ctx); err != nil {
				log.Printf("❌ 任务 %s 执行失败: %v", taskName, err)
			} else {
				log.Printf("✅ 任务 %s 执行成功", taskName)
			}
		}(name, info)

		// 创建定时器
		info.ticker = time.NewTicker(info.Interval)
		ctx, cancel := context.WithCancel(context.Background())
		info.cancel = cancel

		// 启动定时任务
		go s.runTask(ctx, name, info)
	}

	log.Println("✅ 所有任务已启动")
}

// runTask 运行单个任务
func (s *Scheduler) runTask(ctx context.Context, name string, info *TaskInfo) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("⏹️  任务 %s 已停止", name)
			return
		case <-info.ticker.C:
			log.Printf("🔄 执行定时任务: %s", name)
			startTime := time.Now()
			
			if err := info.Task.Run(ctx); err != nil {
				log.Printf("❌ 任务 %s 执行失败: %v", name, err)
			} else {
				duration := time.Since(startTime)
				log.Printf("✅ 任务 %s 执行成功 (耗时: %s)", name, duration)
			}
		}
	}
}

// Stop 停止所有任务
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, info := range s.tasks {
		if info.ticker != nil {
			info.ticker.Stop()
		}
		if info.cancel != nil {
			info.cancel()
		}
		log.Printf("⏹️  停止任务: %s", name)
	}

	log.Println("✅ 所有任务已停止")
}

// GetTasks 获取所有任务信息
func (s *Scheduler) GetTasks() map[string]*TaskInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make(map[string]*TaskInfo)
	for name, info := range s.tasks {
		tasks[name] = info
	}
	return tasks
}

