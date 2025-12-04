package logger

import (
	"log/slog"
	"os"
	"strings"
)

var instance *slog.Logger
var levelVar slog.LevelVar

// Init 初始化日志
func Init() {
	levelVar.Set(slog.LevelInfo)
	// 只输出到标准输出，适合 Docker 环境
	instance = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: &levelVar,
	}))
}

// Get 获取日志实例
func Get() *slog.Logger {
	if instance == nil {
		Init()
	}
	return instance
}

// SetLevel 设置日志级别
func SetLevel(level string) {
	if level == "" {
		level = "info"
	}

	switch strings.ToLower(level) {
	case "debug":
		levelVar.Set(slog.LevelDebug)
	case "info":
		levelVar.Set(slog.LevelInfo)
	case "warn", "warning":
		levelVar.Set(slog.LevelWarn)
	case "error":
		levelVar.Set(slog.LevelError)
	default:
		levelVar.Set(slog.LevelInfo)
	}

	if instance != nil {
		instance.Info("日志级别已更新", "level", levelVar.Level())
	}
}
