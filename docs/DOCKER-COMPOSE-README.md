# Docker Compose 部署说明

## 服务架构

此 Docker Compose 配置包含三个服务：

1. **video-exporter** - 视频流监控导出器
2. **prometheus** - 指标采集和存储
3. **grafana** - 可视化仪表板

## 快速启动

### 1. 启动所有服务

```bash
# 使用启动脚本（推荐）
./start.sh

# 或手动启动
docker-compose up -d
```

### 2. 查看服务状态

```bash
docker-compose ps
```

### 3. 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f video-exporter
docker-compose logs -f prometheus
docker-compose logs -f grafana
```

## 访问地址

- **Video Exporter**: http://localhost:8080
  - Metrics endpoint: http://localhost:8080/metrics

- **Prometheus**: http://localhost:9090
  - 可以查看采集到的指标和执行 PromQL 查询

- **Grafana**: http://localhost:3000
  - 默认用户名: `admin`
  - 默认密码: `admin`
  - 首次登录后建议修改密码

## 配置说明

### Video Exporter

- 配置文件: `config.yml`
- 端口: `8080`
- 时区: `Asia/Shanghai`

### Prometheus

- 配置文件: `prometheus.yml`
- 端口: `9090`
- 数据保留时间: 30天
- 抓取间隔: 30秒（针对 video-exporter）
- 数据存储: Docker volume `prometheus_data`

### Grafana

- 端口: `3000`
- 数据存储: Docker volume `grafana_data`
- 自动配置 Prometheus 数据源
- 自动导入 `video-stream-dashboard.json` 仪表板

## 管理操作

### 停止服务

```bash
./stop.sh
# 或
docker-compose stop
```

### 重启服务

```bash
docker-compose restart
```

### 停止并删除容器

```bash
docker-compose down
```

### 停止并删除容器及数据卷

```bash
docker-compose down -v
```

## 网络指标说明

新增的网络稳定性指标：

- **video_stream_rtt_ms**: RTT往返时间（毫秒）
- **video_stream_packet_loss_ratio**: 丢包率（0.0-1.0）
- **video_stream_network_jitter_ms**: 网络抖动（毫秒）
- **video_stream_reconnect_total**: 重连总次数（累积）

这些指标会自动显示在 Grafana 仪表板中。
