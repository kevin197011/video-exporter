package main

import (
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	urlpkg "net/url"
	pathpkg "path"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/nareix/joy5/av"
	"github.com/nareix/joy5/format/flv"
)

// StreamChecker 流检查器
type StreamChecker struct {
	id      string
	url     string
	project string
	name    string

	// 统计数据（当前检查的值，不累积）
	mu               sync.RWMutex
	totalPackets     int64 // 本次检查的总包数
	videoPackets     int64 // 本次检查的视频包数
	audioPackets     int64 // 本次检查的音频包数
	keyframes        int64 // 本次检查的关键帧数
	currentBitrate   float64
	avgBitrate       float64
	bitrateHistory   []float64
	framerate        float64
	codec            string
	response         int64
	gopSize          int
	width            int
	height           int
	quality          string
	playable         bool
	bitrateStability string
	healthy          bool
	lastCheckTime    time.Time
	consecutiveFails int

	log *slog.Logger
}

// extractStreamName 从 URL 和 ID 提取流名称
// 例如: project=project1, id=stream-01, url=https://example.com/path/stream.flv
// 结果: project1_example_stream-01_path_stream
func extractStreamName(project, id, rawURL string) string {
	hostSegment := "unknown"
	pathSegment := "unknown"

	if parsed, err := urlpkg.Parse(rawURL); err == nil {
		// host 取第一个子域（例如：example.com -> example）
		if host := parsed.Hostname(); host != "" {
			if parts := strings.Split(host, "."); len(parts) > 0 && parts[0] != "" {
				hostSegment = parts[0]
			}
		}

		// path: 去掉扩展名，替换斜杠
		p := strings.TrimPrefix(parsed.Path, "/")
		if p != "" {
			if ext := pathpkg.Ext(p); ext != "" {
				p = strings.TrimSuffix(p, ext)
			}
			p = strings.ReplaceAll(p, "/", "_")
			if p != "" {
				pathSegment = p
			}
		}
	} else {
		// 解析失败的兜底
		re := regexp.MustCompile(`https?://([^/]+)/(.+)`)
		if matches := re.FindStringSubmatch(rawURL); len(matches) >= 3 {
			host := matches[1]
			if host != "" {
				if parts := strings.Split(host, "."); len(parts) > 0 && parts[0] != "" {
					hostSegment = parts[0]
				}
			}

			p := matches[2]
			p = strings.TrimSuffix(p, ".flv")
			p = strings.TrimSuffix(p, ".m3u8")
			p = strings.ReplaceAll(p, "/", "_")
			if p != "" {
				pathSegment = p
			}
		}
	}

	return fmt.Sprintf("%s_%s_%s_%s", project, hostSegment, id, pathSegment)
}

// NewStreamChecker 创建流检查器
func NewStreamChecker(id, url, project string) *StreamChecker {
	return &StreamChecker{
		id:             id,
		url:            url,
		project:        project,
		name:           extractStreamName(project, id, url),
		healthy:        false,
		playable:       false,
		quality:        "unknown",
		bitrateHistory: make([]float64, 0, 10),
		log:            GetLogger(),
	}
}

// Check 执行一次流检查
func (sc *StreamChecker) Check(timeout time.Duration) error {
	sc.log.Debug("开始检查流", "流ID", sc.id, "URL", sc.url)

	startTime := time.Now()

	// 创建 HTTP 客户端 - 不设置超时，因为我们需要持续读取
	client := &http.Client{
		Timeout: 0, // 不限制超时，由我们自己控制
	}

	// 记录请求开始时间，用于计算HTTP-FLV请求响应时间
	reqStart := time.Now()
	req, err := http.NewRequest("GET", sc.url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
	}

	// 将延迟定义为 FLV 的 HTTP 请求响应时间（ms）
	sc.mu.Lock()
	sc.response = time.Since(reqStart).Milliseconds()
	sc.mu.Unlock()

	// 创建解复用器
	demuxer := flv.NewDemuxer(resp.Body)

	// joy5 不需要预先获取流信息，直接读取包即可
	hasVideo := false
	hasMetadata := false

	// 采样数据包 - 基于时间采样，更真实
	packetCount := 0
	videoCount := 0
	audioCount := 0
	keyframeCount := 0
	totalBytes := int64(0)

	// 从配置读取采样参数，如果未配置则使用默认值
	sampleDurationSec := 10
	minKeyframes := 2
	if globalConfig != nil {
		if globalConfig.Exporter.SampleDuration > 0 {
			sampleDurationSec = globalConfig.Exporter.SampleDuration
		}
		if globalConfig.Exporter.MinKeyframes > 0 {
			minKeyframes = globalConfig.Exporter.MinKeyframes
		}
	}
	sampleDuration := time.Duration(sampleDurationSec) * time.Second
	sampleStartTime := time.Now()

	// 用于延迟计算的变量
	firstPacketTime := time.Time{} // 第一个视频包到达的系统时间（用于是否读到包的判定）
	firstDTS := int64(0)           // 第一个视频包的DTS
	lastDTS := int64(0)            // 最后一个视频包的DTS
	keyframeInterval := 0

	for {
		// 基于时间的采样，更能反映实际情况
		if time.Since(sampleStartTime) >= sampleDuration && keyframeCount >= minKeyframes {
			break
		}

		pktRecvTime := time.Now() // 记录包到达时间
		pkt, err := demuxer.ReadPacket()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("读取数据包失败: %w", err)
		}

		packetCount++
		totalBytes += int64(len(pkt.Data))

		// 检查 metadata
		if pkt.Type == av.Metadata && !hasMetadata {
			hasMetadata = true
			// 尝试解析 metadata，看是否有时间信息
			sc.log.Debug("收到Metadata",
				"流ID", sc.id,
				"数据长度", len(pkt.Metadata),
				"数据", string(pkt.Metadata[:min(len(pkt.Metadata), 200)]))
		}

		// joy5: 使用 Type 判断包类型
		switch pkt.Type {
		case av.H264:
			videoCount++
			hasVideo = true

			if pkt.IsKeyFrame {
				keyframeCount++
			}

			// 记录时间戳和到达时间
			if firstPacketTime.IsZero() {
				firstPacketTime = pktRecvTime
				firstDTS = int64(pkt.Time)
			}
			lastDTS = int64(pkt.Time)

			// 获取编码信息
			if sc.codec == "" {
				sc.mu.Lock()
				sc.codec = "H264"
				sc.mu.Unlock()
			}
		case av.AAC:
			audioCount++
		}
	}

	if !hasVideo {
		return fmt.Errorf("未找到视频流")
	}

	duration := time.Since(startTime)

	// 计算 GOP 大小（关键帧间隔的帧数）
	if keyframeCount > 1 {
		// 简单方法：总帧数 / 关键帧数
		keyframeInterval = videoCount / keyframeCount
	} else if keyframeCount == 1 {
		// 只有一个关键帧，GOP就是所有帧
		keyframeInterval = videoCount
	} else {
		// 没有关键帧，设为0
		keyframeInterval = 0
	}

	// 更新统计数据（记录本次检查的值）
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.totalPackets = int64(packetCount)
	sc.videoPackets = int64(videoCount)
	sc.audioPackets = int64(audioCount)
	sc.keyframes = int64(keyframeCount)
	sc.lastCheckTime = time.Now()
	sc.healthy = true
	sc.consecutiveFails = 0
	sc.gopSize = keyframeInterval

	// 计算帧率和码率（基于 DTS 时间，更准确）
	if !firstPacketTime.IsZero() && lastDTS > firstDTS {
		dtsElapsed := float64(lastDTS-firstDTS) / 1e9 // 纳秒转秒
		if dtsElapsed > 0 {
			sc.framerate = float64(videoCount) / dtsElapsed
			// 基于 DTS 时间计算码率更准确
			sc.currentBitrate = (float64(totalBytes) * 8) / dtsElapsed // bps
		}
	} else if duration.Seconds() > 0 {
		// 如果没有 DTS，使用实际耗时
		sc.currentBitrate = (float64(totalBytes) * 8) / duration.Seconds() // bps
	}

	// 更新码率历史
	if sc.currentBitrate > 0 {
		sc.bitrateHistory = append(sc.bitrateHistory, sc.currentBitrate)
		if len(sc.bitrateHistory) > 10 {
			sc.bitrateHistory = sc.bitrateHistory[1:]
		}

		// 计算平均码率
		sum := 0.0
		for _, br := range sc.bitrateHistory {
			sum += br
		}
		sc.avgBitrate = sum / float64(len(sc.bitrateHistory))

		// 评估码率稳定性
		if len(sc.bitrateHistory) >= 3 {
			variance := 0.0
			for _, br := range sc.bitrateHistory {
				diff := br - sc.avgBitrate
				variance += diff * diff
			}
			variance /= float64(len(sc.bitrateHistory))
			stdDev := math.Sqrt(variance)

			// 计算变异系数（CV = 标准差/平均值）
			cv := stdDev / sc.avgBitrate

			// 根据变异系数评估稳定性
			// CV < 0.15 (15%) = 稳定
			// CV < 0.30 (30%) = 中等
			// CV >= 0.30 = 不稳定
			if cv < 0.15 {
				sc.bitrateStability = "stable"
			} else if cv < 0.30 {
				sc.bitrateStability = "moderate"
			} else {
				sc.bitrateStability = "unstable"
			}
		} else {
			sc.bitrateStability = "unknown"
		}
	}

	// 此处延迟已定义为 HTTP-FLV 请求响应时间（在完成HTTP响应后已设置）

	// 评估质量
	sc.playable = keyframeCount >= 2 && videoCount > 10
	if sc.playable {
		// 质量评估：基于帧率、码率和稳定性
		if sc.framerate >= 25 && sc.currentBitrate >= 600000 {
			// 高质量：帧率>=25fps，码率>=600kbps
			sc.quality = "good"
		} else if sc.framerate >= 20 && sc.currentBitrate >= 400000 {
			// 中等质量：帧率>=20fps，码率>=400kbps
			sc.quality = "fair"
		} else {
			// 低质量
			sc.quality = "poor"
		}
	} else {
		sc.quality = "poor"
	}

	// 注意：这里已经持有 mu.Lock()，不需要再加锁
	sc.log.Debug("检查完成",
		"流ID", sc.id,
		"耗时秒", fmt.Sprintf("%.2f", duration.Seconds()),
		"可播放", sc.playable,
		"质量", sc.quality,
		"请求响应ms", sc.response,
		"视频包", videoCount,
		"关键帧", keyframeCount,
		"码率kbps", fmt.Sprintf("%.1f", sc.currentBitrate/1000),
		"平均码率kbps", fmt.Sprintf("%.1f", sc.avgBitrate/1000),
		"稳定性", sc.bitrateStability,
		"帧率fps", fmt.Sprintf("%.1f", sc.framerate),
		"GOP帧", sc.gopSize,
		"编码", sc.codec)

	return nil
}

// MarkFailed 标记检查失败
func (sc *StreamChecker) MarkFailed() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.consecutiveFails++
	sc.healthy = false
	sc.playable = false
	sc.totalPackets = 0
	sc.videoPackets = 0
	sc.audioPackets = 0
	sc.keyframes = 0
	sc.currentBitrate = 0
	sc.avgBitrate = 0
	sc.framerate = 0
	sc.codec = ""
	sc.response = 0
	sc.gopSize = 0
	sc.width = 0
	sc.height = 0
	sc.quality = "poor"
	sc.bitrateStability = "unstable"
	sc.lastCheckTime = time.Now()
}

// GetMetrics 获取指标
func (sc *StreamChecker) GetMetrics() StreamMetrics {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return StreamMetrics{
		ID:               sc.id,
		URL:              sc.url,
		Project:          sc.project,
		Name:             sc.name,
		TotalPackets:     sc.totalPackets,
		VideoPackets:     sc.videoPackets,
		AudioPackets:     sc.audioPackets,
		Keyframes:        sc.keyframes,
		CurrentBitrate:   sc.currentBitrate,
		AvgBitrate:       sc.avgBitrate,
		Framerate:        sc.framerate,
		Codec:            sc.codec,
		Response:         sc.response,
		GOPSize:          sc.gopSize,
		Width:            sc.width,
		Height:           sc.height,
		Quality:          sc.quality,
		Playable:         sc.playable,
		BitrateStability: sc.bitrateStability,
		Healthy:          sc.healthy,
		LastCheckTime:    sc.lastCheckTime,
		ConsecutiveFails: sc.consecutiveFails,
	}
}

// StreamMetrics 流指标
type StreamMetrics struct {
	ID               string
	URL              string
	Project          string
	Name             string
	TotalPackets     int64
	VideoPackets     int64
	AudioPackets     int64
	Keyframes        int64
	CurrentBitrate   float64
	AvgBitrate       float64
	Framerate        float64
	Codec            string
	Response         int64
	GOPSize          int
	Width            int
	Height           int
	Quality          string
	Playable         bool
	BitrateStability string
	Healthy          bool
	LastCheckTime    time.Time
	ConsecutiveFails int
}
