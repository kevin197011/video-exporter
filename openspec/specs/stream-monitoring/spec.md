# Stream Monitoring

## Purpose

Stream Monitoring 是系统的核心能力，负责定期检查视频流的健康状况，收集质量指标，并判断流的可播放性和健康状态。

## Requirements

### Requirement: Stream Health Check

The system SHALL periodically check the health status of configured video streams.

#### Scenario: Successful stream check
- **WHEN** a stream is configured with a valid URL
- **AND** the stream is accessible and streaming
- **THEN** the system SHALL mark the stream as up (`video_stream_up = 1`)
- **AND** collect quality metrics (bitrate, framerate, resolution, etc.)

#### Scenario: Failed stream check
- **WHEN** a stream URL is unreachable or returns an error
- **THEN** the system SHALL mark the stream as down (`video_stream_up = 0`)
- **AND** increment the reconnect counter if previously up
- **AND** reset network metrics (RTT, packet loss, jitter)

#### Scenario: Stream check with retry
- **WHEN** a stream check fails
- **AND** the failure count is less than `max_retries`
- **THEN** the system SHALL retry the connection
- **AND** increment the retry counter

### Requirement: Stream Sampling

The system SHALL sample stream data for a configurable duration to collect quality metrics.

#### Scenario: Sample duration configuration
- **WHEN** `sample_duration` is set to N seconds
- **THEN** the system SHALL sample stream data for at least N seconds
- **OR** until `min_keyframes` keyframes are collected (whichever comes first)

#### Scenario: Early termination on keyframes
- **WHEN** during sampling, `min_keyframes` keyframes are collected
- **THEN** the system MAY terminate sampling early
- **AND** use the collected data for metrics calculation

### Requirement: Concurrent Stream Checking

The system SHALL support concurrent checking of multiple streams with configurable limits.

#### Scenario: Concurrent limit enforcement
- **WHEN** the number of active checks exceeds `max_concurrent`
- **THEN** the system SHALL queue additional checks
- **AND** process them as slots become available

#### Scenario: Check interval scheduling
- **WHEN** `check_interval` is set to N seconds
- **THEN** each stream SHALL be checked every N seconds
- **AND** checks SHALL be distributed to avoid simultaneous bursts

### Requirement: Stream Identification

The system SHALL generate unique identifiers for each stream based on project, ID, and URL.

#### Scenario: Stream name generation
- **WHEN** a stream is configured with project="g01", id="D001", url="https://example.com/path/stream.flv"
- **THEN** the system SHALL generate a stream name combining project, host, and path segments
- **AND** use this name as the `name` label in Prometheus metrics

#### Scenario: Label consistency
- **WHEN** metrics are exported for a stream
- **THEN** all metrics for that stream SHALL use consistent labels: `project`, `id`, `name`, `url`

### Requirement: Stream State Management

The system SHALL track stream state transitions and maintain state between checks.

#### Scenario: State persistence
- **WHEN** a stream check completes
- **THEN** the system SHALL store the current state (up/down, metrics)
- **AND** make it available for the next check cycle

#### Scenario: Consecutive failure tracking
- **WHEN** a stream check fails
- **THEN** the system SHALL increment the consecutive failure counter
- **WHEN** a stream check succeeds after failures
- **THEN** the system SHALL reset the consecutive failure counter

### Requirement: HTTP Client Optimization

The system SHALL use an optimized HTTP client for stream connections.

#### Scenario: Connection pooling
- **WHEN** multiple streams are checked
- **THEN** the system SHALL reuse HTTP connections when possible
- **AND** configure connection pool limits (MaxIdleConns, MaxIdleConnsPerHost)

#### Scenario: Timeout handling
- **WHEN** a stream connection times out
- **THEN** the system SHALL handle the timeout gracefully
- **AND** mark the stream as failed for this check cycle

