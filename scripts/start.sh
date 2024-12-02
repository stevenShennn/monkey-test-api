#!/bin/bash

# 清理已存在的容器
echo "清理已存在的容器..."
docker-compose down

# 构建并启动服务
echo "构建并启动服务..."
docker-compose up --build -d

# 等待 MongoDB 就绪
echo "等待 MongoDB 就绪..."
max_retries=30
count=0
while [ $count -lt $max_retries ]; do
    if docker-compose exec mongodb mongosh --eval "db.runCommand('ping').ok" --quiet; then
        echo "MongoDB 已就绪"
        break
    fi
    echo "等待 MongoDB 就绪... ($(( count + 1 ))/$max_retries)"
    sleep 2
    count=$(( count + 1 ))
done

if [ $count -eq $max_retries ]; then
    echo "错误: MongoDB 启动失败"
    docker-compose logs mongodb
    exit 1
fi

# 等待 API 服务就绪
echo "等待 API 服务就绪..."
max_retries=30
count=0
while [ $count -lt $max_retries ]; do
    if curl -s http://localhost:8080/health > /dev/null; then
        echo "API 服务已就绪"
        break
    fi
    echo "等待 API 服务就绪... ($(( count + 1 ))/$max_retries)"
    sleep 2
    count=$(( count + 1 ))
done

if [ $count -eq $max_retries ]; then
    echo "错误: API 服务启动失败"
    docker-compose logs api
    exit 1
fi

# 检查所有服务状态
echo "检查服务状态..."
docker-compose ps

echo "所有服务已启动成功:"
echo "- API 服务: http://localhost:8080"
echo "- MongoDB: localhost:27017"
echo "- Mongo Express: http://localhost:8081"

# 测试 API 连接
echo "测试 API 连接..."
sleep 5  # 额外等待以确保服务完全就绪

response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/api/v1/test-requests \
  -H "Content-Type: application/json" \
  -d '{
    "curl": "curl \"http://example.com/api/test\""
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo "API 测试成功！"
    echo "响应: $body"
else
    echo "API 测试失败！"
    echo "HTTP 状态码: $http_code"
    echo "响应: $body"
    echo "查看详细日志..."
    docker-compose logs api
fi 