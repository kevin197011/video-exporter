# Go 代码结构说明

本文档说明 Video Exporter 项目的 Go 代码组织方式。

## 目录结构

```
video-exporter/
├── cmd/                        # 程序入口目录
│   └── video-exporter/        # 主程序
│       └── main.go            # 入口文件
│
├── internal/                   # 内部包（不对外暴露）
│   ├── config/                # 配置管理模块
│   │   └── config.go          # 配置加载和结构体定义
│   │
│   ├── logger/                # 日志系统模块
│   │   └── logger.go          # 日志初始化和管理
│   │
│   ├── stream/                # 流检查核心模块
│   │   └── stream.go          # 流检查器实现
│   │
│   ├── scheduler/             # 任务调度模块
│   │   └── scheduler.go       # 并发调度和重试逻辑
│   │
│   └── exporter/              # Prometheus 导出器模块
│       └── exporter.go        # 指标定义和 HTTP 服务器
│
├── go.mod                      # Go 模块定义
└── go.sum                      # 依赖校验文件
```

## 模块说明

### 1. cmd/video-exporter

**用途**: 程序入口

**文件**: `main.go`

**职责**:
- 初始化各个模块
- 加载配置
- 启动服务
- 处理信号（优雅退出）

**示例**:
```go
package main

import (
    "video-exporter/internal/config"
    "video-exporter/internal/logger"
    "video-exporter/internal/scheduler"
    "video-exporter/internal/exporter"
)

func main() {
    logger.Init()
    cfg, _ := config.Load("config.yml")
    sched := scheduler.New(cfg)
    exp := exporter.New(sched)
    // ...
}
```

### 2. internal/config

**用途**: 配置管理

**导出类型**:
- `Config` - 主配置结构
- `ExporterConfig` - 导出器配置
- `StreamConfig` - 流配置

**导出函数**:
- `Load(filename string) (*Config, error)` - 加载配置文件
- `SetGlobal(cfg *Config)` - 设置全局配置
- `GetGlobal() *Config` - 获取全局配置

**用法**:
```go
cfg, err := config.Load("config.yml")
if err != nil {
    log.Fatal(err)
}
config.SetGlobal(cfg)
```

### 3. internal/logger

**用途**: 日志系统

**导出函数**:
- `Init()` - 初始化日志系统
- `Get() *slog.Logger` - 获取日志实例
- `SetLevel(level string)` - 设置日志级别

**用法**:
```go
logger.Init()
log := logger.Get()
log.Info("启动服务", "port", 8080)

logger.SetLevel("debug")
```

### 4. internal/stream

**用途**: 流检查核心逻辑

**导出类型**:
- `Checker` - 流检查器
- `Metrics` - 流指标数据

**导出函数**:
- `NewChecker(id, url, project string) *Checker` - 创建检查器
- `(c *Checker) Check(timeout time.Duration) error` - 执行检查
- `(c *Checker) MarkFailed()` - 标记失败
- `(c *Checker) GetMetrics() Metrics` - 获取指标
- `(c *Checker) ID() string` - 获取流ID

**用法**:
```go
checker := stream.NewChecker("stream-01", "http://example.com/stream.flv", "project1")
err := checker.Check(30 * time.Second)
if err != nil {
    checker.MarkFailed()
}
metrics := checker.GetMetrics()
```

### 5. internal/scheduler

**用途**: 任务调度和并发控制

**导出类型**:
- `Scheduler` - 调度器

**导出函数**:
- `New(cfg *config.Config) *Scheduler` - 创建调度器
- `(s *Scheduler) AddStream(id, url, project string)` - 添加流
- `(s *Scheduler) Start()` - 启动调度器
- `(s *Scheduler) Stop()` - 停止调度器
- `(s *Scheduler) GetAllMetrics() []stream.Metrics` - 获取所有指标

**用法**:
```go
sched := scheduler.New(cfg)
sched.AddStream("stream-01", "http://example.com/stream.flv", "project1")
go sched.Start()
// ...
sched.Stop()
```

### 6. internal/exporter

**用途**: Prometheus 指标导出

**导出类型**:
- `Exporter` - Prometheus 导出器

**导出函数**:
- `New(s *scheduler.Scheduler) *Exporter` - 创建导出器
- `(e *Exporter) StartHTTPServer(addr string) error` - 启动 HTTP 服务器
- `(e *Exporter) UpdateMetrics()` - 更新指标

**用法**:
```go
exp := exporter.New(sched)
go exp.StartHTTPServer(":8080")
```

## 依赖关系

```
main (cmd/video-exporter)
  ├─> config
  ├─> logger
  ├─> scheduler
  │     ├─> config
  │     ├─> logger
  │     └─> stream
  │           ├─> config
  │           └─> logger
  └─> exporter
        ├─> logger
        └─> scheduler
```

## 设计原则

### 1. 模块化

每个模块负责单一职责：
- `config` - 只负责配置管理
- `logger` - 只负责日志
- `stream` - 只负责流检查
- `scheduler` - 只负责调度
- `exporter` - 只负责指标导出

### 2. 封装性

使用 `internal/` 目录确保包不会被外部项目导入，保持 API 的灵活性。

### 3. 依赖方向

- 高层模块依赖低层模块
- 避免循环依赖
- `config` 和 `logger` 是基础模块，被其他模块依赖

### 4. 命名规范

- 包名使用小写单数形式（如 `config` 不是 `configs`）
- 导出函数首字母大写
- 私有函数首字母小写
- 构造函数命名为 `New` 或 `New{Type}`

## 构建和运行

### 构建

```bash
# 使用 Makefile
make build

# 直接使用 go build
go build -o video-exporter ./cmd/video-exporter

# 跨平台编译
make build-all
```

### 运行

```bash
# 使用 Makefile
make run

# 直接运行
go run ./cmd/video-exporter

# 运行编译后的二进制
./video-exporter
```

### 测试

```bash
# 运行测试
make test

# 或
go test ./...

# 查看覆盖率
go test -cover ./...
```

## 添加新模块

### 1. 创建模块目录

```bash
mkdir -p internal/newmodule
```

### 2. 创建模块文件

```go
// internal/newmodule/newmodule.go
package newmodule

import (
    "video-exporter/internal/logger"
)

type NewModule struct {
    log *slog.Logger
}

func New() *NewModule {
    return &NewModule{
        log: logger.Get(),
    }
}

func (n *NewModule) DoSomething() error {
    n.log.Info("doing something")
    return nil
}
```

### 3. 在 main.go 中使用

```go
import "video-exporter/internal/newmodule"

func main() {
    // ...
    mod := newmodule.New()
    mod.DoSomething()
}
```

## 代码风格

### 1. 错误处理

```go
// 推荐：返回详细错误信息
if err != nil {
    return fmt.Errorf("操作失败: %w", err)
}

// 避免：吞掉错误或返回通用错误
if err != nil {
    return err
}
```

### 2. 日志记录

```go
// 推荐：使用结构化日志
log.Info("处理请求", "method", "GET", "path", "/metrics")

// 避免：使用字符串拼接
log.Info(fmt.Sprintf("处理请求: %s %s", method, path))
```

### 3. 配置传递

```go
// 推荐：通过参数传递
func New(cfg *config.Config) *Module {
    return &Module{config: cfg}
}

// 避免：直接访问全局变量
func New() *Module {
    return &Module{config: globalConfig}
}
```

## 性能优化

### 1. 减少锁的使用

```go
// 好：读写锁分离
type Cache struct {
    mu    sync.RWMutex
    data  map[string]string
}

func (c *Cache) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}
```

### 2. 复用对象

```go
// 好：复用 HTTP 客户端
var globalHTTPClient *http.Client
var httpClientOnce sync.Once

func initHTTPClient() {
    httpClientOnce.Do(func() {
        globalHTTPClient = &http.Client{
            Transport: &http.Transport{
                MaxIdleConns: 100,
            },
        }
    })
}
```

### 3. 并发控制

```go
// 好：使用信号量控制并发数
semaphore := make(chan struct{}, maxConcurrent)
for _, item := range items {
    semaphore <- struct{}{}
    go func(i item) {
        defer func() { <-semaphore }()
        process(i)
    }(item)
}
```

## 参考资源

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go 项目标准布局](https://github.com/golang-standards/project-layout)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

