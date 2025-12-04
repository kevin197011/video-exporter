#!/bin/bash
# Video Exporter 日志查看脚本

set -e

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 检查参数
SERVICE=${1:-""}

if [ -z "$SERVICE" ]; then
    echo "查看所有服务日志..."
    docker-compose -f deployments/docker/docker-compose.yml logs -f
else
    echo "查看 $SERVICE 服务日志..."
    docker-compose -f deployments/docker/docker-compose.yml logs -f "$SERVICE"
fi

