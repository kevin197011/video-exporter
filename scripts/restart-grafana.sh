#!/bin/bash
# Grafana 重启和验证脚本

set -e

# 切换到项目根目录
cd "$(dirname "$0")/.."

echo "=========================================="
echo "重启 Grafana 服务"
echo "=========================================="
echo ""

# 1. 停止并删除 Grafana 容器和数据
echo "1. 停止并清理 Grafana..."
docker-compose stop grafana
docker-compose rm -f grafana
docker volume rm video-exporter_grafana_data 2>/dev/null || true
echo "✅ Grafana 已清理"
echo ""

# 2. 重新启动
echo "2. 重新启动 Grafana..."
docker-compose up -d grafana
echo "✅ Grafana 已启动"
echo ""

# 3. 等待启动
echo "3. 等待 Grafana 初始化（30秒）..."
for i in {1..30}; do
    echo -n "."
    sleep 1
done
echo ""
echo "✅ 等待完成"
echo ""

# 4. 检查日志
echo "4. 检查启动日志..."
docker logs grafana --tail 30 | grep -i "provisioning\|datasource\|error\|HTTP Server Listen" || true
echo ""

# 5. 检查容器状态
echo "5. 检查容器状态..."
if docker ps | grep grafana | grep -q "Up"; then
    echo "✅ Grafana 容器运行正常"
else
    echo "❌ Grafana 容器未运行"
    echo "查看完整日志: docker logs grafana"
    exit 1
fi
echo ""

# 6. 测试 HTTP 连接
echo "6. 测试 Grafana HTTP 连接..."
for i in {1..10}; do
    if curl -s http://localhost:3000/api/health > /dev/null 2>&1; then
        echo "✅ Grafana HTTP 服务正常"
        break
    fi
    if [ $i -eq 10 ]; then
        echo "❌ Grafana HTTP 服务无响应"
        exit 1
    fi
    sleep 2
done
echo ""

# 7. 检查数据源
echo "7. 检查数据源配置（等待5秒后检查）..."
sleep 5
DATASOURCES=$(curl -s -u admin:admin http://localhost:3000/api/datasources)
if echo "$DATASOURCES" | grep -q "Prometheus"; then
    echo "✅ Prometheus 数据源已自动配置"
    echo ""
    echo "数据源详情:"
    echo "$DATASOURCES" | python3 -m json.tool 2>/dev/null | grep -A 5 "name" | head -10 || echo "$DATASOURCES"
else
    echo "⚠️  数据源可能还在加载中，请稍后检查"
    echo "手动检查: curl -u admin:admin http://localhost:3000/api/datasources"
fi
echo ""

echo "=========================================="
echo "Grafana 重启完成"
echo "=========================================="
echo ""
echo "访问地址: http://localhost:3000"
echo "用户名: admin"
echo "密码: admin"
echo ""
echo "如需查看详细日志:"
echo "  docker logs -f grafana"
echo ""

