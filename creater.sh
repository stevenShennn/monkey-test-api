#!/bin/bash

# 创建cmd层目录结构
mkdir -p cmd/api-monkey-tester

# 创建内部目录结构
mkdir -p internal/{api,config,utils}

# 创建空文件
touch cmd/api-monkey-tester/main.go

touch internal/api/request.go
touch internal/api/mutation.go
touch internal/api/response.go
touch internal/config/config.go
touch internal/utils/logger.go
touch internal/utils/utils.go

# 创建脚本目录和文件
mkdir -p scripts
touch scripts/test.sh

# 创建Go模块文件
touch go.mod
touch go.sum

# 创建项目说明文件
touch README.md

echo "所有文件和目录已创建！"
