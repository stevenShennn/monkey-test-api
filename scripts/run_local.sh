#!/bin/bash

# 创建必要的目录
mkdir -p logs

# 设置环境变量
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="monkey_test"

# 构建并运行应用
echo "构建应用..."
go build -o bin/api-monkey-tester ./cmd/api-monkey-tester

echo "启动应用..."
./bin/api-monkey-tester 