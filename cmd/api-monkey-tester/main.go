package main

import (
	"context"
	"log"
	"monkey-test-api/internal/api"
	"monkey-test-api/internal/config"
	"monkey-test-api/internal/input"
	"monkey-test-api/internal/input/curl"
	"monkey-test-api/internal/logger"
	"monkey-test-api/internal/store"
	"monkey-test-api/internal/task"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("./conf/config.toml")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 2. 初始化日志
	if err := logger.InitLogger(&cfg.Log); err != nil {
		log.Fatal("初始化日志失败:", err)
	}

	// 3. 连接数据库
	ctx := context.Background()

	// 初始化存储层
	var storage store.Store
	var err error

	switch cfg.Database.Type {
	case "mysql":
		storage, err = store.NewMySQLStore(cfg.Database.MySQL.DSN)
		if err != nil {
			logger.Fatal("初始化 MySQL 失败:", err)
		}
	case "mongodb":
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.MongoDB.URI))
		if err != nil {
			logger.Fatal("连接 MongoDB 失败:", err)
		}
		defer client.Disconnect(ctx)
		storage = store.NewMongoStore(client.Database(cfg.Database.MongoDB.Database))
	default:
		logger.Fatal("不支持的数据库类型:", cfg.Database.Type)
	}

	// 4. 初始化各个组件
	parser := curl.NewParser()
	taskGenerator := task.NewGenerator(&cfg.Task)

	// 5. 创建 API 处理器
	handler := api.NewHandler(parser, storage, taskGenerator)

	// 6. 设置路由
	router := api.SetupRouter(handler)

	// 7. 启动服务器
	logger.Info("服务器启动在 :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Fatal("启动服务器失败:", err)
	}
}
