# Metrics Export

## Purpose

Metrics Export 能力负责将收集的流质量指标和网络指标导出为 Prometheus 格式，通过 HTTP 端点提供指标查询服务。

## Requirements

### Requirement: Prometheus Metrics Endpoint

The system SHALL expose a Prometheus-compatible metrics endpoint.

#### Scenario: Metrics endpoint availability
- **WHEN** the exporter starts
- **THEN** the system SHALL start an HTTP server on the configured `listen_addr`
- **AND** expose metrics at `/metrics` endpoint
- **AND** respond to GET requests with Prometheus text format

#### Scenario: Metrics endpoint access
- **WHEN** a client requests `http://<listen_addr>/metrics`
- **THEN** the system SHALL return all current metrics in Prometheus text format
- **AND** include proper Content-Type header

### Requirement: Metric Labeling

The system SHALL label all metrics with consistent labels.

#### Scenario: Standard labels
- **WHEN** metrics are exported
- **THEN** each metric SHALL include labels: `project`, `id`, `name`, `url`
- **AND** optionally `service` label set to "video-exporter"

#### Scenario: Label consistency
- **WHEN** multiple metrics are exported for the same stream
- **THEN** all metrics SHALL use identical label values for that stream
- **AND** ensure label uniqueness per stream

### Requirement: Gauge Metrics

The system SHALL export current state metrics as Gauge type.

#### Scenario: Stream status gauge
- **WHEN** a stream check completes
- **THEN** the system SHALL export `video_stream_up` as Gauge (1=up, 0=down)
- **AND** update the value based on current state

#### Scenario: Quality metrics as gauges
- **WHEN** quality metrics are calculated
- **THEN** the system SHALL export them as Gauge type:
  - `video_stream_bitrate_bps`
  - `video_stream_avg_bitrate_bps`
  - `video_stream_framerate`
  - `video_stream_response_ms`
  - `video_stream_gop_size`
  - `video_stream_quality_score`
  - `video_stream_stability_score`

#### Scenario: Network metrics as gauges
- **WHEN** network metrics are calculated
- **THEN** the system SHALL export them as Gauge type:
  - `video_stream_rtt_ms`
  - `video_stream_packet_loss_ratio`
  - `video_stream_network_jitter_ms`
  - `video_stream_reconnect_total` (Gauge, not Counter)

### Requirement: Counter Metrics

The system SHALL use Counter type only for cumulative values that never decrease.

#### Scenario: No counter metrics currently
- **WHEN** metrics are exported
- **THEN** the system SHALL NOT use Counter type for reconnect tracking
- **AND** use Gauge instead to represent current-period reconnects
- **NOTE**: Reconnect tracking uses Gauge, not Counter, as it resets per cycle

### Requirement: Metric Updates

The system SHALL update metrics after each stream check cycle.

#### Scenario: Periodic metric updates
- **WHEN** a stream check completes
- **THEN** the system SHALL update all relevant metrics for that stream
- **AND** make updated values immediately available via `/metrics` endpoint

#### Scenario: Metric removal for removed streams
- **WHEN** a stream is removed from configuration
- **THEN** the system SHALL stop updating metrics for that stream
- **AND** metrics MAY remain in Prometheus until expiration (Prometheus behavior)

### Requirement: Metric Help Text

The system SHALL provide descriptive help text for all metrics.

#### Scenario: Help text for each metric
- **WHEN** metrics are registered
- **THEN** each metric SHALL include a `Help` field describing:
  - What the metric measures
  - Unit of measurement (if applicable)
  - Value range or meaning

#### Scenario: Help text examples
- `video_stream_up`: "Stream is up (1) or down (0)"
- `video_stream_bitrate_bps`: "Current bitrate in bits per second"
- `video_stream_rtt_ms`: "Round-trip time in milliseconds"

### Requirement: Metric Registration

The system SHALL register all metrics with Prometheus registry.

#### Scenario: Metric registration on startup
- **WHEN** the exporter is created
- **THEN** the system SHALL create and register all metric vectors
- **AND** register them with the default Prometheus registry

#### Scenario: Metric vector creation
- **WHEN** creating metric vectors
- **THEN** the system SHALL use `prometheus.NewGaugeVec` for gauge metrics
- **AND** specify label names: `[]string{"project", "id", "name", "url"}`

### Requirement: Metrics Collection from Scheduler

The system SHALL collect metrics from the scheduler for all streams.

#### Scenario: Metrics collection
- **WHEN** updating metrics
- **THEN** the system SHALL call `scheduler.GetAllMetrics()`
- **AND** iterate through returned stream metrics
- **AND** update Prometheus metrics for each stream

#### Scenario: Handling missing streams
- **WHEN** a stream is no longer in scheduler metrics
- **THEN** the system SHALL handle gracefully
- **AND** not update metrics for that stream (they will expire in Prometheus)

### Requirement: HTTP Server Management

The system SHALL manage the HTTP server lifecycle.

#### Scenario: Server startup
- **WHEN** the exporter starts
- **THEN** the system SHALL start the HTTP server in a goroutine
- **AND** handle errors gracefully

#### Scenario: Server shutdown
- **WHEN** the application receives a shutdown signal
- **THEN** the system SHALL gracefully shut down the HTTP server
- **AND** allow in-flight requests to complete

### Requirement: Metrics Format Compliance

The system SHALL export metrics in Prometheus text format specification.

#### Scenario: Format compliance
- **WHEN** metrics are exported
- **THEN** the output SHALL comply with Prometheus text format:
  - Metric name and labels on one line
  - Value on the next line
  - Help text with `# HELP` prefix
  - Type information with `# TYPE` prefix

#### Scenario: Metric name validation
- **WHEN** creating metrics
- **THEN** metric names SHALL:
  - Start with a letter
  - Contain only letters, digits, and underscores
  - Match Prometheus naming conventions

