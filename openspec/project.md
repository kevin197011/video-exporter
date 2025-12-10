# Project Context

## Purpose

Video Exporter 是一个基于 FFmpeg 的视频流监控导出系统，用于实时监控直播流的健康状况和质量指标。系统通过定期采样流数据，分析码率、帧率、分辨率、GOP 等质量指标，并通过 Prometheus 指标和 Grafana 仪表板提供可视化和告警能力。

## Tech Stack

- **语言**: Go 1.24+
- **流处理**: FFmpeg (通过 joy5 库解析 FLV 格式)
- **监控**: Prometheus (指标导出)
- **可视化**: Grafana (仪表板)
- **容器化**: Docker, Docker Compose
- **CI/CD**: GitHub Actions

## Project Conventions

### Code Style

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 包名使用小写，简短且有意义
- 导出函数/类型使用大写开头
- 注释使用中文，遵循 Go 文档注释规范

### Architecture Patterns

- **标准 Go 项目布局**:
  - `cmd/` - 可执行程序入口
  - `internal/` - 内部包（不对外暴露）
  - `deployments/` - 部署相关配置
  - `scripts/` - 工具脚本
  - `docs/` - 文档

- **模块化设计**:
  - `internal/config` - 配置管理
  - `internal/logger` - 日志系统
  - `internal/stream` - 流检查核心逻辑
  - `internal/scheduler` - 调度和并发控制
  - `internal/exporter` - Prometheus 指标导出

- **并发模型**:
  - 使用 goroutine 实现并发流检查
  - 使用 channel 进行通信
  - 使用 sync.RWMutex 保护共享状态

### Testing Strategy

- 单元测试覆盖核心业务逻辑
- 集成测试验证端到端流程
- 使用 Go 标准 testing 包
- 测试文件命名: `*_test.go`

### Git Workflow

- 遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范
- 提交类型: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
- 主分支: `main`
- 通过 Pull Request 合并代码

## Domain Context

### 核心概念

- **流 (Stream)**: 一个视频流地址，包含 URL、ID、项目标识
- **检查周期 (Check Interval)**: 定期检查流的间隔时间（默认 30 秒）
- **采样时长 (Sample Duration)**: 每次检查时采样的数据时长（默认 10 秒）
- **关键帧 (Keyframe/I-Frame)**: 视频编码中的独立帧，用于 GOP 分析
- **GOP (Group of Pictures)**: 关键帧之间的帧组，影响视频质量和延迟

### 质量指标

- **码率 (Bitrate)**: 数据传输速率，单位 bps
- **帧率 (Framerate)**: 每秒帧数，单位 fps
- **分辨率 (Resolution)**: 视频宽高，如 1920x1080
- **编码格式 (Codec)**: 如 H.264, H.265
- **可播放性 (Playable)**: 流是否可以被播放器正常播放
- **健康状态 (Healthy)**: 综合质量评估结果

### 网络指标

- **RTT (Round-Trip Time)**: 往返时间，单位毫秒
- **丢包率 (Packet Loss Ratio)**: 0.0-1.0，表示丢失的数据包比例
- **网络抖动 (Network Jitter)**: 包间隔的标准差，单位毫秒
- **重连次数 (Reconnect Count)**: 流断开后重新连接的次数

### Prometheus 标签

- `project`: 项目标识
- `id`: 流 ID（桌台ID）
- `name`: 流名称（自动生成）
- `url`: 流地址
- `service`: 服务标识（固定为 "video-exporter"）

## Important Constraints

- **FFmpeg 依赖**: 系统必须安装 FFmpeg，用于流解析
- **内存限制**: 每个流检查会占用一定内存，需要根据流数量合理配置
- **并发限制**: `max_concurrent` 控制同时检查的流数量，避免资源耗尽
- **网络超时**: 流检查可能因网络问题超时，需要合理的重试机制
- **采样时长**: 采样时长影响检查速度和指标准确性，需要平衡

## External Dependencies

- **FFmpeg**: 视频流处理（系统级依赖）
- **Prometheus**: 指标收集和存储
- **Grafana**: 数据可视化
- **joy5**: Go 语言 FLV 格式解析库
- **prometheus/client_golang**: Prometheus Go 客户端库

## Performance Characteristics

| 流数量 | 内存占用 | CPU占用 | 检查周期 |
|--------|----------|---------|----------|
| 1路    | ~10MB    | <1%     | 30s      |
| 10路   | ~30MB    | ~5%     | 30s      |
| 100路  | ~200MB   | ~20%    | 30s      |

## Configuration

配置文件: `config.yml`

主要参数:
- `check_interval`: 检查间隔（秒），建议 20-60
- `sample_duration`: 采样时长（秒），建议 5-15
- `min_keyframes`: 最小关键帧数，建议 2-5
- `max_concurrent`: 最大并发数，建议 100-1000
- `max_retries`: 重试次数，建议 3-5
- `listen_addr`: HTTP 监听地址，默认 ":8080"

## Deployment

- **本地运行**: `make run` 或 `go run ./cmd/video-exporter`
- **Docker**: `docker-compose up -d`
- **生产环境**: 建议使用 systemd 或 Kubernetes

## Monitoring Endpoints

- **Metrics**: `http://localhost:8080/metrics` (Prometheus 格式)
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (默认 admin/admin)
