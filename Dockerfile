# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o api-monkey-tester ./cmd/api-monkey-tester

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/api-monkey-tester .

# 复制配置文件
COPY conf/ ./conf/

# 创建日志目录
RUN mkdir -p logs

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./api-monkey-tester"] 