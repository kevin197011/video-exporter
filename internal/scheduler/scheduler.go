package scheduler

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"video-exporter/internal/config"
	"video-exporter/internal/logger"
	"video-exporter/internal/stream"
)

// Scheduler 调度器
type Scheduler struct {
	checkers map[string]*stream.Checker
	config   *config.Config
	mu       sync.RWMutex
	stopChan chan struct{}
	log      *slog.Logger
}

// New 创建调度器
func New(cfg *config.Config) *Scheduler {
	return &Scheduler{
		checkers: make(map[string]*stream.Checker),
		config:   cfg,
		stopChan: make(chan struct{}),
		log:      logger.Get(),
	}
}

// AddStream 添加流
func (s *Scheduler) AddStream(id, url, project string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s::%s", project, url)
	checker := stream.NewChecker(id, url, project)
	s.checkers[key] = checker

	s.log.Info("添加流", "流ID", id, "URL", url, "项目", project)
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.log.Info("启动调度器",
		"流数量", len(s.checkers),
		"检查间隔秒", s.config.Exporter.CheckInterval,
		"最大并发", s.config.Exporter.MaxConcurrent,
		"最大重试", s.config.Exporter.MaxRetries)

	// 立即执行第一次检查（在后台）
	go s.runCheckCycle()

	// 定时执行检查
	ticker := time.NewTicker(time.Duration(s.config.Exporter.CheckInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			s.log.Info("调度器已停止")
			return
		case <-ticker.C:
			go s.runCheckCycle()
		}
	}
}

// runCheckCycle 执行一轮检查
func (s *Scheduler) runCheckCycle() {
	s.mu.RLock()
	checkers := make([]*stream.Checker, 0, len(s.checkers))
	for _, checker := range s.checkers {
		checkers = append(checkers, checker)
	}
	s.mu.RUnlock()

	s.log.Info("开始检查周期", "流数量", len(checkers))

	// 重置所有流的周期指标（重连次数等）
	for _, checker := range checkers {
		checker.ResetCycleMetrics()
	}

	// 使用信号量控制并发
	semaphore := make(chan struct{}, s.config.Exporter.MaxConcurrent)
	var wg sync.WaitGroup

	for _, checker := range checkers {
		wg.Add(1)
		go func(c *stream.Checker) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 执行检查，带重试
			s.checkWithRetry(c)
		}(checker)
	}

	wg.Wait()
	s.log.Info("检查周期完成")
}

// checkWithRetry 带重试的检查
func (s *Scheduler) checkWithRetry(checker *stream.Checker) {
	// 超时时间：采样时间(10秒) + 网络缓冲(5秒)
	timeout := 15 * time.Second

	// 如果检查间隔很长，可以给更多时间
	if s.config.Exporter.CheckInterval > 20 {
		timeout = time.Duration(s.config.Exporter.CheckInterval-5) * time.Second
	}

	var lastErr error
	for attempt := 0; attempt <= s.config.Exporter.MaxRetries; attempt++ {
		if attempt > 0 {
			// 重试前等待
			retryDelay := time.Duration(attempt*2) * time.Second
			s.log.Info("等待重试", "流ID", checker.ID(), "尝试次数", attempt, "延迟秒", retryDelay.Seconds())
			time.Sleep(retryDelay)
		}

		err := checker.Check(timeout)
		if err == nil {
			// 成功
			return
		}

		lastErr = err
		s.log.Error("检查失败", "流ID", checker.ID(), "尝试次数", attempt+1, "错误", err)
	}

	// 所有重试都失败
	checker.MarkFailed()
	s.log.Error("达到最大重试次数", "流ID", checker.ID(), "最后错误", lastErr)
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	close(s.stopChan)
}

// GetAllMetrics 获取所有流的指标
func (s *Scheduler) GetAllMetrics() []stream.Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := make([]stream.Metrics, 0, len(s.checkers))
	for _, checker := range s.checkers {
		metrics = append(metrics, checker.GetMetrics())
	}

	return metrics
}
