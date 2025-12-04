#!/bin/bash
# Video Exporter 快速启动脚本

set -e

# 切换到项目根目录
cd "$(dirname "$0")/.."

echo "=========================================="
echo "Video Exporter + Prometheus + Grafana"
echo "=========================================="
echo ""

# 检查 docker 和 docker-compose 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装，请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "错误: Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 检查必要的配置文件
if [ ! -f "config.yml" ]; then
    echo "错误: config.yml 不存在"
    exit 1
fi

if [ ! -f "deployments/docker/prometheus.yml" ]; then
    echo "错误: deployments/docker/prometheus.yml 不存在"
    exit 1
fi

echo "✓ 环境检查通过"
echo ""

# 启动服务
echo "正在启动服务..."
docker-compose -f deployments/docker/docker-compose.yml up -d

echo ""
echo "=========================================="
echo "服务启动完成！"
echo "=========================================="
echo ""
echo "访问地址："
echo "  - Video Exporter:  http://localhost:8080"
echo "  - Video Metrics:   http://localhost:8080/metrics"
echo "  - Prometheus:      http://localhost:9090"
echo "  - Grafana:         http://localhost:3000"
echo ""
echo "Grafana 默认登录信息："
echo "  - 用户名: admin"
echo "  - 密码:   admin"
echo ""
echo "查看日志: ./scripts/logs.sh 或 docker-compose -f deployments/docker/docker-compose.yml logs -f"
echo "停止服务: ./scripts/stop.sh"
echo "删除服务: docker-compose -f deployments/docker/docker-compose.yml down"
echo ""
echo "等待服务启动..."
sleep 5

echo ""
echo "检查服务状态..."
docker-compose -f deployments/docker/docker-compose.yml ps

echo ""
echo "=========================================="
echo "提示："
echo "  - 首次启动 Grafana 可能需要1-2分钟初始化"
echo "  - 仪表板会自动导入，无需手动配置"
echo "  - 更多信息请查看 docs/DOCKER-COMPOSE-README.md"
echo "=========================================="

