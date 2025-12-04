#!/bin/bash
# Grafana 配置检查脚本

set -e

# 切换到项目根目录
cd "$(dirname "$0")/.."

echo "=========================================="
echo "Grafana 配置检查"
echo "=========================================="
echo ""

# 1. 检查配置文件
echo "1. 检查 provisioning 配置文件..."
if [ -f "deployments/grafana/grafana-provisioning/datasources/prometheus.yml" ]; then
    echo "✅ datasources/prometheus.yml 存在"
else
    echo "❌ datasources/prometheus.yml 不存在"
    exit 1
fi

if [ -f "deployments/grafana/grafana-provisioning/dashboards/dashboard.yml" ]; then
    echo "✅ dashboards/dashboard.yml 存在"
else
    echo "❌ dashboards/dashboard.yml 不存在"
    exit 1
fi

if [ -f "deployments/grafana/grafana-provisioning/dashboards/video-stream-dashboard.json" ]; then
    echo "✅ video-stream-dashboard.json 存在"
else
    echo "❌ video-stream-dashboard.json 不存在"
    exit 1
fi

echo ""

# 2. 检查 Grafana 容器状态
echo "2. 检查 Grafana 容器状态..."
if docker ps | grep -q grafana; then
    echo "✅ Grafana 容器正在运行"
else
    echo "❌ Grafana 容器未运行"
    echo "请先启动服务: ./scripts/start.sh"
    exit 1
fi

echo ""

# 3. 检查容器内的 provisioning 目录
echo "3. 检查容器内的 provisioning 目录..."
echo "数据源配置:"
docker exec grafana ls -la /etc/grafana/provisioning/datasources/ 2>/dev/null || echo "❌ 无法访问 datasources 目录"

echo ""
echo "仪表板配置:"
docker exec grafana ls -la /etc/grafana/provisioning/dashboards/ 2>/dev/null || echo "❌ 无法访问 dashboards 目录"

echo ""

# 4. 检查 Grafana 日志
echo "4. 检查 Grafana 日志（最近 20 行）..."
docker logs grafana --tail 20 2>&1 | grep -i "provisioning\|datasource\|error" || echo "未发现 provisioning 相关日志"

echo ""

# 5. 测试 Prometheus 连接
echo "5. 测试 Prometheus 连接..."
if docker exec grafana wget -q -O- http://prometheus:9090/api/v1/status/config > /dev/null 2>&1; then
    echo "✅ Grafana 可以访问 Prometheus"
else
    echo "❌ Grafana 无法访问 Prometheus"
fi

echo ""

# 6. 检查 Grafana API
echo "6. 检查 Grafana 数据源 API..."
DATASOURCES=$(curl -s -u admin:admin http://localhost:3000/api/datasources 2>/dev/null)
if echo "$DATASOURCES" | grep -q "Prometheus"; then
    echo "✅ Prometheus 数据源已配置"
    echo "$DATASOURCES" | grep -o '"name":"[^"]*"' | head -3
else
    echo "❌ 未找到 Prometheus 数据源"
    echo "响应: $DATASOURCES"
fi

echo ""
echo "=========================================="
echo "检查完成"
echo "=========================================="
echo ""
echo "如果数据源未自动配置，请尝试："
echo "1. 重启 Grafana: docker-compose restart grafana"
echo "2. 查看完整日志: docker logs grafana"
echo "3. 手动添加数据源:"
echo "   - 访问 http://localhost:3000"
echo "   - 登录 (admin/admin)"
echo "   - Configuration → Data Sources → Add data source"
echo "   - 选择 Prometheus"
echo "   - URL: http://prometheus:9090"
echo "   - 点击 Save & Test"
echo ""

