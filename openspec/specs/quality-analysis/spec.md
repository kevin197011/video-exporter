# Quality Analysis

## Purpose

Quality Analysis 能力负责从采样的流数据中提取和分析质量指标，包括码率、帧率、分辨率、GOP、编码格式等，并计算质量评分和稳定性指标。

## Requirements

### Requirement: Bitrate Calculation

The system SHALL calculate real-time and average bitrate for each stream.

#### Scenario: Real-time bitrate calculation
- **WHEN** stream data is sampled
- **THEN** the system SHALL calculate the current bitrate based on data received in the sampling window
- **AND** export it as `video_stream_bitrate_bps` (in bits per second)

#### Scenario: Average bitrate calculation
- **WHEN** multiple samples have been collected
- **THEN** the system SHALL maintain a rolling average of bitrate
- **AND** export it as `video_stream_avg_bitrate_bps`

#### Scenario: Bitrate stability assessment
- **WHEN** bitrate history is available
- **THEN** the system SHALL calculate bitrate stability (stable/unstable)
- **AND** classify based on variance threshold

### Requirement: Framerate Analysis

The system SHALL calculate the framerate of video streams.

#### Scenario: Framerate calculation
- **WHEN** video packets are received during sampling
- **THEN** the system SHALL calculate frames per second
- **AND** export it as `video_stream_framerate` (in fps)

#### Scenario: Framerate with insufficient data
- **WHEN** sampling duration is too short or no video packets received
- **THEN** the system SHALL set framerate to 0 or N/A
- **AND** indicate insufficient data in logs

### Requirement: Resolution Detection

The system SHALL detect and report video resolution.

#### Scenario: Resolution extraction
- **WHEN** video stream metadata is available
- **THEN** the system SHALL extract width and height
- **AND** store resolution information
- **NOTE**: Resolution is not exported as a separate metric but used in quality assessment

### Requirement: GOP Analysis

The system SHALL analyze Group of Pictures (GOP) structure.

#### Scenario: GOP size calculation
- **WHEN** keyframes are detected in the stream
- **THEN** the system SHALL calculate the number of frames between keyframes
- **AND** export it as `video_stream_gop_size`

#### Scenario: Keyframe counting
- **WHEN** stream data is sampled
- **THEN** the system SHALL count keyframes (I-frames) received
- **AND** export it as `video_stream_keyframes`

### Requirement: Codec Detection

The system SHALL identify the video codec used by the stream.

#### Scenario: Codec identification
- **WHEN** video stream metadata is parsed
- **THEN** the system SHALL identify the codec (e.g., H.264, H.265)
- **AND** store it for quality assessment
- **NOTE**: Codec information is used internally but not exported as a metric

### Requirement: Quality Scoring

The system SHALL calculate quality scores based on multiple factors.

#### Scenario: Quality score calculation
- **WHEN** quality metrics are collected (bitrate, framerate, stability)
- **THEN** the system SHALL calculate a quality score (0-100)
- **AND** export it as `video_stream_quality_score`

#### Scenario: Stability score calculation
- **WHEN** bitrate history is available
- **THEN** the system SHALL calculate a stability score based on variance
- **AND** export it as `video_stream_stability_score`

### Requirement: Playability Assessment

The system SHALL determine if a stream is playable.

#### Scenario: Playable stream detection
- **WHEN** stream check completes
- **THEN** the system SHALL assess playability based on:
  - Presence of video and/or audio packets
  - Keyframe availability
  - Data continuity
- **AND** export `video_stream_playable = 1` if playable, `0` otherwise

#### Scenario: Unplayable stream detection
- **WHEN** no valid video/audio data is received
- **OR** stream is completely down
- **THEN** the system SHALL mark stream as unplayable (`video_stream_playable = 0`)

### Requirement: Health Status Classification

The system SHALL classify stream health status as good, fair, or poor.

#### Scenario: Health classification
- **WHEN** quality metrics are available
- **THEN** the system SHALL classify health as:
  - `good`: High bitrate, stable, playable
  - `fair`: Moderate quality, some issues
  - `poor`: Low quality or unstable
- **AND** export `video_stream_healthy = 1` for good/fair, `0` for poor

### Requirement: Packet Statistics

The system SHALL track packet-level statistics.

#### Scenario: Total packet counting
- **WHEN** stream data is sampled
- **THEN** the system SHALL count total packets received
- **AND** export as `video_stream_total_packets`

#### Scenario: Video packet counting
- **WHEN** video packets are received
- **THEN** the system SHALL count video packets separately
- **AND** export as `video_stream_video_packets`

#### Scenario: Audio packet counting
- **WHEN** audio packets are received
- **THEN** the system SHALL count audio packets separately
- **AND** export as `video_stream_audio_packets`

### Requirement: Response Time Measurement

The system SHALL measure HTTP response time for FLV streams.

#### Scenario: Response time for HTTP-FLV
- **WHEN** connecting to an HTTP-FLV stream
- **THEN** the system SHALL measure the time from request to first data received
- **AND** export it as `video_stream_response_ms` (in milliseconds)

#### Scenario: Response time for non-HTTP streams
- **WHEN** connecting to a non-HTTP stream (RTMP, RTSP, etc.)
- **THEN** the system SHALL set response time to 0 or N/A
- **AND** indicate in logs that response time is not applicable

