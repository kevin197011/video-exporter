package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Exporter Prometheus 导出器
type Exporter struct {
	streamUp       *prometheus.GaugeVec
	streamHealthy  *prometheus.GaugeVec
	streamPlayable *prometheus.GaugeVec
	totalPackets   *prometheus.GaugeVec
	videoPackets   *prometheus.GaugeVec
	audioPackets   *prometheus.GaugeVec
	keyframes      *prometheus.GaugeVec
	currentBitrate *prometheus.GaugeVec
	avgBitrate     *prometheus.GaugeVec
	framerate      *prometheus.GaugeVec
	responseTime   *prometheus.GaugeVec
	gopSize        *prometheus.GaugeVec
	qualityScore   *prometheus.GaugeVec
	stabilityScore *prometheus.GaugeVec

	scheduler *Scheduler
	log       *slog.Logger
}

// NewExporter 创建导出器
func NewExporter(scheduler *Scheduler) *Exporter {
	exporter := &Exporter{
		scheduler: scheduler,
		log:       GetLogger(),

		streamUp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_up",
				Help: "Stream is up (1) or down (0)",
			},
			[]string{"project", "id", "name", "url"},
		),

		streamHealthy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_healthy",
				Help: "Stream health status (1=healthy, 0=unhealthy)",
			},
			[]string{"project", "id", "name", "url"},
		),

		streamPlayable: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_playable",
				Help: "Stream is playable (1=yes, 0=no)",
			},
			[]string{"project", "id", "name", "url"},
		),

		totalPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_total_packets",
				Help: "Total packets received",
			},
			[]string{"project", "id", "name", "url"},
		),

		videoPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_video_packets",
				Help: "Video packets received",
			},
			[]string{"project", "id", "name", "url"},
		),

		audioPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_audio_packets",
				Help: "Audio packets received",
			},
			[]string{"project", "id", "name", "url"},
		),

		keyframes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_keyframes",
				Help: "Keyframes received",
			},
			[]string{"project", "id", "name", "url"},
		),

		currentBitrate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_bitrate_bps",
				Help: "Current stream bitrate in bits per second",
			},
			[]string{"project", "id", "name", "url"},
		),

		avgBitrate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_avg_bitrate_bps",
				Help: "Average stream bitrate in bits per second",
			},
			[]string{"project", "id", "name", "url"},
		),

		framerate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_framerate",
				Help: "Stream framerate in fps",
			},
			[]string{"project", "id", "name", "url"},
		),

		responseTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_response_ms",
				Help: "FLV HTTP request response time in milliseconds",
			},
			[]string{"project", "id", "name", "url"},
		),

		gopSize: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_gop_size",
				Help: "GOP size in frames",
			},
			[]string{"project", "id", "name", "url"},
		),

		qualityScore: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_quality_score",
				Help: "Stream quality score (0=poor, 1=fair, 2=good)",
			},
			[]string{"project", "id", "name", "url"},
		),

		stabilityScore: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "video_stream_stability_score",
				Help: "Bitrate stability score (0=unstable, 1=moderate, 2=stable)",
			},
			[]string{"project", "id", "name", "url"},
		),

		// resolution: prometheus.NewGaugeVec(
		// 	prometheus.GaugeOpts{
		// 		Name: "video_stream_resolution_pixels",
		// 		Help: "Video resolution in pixels (width * height)",
		// 	},
		// 	[]string{"project", "id", "url", "width", "height"},
		// ),
	}

	// 注册指标
	prometheus.MustRegister(
		exporter.streamUp,
		exporter.streamHealthy,
		exporter.streamPlayable,
		exporter.totalPackets,
		exporter.videoPackets,
		exporter.audioPackets,
		exporter.keyframes,
		exporter.currentBitrate,
		exporter.avgBitrate,
		exporter.framerate,
		exporter.responseTime,
		exporter.gopSize,
		exporter.qualityScore,
		exporter.stabilityScore,
		// exporter.resolution,
	)

	return exporter
}

// updateMetrics 更新指标
func (e *Exporter) updateMetrics() {
	e.log.Debug("开始更新指标")
	metrics := e.scheduler.GetAllMetrics()
	e.log.Debug("获取到指标", "数量", len(metrics))

	for _, m := range metrics {
		labels := []string{m.Project, m.ID, m.Name, m.URL}

		// 流状态
		upValue := 0.0
		if m.Healthy {
			upValue = 1.0
		}
		e.streamUp.WithLabelValues(labels...).Set(upValue)

		// 健康状态
		healthValue := 0.0
		if m.Healthy && m.ConsecutiveFails == 0 {
			healthValue = 1.0
		}
		e.streamHealthy.WithLabelValues(labels...).Set(healthValue)

		// 可播放状态
		playableValue := 0.0
		if m.Playable {
			playableValue = 1.0
		}
		e.streamPlayable.WithLabelValues(labels...).Set(playableValue)

		// 数据包统计
		e.totalPackets.WithLabelValues(labels...).Set(float64(m.TotalPackets))
		e.videoPackets.WithLabelValues(labels...).Set(float64(m.VideoPackets))
		e.audioPackets.WithLabelValues(labels...).Set(float64(m.AudioPackets))
		e.keyframes.WithLabelValues(labels...).Set(float64(m.Keyframes))

		// 码率指标
		e.currentBitrate.WithLabelValues(labels...).Set(m.CurrentBitrate)
		e.avgBitrate.WithLabelValues(labels...).Set(m.AvgBitrate)

		// 其他质量指标
		e.framerate.WithLabelValues(labels...).Set(m.Framerate)
		e.responseTime.WithLabelValues(labels...).Set(float64(m.Response))
		e.gopSize.WithLabelValues(labels...).Set(float64(m.GOPSize))

		// 质量评分
		qualityScore := 0.0
		switch m.Quality {
		case "good":
			qualityScore = 2.0
		case "fair":
			qualityScore = 1.0
		case "poor":
			qualityScore = 0.0
		}
		e.qualityScore.WithLabelValues(labels...).Set(qualityScore)

		// 稳定性评分
		stabilityScore := 0.0
		switch m.BitrateStability {
		case "stable":
			stabilityScore = 2.0
		case "moderate":
			stabilityScore = 1.0
		case "unstable":
			stabilityScore = 0.0
		}
		e.stabilityScore.WithLabelValues(labels...).Set(stabilityScore)

		// 分辨率 - 暂时注释掉
		// if m.Width > 0 && m.Height > 0 {
		// 	resLabels := append(labels, fmt.Sprintf("%d", m.Width), fmt.Sprintf("%d", m.Height))
		// 	pixels := float64(m.Width * m.Height)
		// 	e.resolution.WithLabelValues(resLabels...).Set(pixels)
		// }
	}

	e.log.Debug("指标更新完成")
}

// StartHTTPServer 启动 HTTP 服务器
func (e *Exporter) StartHTTPServer(addr string) error {
	mux := http.NewServeMux()

	// Prometheus metrics endpoint - 每次请求时更新指标
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		e.log.Debug("收到 metrics 请求")
		e.updateMetrics()
		e.log.Debug("指标更新完成")
		promhttp.Handler().ServeHTTP(w, r)
	})

	// 首页
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html>
<head><title>Video Stream Exporter</title></head>
<body>
<h1>Video Stream Exporter</h1>
<p><a href="/metrics">Metrics</a></p>
</body>
</html>`)
	})

	e.log.Info("Prometheus exporter 启动", "地址", addr)
	e.log.Info("访问指标", "URL", fmt.Sprintf("http://localhost%s/metrics", addr))

	return http.ListenAndServe(addr, mux)
}
