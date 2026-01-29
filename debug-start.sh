#!/bin/bash

cd /www/wwwroot/task-trade-platform

export PATH=/usr/local/go/bin:$PATH
export GO111MODULE=on

echo "=== 尝试启动完整Go后端 ==="
echo "配置文件: ./configs/config-optimized.yaml"
echo ""

# 直接运行并捕获所有输出
/usr/local/go/bin/go run cmd/server/main.go -config=./configs/config-optimized.yaml
