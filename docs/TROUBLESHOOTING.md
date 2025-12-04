# 故障排查指南

本文档提供常见问题的排查和解决方案。

## Counter 类型指标无数据

### 问题描述

查询 `video_stream_reconnect_total{id="D001"}` 没有数据。

### 原因分析

`video_stream_reconnect_total` 是 **Counter** 类型指标：

**Counter 的特性**：
- ✅ 只在值 > 0 时才创建和显示
- ❌ 值为 0 时，指标不存在
- 🎯 这是 Prometheus Counter 的标准行为

**三种可能的情况**：

#### 1. 流在线，但从未重连（正常）
```promql
# 查询流状态
video_stream_up{id="D001"}  # 返回 1

# 查询重连次数
video_stream_reconnect_total{id="D001"}  # 无数据（因为 counter = 0）
```

**解决方案**：
- ✅ 已修复：查询添加了 `or vector(0)`
- 现在会显示 0 而不是无数据

#### 2. 流离线或不存在
```promql
# 查询流状态
video_stream_up{id="D001"}  # 无数据

# 所有相关指标都无数据
video_stream_reconnect_total{id="D001"}  # 无数据
```

**排查方法**：
```bash
# 1. 检查配置
grep "D001" config.yml

# 2. 检查日志
docker logs video-exporter | grep "D001"

# 3. 查看 Prometheus targets
curl http://localhost:9090/api/v1/targets
```

#### 3. 流最近才上线，Prometheus 还未抓取
**解决方案**：等待下一个抓取周期（默认30秒）

### 修复后的查询

**Grafana 仪表板查询**（已自动更新）：
```promql
# 总重连次数统计卡片
sum(video_stream_reconnect_total{project=~"$project", id=~"$id", name=~"$name"}) or vector(0)

# 重连次数趋势图
video_stream_reconnect_total{project=~"$project", id=~"$id", name=~"$name"} or vector(0)
```

**效果**：
- 有重连记录的流：显示实际次数
- 在线但无重连的流：显示 0
- 离线的流：不显示（符合预期）

## 验证方法

### 1. 检查当前有哪些流有数据

```bash
# 查看所有在线流
curl -s 'http://localhost:9090/api/v1/query?query=video_stream_up' | \
  python3 -c "import sys, json; data=json.load(sys.stdin); \
  result=data.get('data',{}).get('result',[]); \
  print(f'在线流数: {len(result)}'); \
  ids=sorted(set([r['metric']['id'] for r in result])); \
  print(f'不同ID数: {len(ids)}'); \
  print(f'ID列表: {ids[:20]}')"

# 查看有重连记录的流
curl -s 'http://localhost:9090/api/v1/query?query=video_stream_reconnect_total' | \
  python3 -c "import sys, json; data=json.load(sys.stdin); \
  result=data.get('data',{}).get('result',[]); \
  print(f'有重连的流: {len(result)}条'); \
  [print(f\"  {r['metric']['id']}: {r['value'][1]}次\") for r in result[:10]]"
```

### 2. 测试特定流

```bash
# 检查 D001 是否在线
curl -s 'http://localhost:9090/api/v1/query?query=video_stream_up{id="D001"}' | \
  python3 -m json.tool

# 检查 D072（有重连记录）
curl -s 'http://localhost:9090/api/v1/query?query=video_stream_reconnect_total{id="D072"}' | \
  python3 -m json.tool
```

### 3. 在 Grafana 中验证

1. 访问：http://localhost:3000/d/video-stream-monitoring
2. 设置过滤器：
   - 项目筛选 = g01
   - 桌台ID = D072（或其他有数据的ID）
   - 流名称 = All
3. 查看"🔄 重连次数"面板
4. 应该看到数据

## 其他常见问题

### 问题：所有指标都无数据

**可能原因**：
1. video-exporter 未运行
2. Prometheus 未抓取到数据
3. 时间范围选择错误

**排查步骤**：
```bash
# 1. 检查容器状态
docker-compose ps

# 2. 检查 video-exporter 指标端点
curl http://localhost:8080/metrics | grep video_stream

# 3. 检查 Prometheus targets
http://localhost:9090/targets

# 4. 检查时间范围
# 在 Grafana 右上角选择合适的时间范围
```

### 问题：Grafana 仪表板变量为空

**可能原因**：
1. Prometheus 还没有数据
2. 标签名称不匹配

**排查步骤**：
```bash
# 检查标签值
curl http://localhost:9090/api/v1/label/project/values
curl http://localhost:9090/api/v1/label/id/values
curl http://localhost:9090/api/v1/label/name/values
```

### 问题：数据延迟

**原因**：
- Prometheus 抓取间隔：30秒
- Grafana 自动刷新：30秒

**解决方案**：
- 等待 1-2 分钟
- 手动刷新仪表板

## 调试命令集合

### Prometheus 查询

```bash
# 查看所有流
curl 'http://localhost:9090/api/v1/query?query=video_stream_up'

# 查看特定项目
curl 'http://localhost:9090/api/v1/query?query=video_stream_up{project="g01"}'

# 查看特定桌台ID
curl 'http://localhost:9090/api/v1/query?query=video_stream_up{id="D072"}'

# 查看重连次数（带默认值）
curl 'http://localhost:9090/api/v1/query?query=video_stream_reconnect_total or vector(0)'

# 查看标签值
curl 'http://localhost:9090/api/v1/label/id/values'
```

### Docker 日志

```bash
# video-exporter 日志
docker logs video-exporter --tail 50

# Prometheus 日志
docker logs prometheus --tail 50

# Grafana 日志
docker logs grafana --tail 50

# 实时日志
./scripts/logs.sh [service_name]
```

### 服务状态

```bash
# 查看容器状态
docker-compose ps

# 查看资源使用
docker stats

# 查看网络
docker network ls
docker network inspect video-exporter_monitoring
```

## 性能优化

### 减少数据量

如果流数量很大，可以优化查询：

```promql
# 只查询在线的流
video_stream_reconnect_total{project="g01"} and on(id) video_stream_up{project="g01"} == 1

# 聚合查询
sum(rate(video_stream_reconnect_total[5m])) by (project)
```

### 调整采样参数

编辑 `config.yml`：
```yaml
exporter:
  check_interval: 30      # 增加检查间隔
  sample_duration: 10     # 减少采样时长
  max_concurrent: 100     # 调整并发数
```

## 联系支持

如果问题仍未解决：

1. 查看 [GitHub Issues](../../issues)
2. 提供以下信息：
   - 错误日志
   - 配置文件（删除敏感信息）
   - 环境信息（OS、Docker 版本等）
   - 重现步骤

## 参考文档

- [API 文档](API.md) - 指标详细说明
- [网络指标](NETWORK-METRICS.md) - 网络指标使用
- [部署检查清单](DEPLOYMENT-CHECKLIST.md) - 部署验证

