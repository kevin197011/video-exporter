# API 文档

Video Exporter 通过 HTTP 接口提供 Prometheus 格式的监控指标。

## Endpoints

### GET /

返回服务信息页面。

**响应**:
```html
<html>
<head><title>Video Stream Exporter</title></head>
<body>
<h1>Video Stream Exporter</h1>
<p><a href="/metrics">Metrics</a></p>
</body>
</html>
```

### GET /metrics

返回 Prometheus 格式的监控指标。

**响应格式**: `text/plain`

**示例请求**:
```bash
curl http://localhost:8080/metrics
```

**示例响应**:
```
# HELP video_stream_up Stream is up (1) or down (0)
# TYPE video_stream_up gauge
video_stream_up{id="stream-01",name="project1_example_stream-01",project="project1",url="https://example.com/stream.flv"} 1

# HELP video_stream_healthy Stream health status (1=healthy, 0=unhealthy)
# TYPE video_stream_healthy gauge
video_stream_healthy{id="stream-01",name="project1_example_stream-01",project="project1",url="https://example.com/stream.flv"} 1

# HELP video_stream_playable Stream is playable (1=yes, 0=no)
# TYPE video_stream_playable gauge
video_stream_playable{id="stream-01",name="project1_example_stream-01",project="project1",url="https://example.com/stream.flv"} 1

...
```

## 指标说明

所有指标都包含以下标签（Labels）：
- `project`: 项目名称
- `id`: 流 ID
- `name`: 流名称（自动生成）
- `url`: 流地址

### 基础状态指标

#### video_stream_up
- **类型**: Gauge
- **说明**: 流是否在线
- **值**: 1（在线）/ 0（离线）

#### video_stream_healthy
- **类型**: Gauge
- **说明**: 流健康状态
- **值**: 1（健康）/ 0（不健康）

#### video_stream_playable
- **类型**: Gauge
- **说明**: 流是否可播放
- **值**: 1（可播放）/ 0（不可播放）

### 数据包统计指标

#### video_stream_total_packets
- **类型**: Gauge
- **说明**: 总数据包数
- **单位**: 个

#### video_stream_video_packets
- **类型**: Gauge
- **说明**: 视频数据包数
- **单位**: 个

#### video_stream_audio_packets
- **类型**: Gauge
- **说明**: 音频数据包数
- **单位**: 个

#### video_stream_keyframes
- **类型**: Gauge
- **说明**: 关键帧数量
- **单位**: 个

### 码率指标

#### video_stream_bitrate_bps
- **类型**: Gauge
- **说明**: 当前码率
- **单位**: bits per second (bps)

#### video_stream_avg_bitrate_bps
- **类型**: Gauge
- **说明**: 平均码率
- **单位**: bits per second (bps)

### 视频质量指标

#### video_stream_framerate
- **类型**: Gauge
- **说明**: 帧率
- **单位**: fps (frames per second)

#### video_stream_response_ms
- **类型**: Gauge
- **说明**: HTTP-FLV 请求响应时间
- **单位**: 毫秒 (ms)

#### video_stream_gop_size
- **类型**: Gauge
- **说明**: GOP 大小（关键帧间隔）
- **单位**: 帧数

#### video_stream_quality_score
- **类型**: Gauge
- **说明**: 质量评分
- **值**:
  - 2 = good（高质量）
  - 1 = fair（中等质量）
  - 0 = poor（低质量）

#### video_stream_stability_score
- **类型**: Gauge
- **说明**: 码率稳定性评分
- **值**:
  - 2 = stable（稳定）
  - 1 = moderate（中等）
  - 0 = unstable（不稳定）

### 网络稳定性指标

#### video_stream_rtt_ms
- **类型**: Gauge
- **说明**: RTT 往返时间（使用 HTTP 响应时间作为近似值）
- **单位**: 毫秒 (ms)

#### video_stream_packet_loss_ratio
- **类型**: Gauge
- **说明**: 丢包率
- **值范围**: 0.0 - 1.0 (0% - 100%)
- **计算方法**: 基于视频包 DTS 时间戳连续性估算

#### video_stream_network_jitter_ms
- **类型**: Gauge
- **说明**: 网络抖动（包间隔时间的标准差）
- **单位**: 毫秒 (ms)

#### video_stream_reconnect_total
- **类型**: Counter
- **说明**: 重连总次数（累积值）
- **单位**: 次

## PromQL 查询示例

### 基础查询

```promql
# 查询所有在线的流
video_stream_up == 1

# 查询特定项目的流
video_stream_up{project="project1"}

# 查询不健康的流
video_stream_healthy == 0
```

### 码率查询

```promql
# 查询低于 500kbps 的流
video_stream_bitrate_bps < 500000

# 查询码率在 1Mbps - 3Mbps 之间的流
video_stream_bitrate_bps >= 1000000 and video_stream_bitrate_bps <= 3000000

# 计算码率（kbps）
video_stream_bitrate_bps / 1000
```

### 质量查询

```promql
# 查询低质量流
video_stream_quality_score == 0

# 查询码率不稳定的流
video_stream_stability_score == 0

# 查询帧率低于 20fps 的流
video_stream_framerate < 20
```

### 网络稳定性查询

```promql
# 查询 RTT 过高的流（> 1 秒）
video_stream_rtt_ms > 1000

# 查询丢包率高于 5% 的流
video_stream_packet_loss_ratio > 0.05

# 查询网络抖动严重的流（> 50ms）
video_stream_network_jitter_ms > 50

# 查询重连次数
rate(video_stream_reconnect_total[5m])
```

### 聚合查询

```promql
# 统计在线流数量
count(video_stream_up == 1)

# 统计每个项目的流数量
count by (project) (video_stream_up)

# 计算平均码率
avg(video_stream_bitrate_bps)

# 计算平均帧率
avg(video_stream_framerate)
```

### 趋势查询

```promql
# 过去 5 分钟码率变化
rate(video_stream_bitrate_bps[5m])

# 过去 1 小时的平均码率
avg_over_time(video_stream_bitrate_bps[1h])

# 重连速率（每分钟）
rate(video_stream_reconnect_total[1m]) * 60
```

## 告警规则示例

```yaml
groups:
  - name: video_stream_alerts
    interval: 30s
    rules:
      # 流离线告警
      - alert: StreamDown
        expr: video_stream_up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "视频流离线"
          description: "{{ $labels.project }}/{{ $labels.id }} 已离线超过 1 分钟"

      # 低码率告警
      - alert: LowBitrate
        expr: video_stream_bitrate_bps < 500000
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "码率过低"
          description: "{{ $labels.project }}/{{ $labels.id }} 码率低于 500kbps"

      # 响应过慢告警
      - alert: SlowResponse
        expr: video_stream_response_ms > 2000
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "响应时间过长"
          description: "{{ $labels.project }}/{{ $labels.id }} 响应时间超过 2 秒"

      # 高丢包率告警
      - alert: HighPacketLoss
        expr: video_stream_packet_loss_ratio > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "丢包率过高"
          description: "{{ $labels.project }}/{{ $labels.id }} 丢包率超过 10%"

      # 网络抖动告警
      - alert: HighNetworkJitter
        expr: video_stream_network_jitter_ms > 100
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "网络抖动严重"
          description: "{{ $labels.project }}/{{ $labels.id }} 网络抖动超过 100ms"

      # 频繁重连告警
      - alert: FrequentReconnects
        expr: rate(video_stream_reconnect_total[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "频繁重连"
          description: "{{ $labels.project }}/{{ $labels.id }} 每分钟重连次数超过 0.1 次"
```

## 性能考虑

### 抓取间隔

建议的 Prometheus 抓取间隔：
- 生产环境: 15-30 秒
- 开发环境: 5-15 秒

配置示例：
```yaml
scrape_configs:
  - job_name: 'video-exporter'
    scrape_interval: 30s
    scrape_timeout: 10s
    static_configs:
      - targets: ['localhost:8080']
```

### 资源占用

- 每个流的指标数量: ~20 个
- 每次抓取的数据量: ~2KB（单流）
- 1000 个流: ~2MB/次抓取

## 故障排查

### 指标不更新

1. 检查流配置是否正确
2. 查看 exporter 日志
3. 验证流地址是否可访问

### 指标值异常

1. 检查 `video_stream_up` 指标
2. 查看 `video_stream_healthy` 状态
3. 确认采样参数配置合理

### 性能问题

1. 调整 `check_interval` 参数
2. 限制并发数 `max_concurrent`
3. 优化 Prometheus 抓取频率

