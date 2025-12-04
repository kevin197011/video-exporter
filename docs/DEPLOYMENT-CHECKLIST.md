# 部署检查清单

## 配置文件清单

在启动 Docker Compose 之前，请确保以下文件存在：

### ✅ 必需文件

- [x] `docker-compose.yml` - Docker Compose 配置
- [x] `config.yml` - Video Exporter 配置（从 config.example.yaml 复制并修改）
- [x] `prometheus.yml` - Prometheus 配置
- [x] `grafana-provisioning/dashboards/video-stream-dashboard.json` - Grafana 仪表板

### ✅ Grafana Provisioning 配置

- [x] `grafana-provisioning/datasources/prometheus.yml` - 数据源自动配置
- [x] `grafana-provisioning/dashboards/dashboard.yml` - 仪表板自动导入配置

### ✅ 辅助脚本

- [x] `start.sh` - 快速启动脚本
- [x] `stop.sh` - 快速停止脚本

## 部署步骤

### 1. 检查配置文件

```bash
# 确保 config.yml 存在且配置正确
ls -l config.yml

# 如果不存在，从示例复制并修改
cp config.example.yaml config.yml
vim config.yml
```

### 2. 验证 Docker Compose 配置

```bash
docker-compose config
```

### 3. 启动服务

```bash
# 方式一：使用启动脚本（推荐）
./start.sh

# 方式二：手动启动
docker-compose up -d
```

### 4. 检查服务状态

```bash
# 查看容器状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 5. 验证服务

访问以下地址确认服务正常：

- [ ] Video Exporter: http://localhost:8080
- [ ] Prometheus: http://localhost:9090/targets
- [ ] Grafana: http://localhost:3000

### 6. Grafana 配置验证

登录 Grafana (admin/admin) 后检查：

- [ ] 数据源已自动配置（Configuration → Data Sources）
- [ ] 仪表板已自动导入（Dashboards → Browse）
- [ ] 数据正常显示

## 网络指标验证

在 Prometheus 中执行以下查询，确认新指标可用：

```promql
# RTT 指标
video_stream_rtt_ms

# 丢包率指标
video_stream_packet_loss_ratio

# 网络抖动指标
video_stream_network_jitter_ms

# 重连次数指标
video_stream_reconnect_total
```

## 常见问题

### 端口冲突

如果端口已被占用，修改 `docker-compose.yml` 中的端口映射：

```yaml
ports:
  - "8081:8080"  # 改为其他端口
```

### 数据不显示

1. 检查 video-exporter 配置是否正确
2. 检查 Prometheus targets 状态
3. 检查时间范围设置

### Grafana 仪表板未自动导入

手动导入仪表板：
1. 登录 Grafana
2. 点击 "+" → "Import"
3. 上传 `deployments/grafana/grafana-provisioning/dashboards/video-stream-dashboard.json` 文件

## 文档参考

- [DOCKER-COMPOSE-README.md](./DOCKER-COMPOSE-README.md) - 详细部署文档
- [README.md](./README.md) - 项目说明
- [config.example.yaml](./config.example.yaml) - 配置示例

## 升级和维护

### 更新 Video Exporter

```bash
# 重新构建镜像
docker-compose build video-exporter

# 重启服务
docker-compose up -d video-exporter
```

### 备份数据

```bash
# 备份所有数据卷
docker-compose down
tar czf backup-$(date +%Y%m%d).tar.gz \
  config.yml \
  prometheus.yml \
  grafana-provisioning/
```

### 清理和重置

```bash
# 停止并删除所有容器和数据
docker-compose down -v

# 重新开始
./start.sh
```

---

**部署完成后，记得修改 Grafana 的默认密码！**

