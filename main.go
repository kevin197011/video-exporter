package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化日志
	InitLogger()
	log := GetLogger()

	log.Info("启动 Video Stream Exporter")

	// 加载配置
	cfg, err := LoadConfig("config.yml")
	if err != nil {
		log.Error("加载配置失败", "错误", err)
		os.Exit(1)
	}

	// 设置日志级别
	SetLogLevel(cfg.Exporter.LogLevel)

	// 设置全局配置
	SetGlobalConfig(cfg)

	// 创建调度器
	scheduler := NewScheduler(cfg)

	// 添加所有流
	totalStreams := 0
	for project, streams := range cfg.Streams {
		log.Info("加载项目", "项目", project, "流数量", len(streams))
		for _, stream := range streams {
			scheduler.AddStream(stream.ID, stream.URL, project)
			totalStreams++
		}
	}

	log.Info("已加载流", "总数", totalStreams)

	// 启动调度器
	go scheduler.Start()

	// 创建并启动 Prometheus exporter
	exporter := NewExporter(scheduler)

	listenAddr := cfg.Exporter.ListenAddr
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	if listenAddr[0] != ':' {
		listenAddr = ":" + listenAddr
	}

	// 启动 HTTP 服务器
	go func() {
		if err := exporter.StartHTTPServer(listenAddr); err != nil {
			log.Error("HTTP 服务器错误", "错误", err)
		}
	}()

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info("服务已启动，按 Ctrl+C 停止")

	<-sigChan
	log.Info("收到停止信号")

	// 停止调度器
	scheduler.Stop()

	log.Info("服务已停止")
}
