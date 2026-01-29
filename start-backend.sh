#!/bin/bash

# 任务平台后端启动脚本

cd /www/wwwroot/task-trade-platform

# 设置Go环境
export PATH=/usr/local/go/bin:$PATH
export GO111MODULE=on

# 创建必要的目录
mkdir -p logs
mkdir -p /var/log/task-platform

echo "=== 启动任务平台后端服务 ==="
echo "时间: $(date)"
echo "配置文件: ./configs/config-optimized.yaml"
echo ""

# 检查8080端口是否被占用
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "警告: 端口8080已被占用"
    echo "正在尝试停止旧进程..."
    lsof -Pi :8080 -sTCP:LISTEN -t | xargs kill -9
    sleep 2
fi

# 启动服务
nohup /usr/local/go/bin/go run cmd/server/main.go -config=./configs/config-optimized.yaml >> logs/backend.log 2>&1 &

PID=$!
echo "服务已启动, PID: $PID"

# 等待几秒查看启动情况
sleep 3

if ps -p $PID > /dev/null; then
    echo "✓ 服务运行正常"
    echo "  日志文件: logs/backend.log"
    echo "  访问地址: http://49.234.39.189:8080/api/"
    echo ""
    echo "最近的日志输出:"
    tail -20 logs/backend.log
else
    echo "✗ 服务启动失败"
    echo "错误日志:"
    cat logs/backend.log
    exit 1
fi
