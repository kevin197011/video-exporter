package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置结构
type Config struct {
	Exporter ExporterConfig            `yaml:"exporter"`
	Streams  map[string][]StreamConfig `yaml:"streams"` // project -> streams
}

// ExporterConfig 导出器配置
type ExporterConfig struct {
	CheckInterval  int    `yaml:"check_interval"`  // 检查间隔（秒）
	SampleDuration int    `yaml:"sample_duration"` // 采样时长（秒），默认10秒
	MinKeyframes   int    `yaml:"min_keyframes"`   // 最小关键帧数，默认2
	MaxConcurrent  int    `yaml:"max_concurrent"`
	MaxRetries     int    `yaml:"max_retries"`
	ListenAddr     string `yaml:"listen_addr"` // Prometheus exporter 监听地址
	LogLevel       string `yaml:"log_level"`   // 日志级别
}

// StreamConfig 流配置
type StreamConfig struct {
	URL string `yaml:"url"`
	ID  string `yaml:"id"`
}

// Load 加载配置文件
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// 全局配置
var globalConfig *Config

// SetGlobal 设置全局配置
func SetGlobal(cfg *Config) {
	globalConfig = cfg
}

// GetGlobal 获取全局配置
func GetGlobal() *Config {
	return globalConfig
}
