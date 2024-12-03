package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"monkey-test-api/internal/api"
	"monkey-test-api/internal/config"
	"monkey-test-api/internal/logger"
	"monkey-test-api/internal/parser/curl"
	"monkey-test-api/internal/store"
	"monkey-test-api/internal/task"
	"monkey-test-api/internal/types"

	"github.com/BurntSushi/toml"
)

func main() {
	// 1. 创建带取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. 处理系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("接收到关闭信号，开始优雅关闭...")
		cancel()
	}()

	// 3. 加载配置
	cfg, err := config.LoadConfig("./conf/config.toml")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 4. 加载参数配置
	var paramConfig types.ParamConfig
	if _, err := toml.DecodeFile("./conf/param.toml", &paramConfig); err != nil {
		log.Fatal("加载参数配置失败:", err)
	}

	// 5. 初始化日志
	if err := logger.InitLogger(&cfg.Log); err != nil {
		log.Fatal("初始化日志失败:", err)
	}

	// 6. 初始化存储层
	var storage store.Store

	switch cfg.Database.Type {
	case "mysql":
		storage, err = store.NewMySQLStore(cfg.Database.MySQL.DSN)
		if err != nil {
			logger.Fatal("初始化 MySQL 失败:", err)
		}
	case "mongodb":
		storage, err = store.NewMongoStore(cfg.Database.MongoDB.URI, cfg.Database.MongoDB.Database)
		if err != nil {
			logger.Fatal("初始化 MongoDB 失败:", err)
		}
	default:
		logger.Fatal("不支持的数据库类型:", cfg.Database.Type)
	}
	defer storage.Close()

	// 7. 初始化各个组件
	parser := curl.NewParser()
	taskGenerator := task.NewGenerator(&paramConfig)

	// 8. 创建 API 处理器
	handler := api.NewHandler(parser, storage, taskGenerator)

	// 9. 设置路由
	router := api.SetupRouter(handler)

	// 10. 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 11. 启动服务器
	logger.Info("服务器启动在 :8080")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("服务器运行错误:", err)
			cancel()
		}
	}()

	// 12. 等待上下文取消
	<-ctx.Done()

	// 13. 优雅关闭
	logger.Info("开始关闭服务...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// 关闭 HTTP 服务器
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP 服务器关闭失败:", err)
	}

	// 关闭存储连接
	if err := storage.Close(); err != nil {
		logger.Error("关闭存储连接失败:", err)
	}

	logger.Info("服务已关闭")
}
