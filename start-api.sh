#!/bin/bash

cd /www/wwwroot/task-trade-platform

# 查找并杀掉旧进程
OLD_PID=$(ss -tulnp | grep :8080 | grep -oP 'pid=\K[0-9]+' | head -1)
if [ ! -z "$OLD_PID" ]; then
    echo "正在停止旧进程 PID: $OLD_PID"
    kill -9 $OLD_PID
    sleep 1
fi

# 启动新服务
echo "正在启动API服务..."
nohup ./simple-server-db > logs/api.log 2>&1 &
NEW_PID=$!

sleep 2

# 检查是否成功启动
if ps -p $NEW_PID > /dev/null; then
    echo "✅ API服务已启动，PID: $NEW_PID"
    echo "📡 访问地址: http://49.234.39.189:8080/api/"
    ss -tulnp | grep :8080
else
    echo "❌ API服务启动失败"
    echo "查看日志:"
    tail -20 logs/api.log
    exit 1
fi
