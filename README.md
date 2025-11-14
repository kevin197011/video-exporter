# Video Exporter

基于 FFmpeg 的视频流监控导出系统，用于实时监控直播流的健康状况和质量指标。

## 功能特性

- ✅ 实时流监控（多协程并发）
- ✅ 深度质量分析（码率、帧率、分辨率、GOP等）
- ✅ 健康评估系统（可播放性、质量等级）
- ✅ 延迟分析（流延迟计算）
- ✅ 自动重连机制
- ✅ 支持多种流格式（FLV、RTMP、HLS、RTSP等）
- ✅ Prometheus 指标导出
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
# 方式1: 直接运行
go run *.go

# 方式2: 编译后运行
go build -o video-exporter
./video-exporter
```

## 项目结构

```
video-exporter/
├── main.go                 # 程序入口
├── config.go               # 配置加载/结构体
├── logger.go               # 日志系统
├── exporter.go             # Prometheus 指标导出
├── scheduler.go            # 调度与并发检查
├── stream.go               # 核心流检查逻辑
├── config.yml              # 配置文件（挂载到容器 /app/config.yml）
├── Dockerfile              # 多阶段构建镜像
├── docker-compose.yml      # 本地/服务器编排与配置挂载
├── grafana-dashboard.json  # Grafana 仪表盘（按项目过滤）
├── Makefile                # 常用命令
├── run.sh                  # 运行脚本
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
- 数据包接收时间

### 深度指标
- **码率**: 实时码率、平均码率、码率稳定性
- **帧率**: 实时帧率计算
- **分辨率**: 视频分辨率识别
- **GOP**: 关键帧间隔分析
- **编码**: 视频编码格式（H.264/H.265等）

### 健康评估
- 可播放性判断
- 质量等级（good/fair/poor）
- 响应时长（FLV HTTP 请求响应时间，单位：ms）
- 异常检测

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
# 本地编译
go build -o video-exporter

# Linux
GOOS=linux GOARCH=amd64 go build -o video-exporter-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o video-exporter.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o video-exporter-mac
```

## 部署

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

## 许可证

MIT License
