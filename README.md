# Video Exporter

基于 FFmpeg 的视频流监控导出系统，用于实时监控直播流的健康状况和质量指标。

## 功能特性

### 核心功能
- ✅ 实时流监控（多协程并发）
- ✅ 深度质量分析（码率、帧率、分辨率、GOP等）
- ✅ 健康评估系统（可播放性、质量等级）
- ✅ 网络稳定性监控（RTT、丢包率、抖动、重连）
- ✅ 延迟分析（流延迟计算）
- ✅ 自动重连机制
- ✅ 支持多种流格式（FLV、RTMP、HLS、RTSP等）

### 监控与可视化
- ✅ Prometheus 指标导出
- ✅ Grafana 仪表板（自动配置）
- ✅ Docker Compose 一键部署
- ✅ 结构化日志输出

## 快速开始

### 1. 安装 FFmpeg

**macOS:**
```bash
brew install ffmpeg pkg-config
```

**Ubuntu/Debian:**
```bash
sudo apt-get install -y libavcodec-dev libavformat-dev libavutil-dev pkg-config
```

### 2. 配置流地址

编辑 `config.yml`：

```yaml
exporter:
  check_interval: 30    # 检查间隔（秒）
  max_concurrent: 1000  # 最大并发数
  max_retries: 3        # 最大重试次数
  listen_addr: "8080"   # Prometheus 监听端口

streams:
  project1:  # 项目名称（用于 Prometheus 标签）
    - url: https://example.com/live/stream.flv
      id: stream-01
```

### 3. 运行

```bash
# 方式1: 使用 Makefile（推荐）
make run

# 方式2: 直接运行
go run ./cmd/video-exporter

# 方式3: 编译后运行
make build
./video-exporter
```

## 项目结构

```
video-exporter/
├── main.go                       # 程序入口
├── config.go                     # 配置加载/结构体
├── logger.go                     # 日志系统
├── exporter.go                   # Prometheus 指标导出
├── scheduler.go                  # 调度与并发检查
├── stream.go                     # 核心流检查逻辑
├── config.yml                    # 配置文件
├── config.example.yaml           # 配置文件示例
├── deployments/                  # 部署配置
│   └── grafana/                  # Grafana 配置
│       └── grafana-provisioning/ # 自动配置
│           ├── dashboards/       # 仪表板
│           │   └── video-stream-dashboard.json
│           └── datasources/      # 数据源
│               └── prometheus.yml
├── docker-compose.yml            # Docker Compose 编排配置
├── Dockerfile                    # Docker 镜像构建
├── prometheus.yml                # Prometheus 配置
├── scripts/                      # 脚本工具
│   ├── start.sh                 # 启动脚本
│   ├── stop.sh                  # 停止脚本
│   └── logs.sh                  # 日志查看脚本
├── docs/                         # 文档
│   ├── API.md                   # API 文档
│   ├── DOCKER-COMPOSE-README.md # Docker Compose 使用说明
│   └── DEPLOYMENT-CHECKLIST.md  # 部署检查清单
├── CONTRIBUTING.md               # 贡献指南
├── CHANGELOG.md                  # 更新日志
├── Makefile                      # 常用命令
├── go.mod
└── go.sum
```

## 监控输出

### 控制台输出
```
检查 #001 stream-01 stream-01 (https://...)
可播放: true | 质量: good | 响应: 150ms
视频包: 1234 | 关键帧: 45
码率: 2500.5kbps (平均: 2480.3kbps) | 稳定性: stable
帧率: 25.0fps | 分辨率: 1920x1080
编码: H.264 | GOP: 75帧
```

### Prometheus 指标
访问 `http://localhost:8080/metrics` 查看所有指标：
```
video_stream_up{project="project1",id="stream-01",url="https://..."} 1
video_stream_bitrate_bps{project="project1",id="stream-01",url="https://..."} 753000.0
video_stream_framerate{project="project1",id="stream-01",url="https://..."} 42.0
video_stream_response_ms{project="project1",id="stream-01",url="https://..."} 150
```

## 监控指标

### 基础指标
- 总包数、视频包数、音频包数
- 关键帧数量
- 流状态（在线/离线）

### 质量指标
- **码率**: 实时码率、平均码率、码率稳定性
- **帧率**: 实时帧率计算
- **分辨率**: 视频分辨率识别
- **GOP**: 关键帧间隔分析
- **编码**: 视频编码格式（H.264/H.265等）

### 网络指标 🆕
- **RTT**: 往返时间（毫秒）
- **丢包率**: 0.0-1.0（0%-100%）
- **网络抖动**: 包间隔标准差（毫秒）
- **重连次数**: 累积重连统计

### 健康评估
- 可播放性判断
- 质量等级（good/fair/poor）
- 响应时长（FLV HTTP 请求响应时间，单位：ms）
- 异常检测

> 💡 详细的 API 文档请查看 [docs/API.md](docs/API.md)

## 配置说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| check_interval | 健康检查间隔（秒） | 30 |
| max_concurrent | 最大并发监控数 | 1000 |
| max_retries | 连接失败最大重试次数 | 3 |
| listen_addr | Prometheus 监听端口 | 8080 |

## 支持的流格式

- FLV / HTTP-FLV
- RTMP / RTMPS
- HLS (m3u8)
- RTSP
- 其他 FFmpeg 支持的格式

## 性能

| 流数量 | 内存占用 | CPU占用 |
|--------|----------|---------|
| 1路    | ~10MB    | <1%     |
| 10路   | ~30MB    | ~5%     |
| 100路  | ~200MB   | ~20%    |


## 编译

```bash
# 使用 Makefile（推荐）
make build              # 本地编译
make build-all          # 所有平台

# 手动编译
go build -o video-exporter ./cmd/video-exporter

# 跨平台编译
GOOS=linux GOARCH=amd64 go build -o video-exporter-linux ./cmd/video-exporter
GOOS=windows GOARCH=amd64 go build -o video-exporter.exe ./cmd/video-exporter
GOOS=darwin GOARCH=amd64 go build -o video-exporter-mac ./cmd/video-exporter
```

## 部署

### Docker Compose（推荐）

一键启动 Video Exporter + Prometheus + Grafana：

```bash
# 启动
./scripts/start.sh

# 查看日志
./scripts/logs.sh

# 停止
./scripts/stop.sh
```

访问：
- Video Exporter: http://localhost:8080
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)

> 📖 详细部署文档：
> - [Docker Compose 使用说明](docs/DOCKER-COMPOSE-README.md)
> - [部署检查清单](docs/DEPLOYMENT-CHECKLIST.md)

### 后台运行

```bash
nohup ./video-exporter > monitor.log 2>&1 &
```

### Systemd 服务

```ini
[Unit]
Description=Video Exporter
After=network.target

[Service]
Type=simple
User=nobody
WorkingDirectory=/opt/video-exporter
ExecStart=/opt/video-exporter/video-exporter
Restart=always

[Install]
WantedBy=multi-user.target
```

## Prometheus 集成

### 访问指标
```bash
# 查看所有指标
curl http://localhost:8080/metrics

# 在浏览器中访问
http://localhost:8080/metrics
```

### Prometheus 配置
```yaml
scrape_configs:
  - job_name: 'video-exporter'
    static_configs:
      - targets: ['localhost:8080']
    scrape_interval: 15s
```

### 告警示例
```yaml
# 流离线告警
- alert: StreamDown
  expr: video_stream_up == 0
  for: 1m

# 低码率告警
- alert: LowBitrate
  expr: video_stream_bitrate_kbps < 500
  for: 2m

# 响应过慢告警（FLV HTTP 请求响应时间）
- alert: SlowResponse
  expr: video_stream_response_ms > 2000
  for: 1m
```


## 常见问题

### Q: 连接失败
A: 检查流地址是否正确，网络是否可达

### Q: 码率为0
A: 等待1-2个检查周期，让系统收集足够数据

### Q: 如何查看 Prometheus 指标
A: 访问 http://localhost:8080/metrics

### Q: 响应时间显示 N/A
A: 需要成功完成 HTTP 连接才会产生响应时间

## 文档

- [API 文档](docs/API.md) - Prometheus 指标和 API 说明
- [贡献指南](CONTRIBUTING.md) - 如何贡献代码
- [更新日志](CHANGELOG.md) - 版本更新记录
- [Docker Compose 使用](docs/DOCKER-COMPOSE-README.md) - 容器化部署
- [部署检查清单](docs/DEPLOYMENT-CHECKLIST.md) - 部署前检查

## 贡献

欢迎贡献代码！请查看 [贡献指南](CONTRIBUTING.md)。

提交代码前请确保：
- 遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范
- 代码通过 `go fmt` 和 `go vet` 检查
- 添加必要的测试和文档

## 许可证

MIT License
