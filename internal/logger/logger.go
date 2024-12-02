package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	// Log 全局日志实例
	Log *zap.SugaredLogger
)

// Config 日志配置
type Config struct {
	Level      string `toml:"level"`       // 日志级别
	Filename   string `toml:"filename"`    // 日志文件路径
	MaxSize    int    `toml:"maxsize"`     // 每个日志文件最大尺寸，单位MB
	MaxBackups int    `toml:"maxbackups"`  // 保留的旧日志文件最大数量
	MaxAge     int    `toml:"maxage"`      // 保留的旧日志文件最大天数
	Compress   bool   `toml:"compress"`    // 是否压缩旧日志文件
	Console    bool   `toml:"console"`     // 是否同时输出到控制台
}

// InitLogger 初始化日志
func InitLogger(cfg *Config) error {
	// 设置日志级别
	level := zap.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return err
	}

	// 创建核心日志配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置日志输出
	var cores []zapcore.Core

	// 文件输出
	if cfg.Filename != "" {
		writer := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(writer),
			level,
		))
	}

	// 控制台输出
	if cfg.Console {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// 创建日志实例
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 设置全局日志实例
	Log = logger.Sugar()
	return nil
}

// Debug 输出调试日志
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Debugf 输出格式化的调试日志
func Debugf(template string, args ...interface{}) {
	Log.Debugf(template, args...)
}

// Info 输出信息日志
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Infof 输出格式化的信息日志
func Infof(template string, args ...interface{}) {
	Log.Infof(template, args...)
}

// Warn 输出警告日志
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Warnf 输出格式化的警告日志
func Warnf(template string, args ...interface{}) {
	Log.Warnf(template, args...)
}

// Error 输出错误日志
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Errorf 输出格式化的错误日志
func Errorf(template string, args ...interface{}) {
	Log.Errorf(template, args...)
}

// Fatal 输出致命错误日志并退出程序
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// Fatalf 输出格式化的致命错误日志并退出程序
func Fatalf(template string, args ...interface{}) {
	Log.Fatalf(template, args...)
} 