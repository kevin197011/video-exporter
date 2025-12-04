#!/bin/bash
# Video Exporter 停止脚本

set -e

# 切换到项目根目录
cd "$(dirname "$0")/.."

echo "=========================================="
echo "停止 Video Exporter 服务"
echo "=========================================="
echo ""

# 检查是否有运行的容器
if ! docker-compose -f deployments/docker/docker-compose.yml ps | grep -q "Up"; then
    echo "没有运行中的服务"
    exit 0
fi

echo "正在停止服务..."
docker-compose -f deployments/docker/docker-compose.yml stop

echo ""
echo "服务已停止"
echo ""
echo "如需完全删除容器（保留数据）："
echo "  docker-compose -f deployments/docker/docker-compose.yml down"
echo ""
echo "如需删除容器和数据："
echo "  docker-compose -f deployments/docker/docker-compose.yml down -v"
echo ""

