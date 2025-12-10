# Network Monitoring

## Purpose

Network Monitoring 能力负责收集和分析网络层面的指标，包括 RTT（往返时间）、丢包率、网络抖动和重连统计，用于评估流的网络稳定性。

## Requirements

### Requirement: RTT Measurement

The system SHALL measure Round-Trip Time (RTT) for stream connections.

#### Scenario: RTT calculation
- **WHEN** establishing a connection to a stream
- **THEN** the system SHALL measure the time from connection initiation to first data received
- **AND** export it as `video_stream_rtt_ms` (in milliseconds)

#### Scenario: RTT update on reconnection
- **WHEN** a stream reconnects after a failure
- **THEN** the system SHALL measure and update RTT for the new connection
- **AND** replace the previous RTT value

#### Scenario: RTT reset on failure
- **WHEN** a stream check fails
- **THEN** the system SHALL reset RTT to 0 or leave it as last known value
- **AND** clear it on successful reconnection

### Requirement: Packet Loss Ratio Calculation

The system SHALL calculate packet loss ratio based on expected vs received packets.

#### Scenario: Packet loss calculation
- **WHEN** stream data is sampled
- **THEN** the system SHALL estimate packet loss based on:
  - Sequence number gaps (if available)
  - Expected packet rate vs actual rate
  - Connection errors
- **AND** export it as `video_stream_packet_loss_ratio` (0.0-1.0, where 1.0 = 100% loss)

#### Scenario: No packet loss
- **WHEN** all expected packets are received
- **THEN** the system SHALL set packet loss ratio to 0.0
- **AND** export the metric

#### Scenario: Packet loss reset
- **WHEN** a stream check fails
- **THEN** the system SHALL reset packet loss ratio to 0.0
- **AND** recalculate on next successful check

### Requirement: Network Jitter Measurement

The system SHALL measure network jitter as the standard deviation of packet inter-arrival times.

#### Scenario: Jitter calculation
- **WHEN** multiple packets are received during sampling
- **THEN** the system SHALL calculate inter-arrival time between consecutive packets
- **AND** compute standard deviation of inter-arrival times
- **AND** export it as `video_stream_network_jitter_ms` (in milliseconds)

#### Scenario: Insufficient data for jitter
- **WHEN** fewer than 2 packets are received
- **THEN** the system SHALL set jitter to 0 or N/A
- **AND** indicate insufficient data

#### Scenario: Jitter reset
- **WHEN** a stream check fails
- **THEN** the system SHALL reset jitter to 0
- **AND** recalculate on next successful check

### Requirement: Reconnect Tracking

The system SHALL track reconnection events for streams.

#### Scenario: Reconnect detection
- **WHEN** a stream transitions from up to down
- **AND** then back to up
- **THEN** the system SHALL detect this as a reconnect event
- **AND** increment the reconnect counter

#### Scenario: Reconnect counter as Gauge
- **WHEN** reconnects occur within a check cycle
- **THEN** the system SHALL track the count for the current cycle
- **AND** export it as `video_stream_reconnect_total` (Gauge type, not Counter)
- **NOTE**: This is a Gauge that represents reconnects in the current period, not a cumulative counter

#### Scenario: Reconnect counter reset
- **WHEN** a new check cycle begins
- **THEN** the system SHALL reset the reconnect counter for that cycle
- **AND** start counting from 0

#### Scenario: Counter behavior for zero reconnects
- **WHEN** a stream has no reconnects in a cycle
- **THEN** the metric MAY not be exported (Prometheus Counter/Gauge behavior)
- **AND** queries using `or vector(0)` SHALL display 0 for streams without reconnects

### Requirement: Network Metrics Integration

The system SHALL integrate network metrics with stream health assessment.

#### Scenario: Network metrics in health check
- **WHEN** network metrics are collected
- **THEN** the system SHALL use them as factors in health assessment
- **AND** consider high RTT, packet loss, or jitter as indicators of poor network conditions

#### Scenario: Network metrics logging
- **WHEN** network metrics are calculated
- **THEN** the system SHALL log RTT, packet loss ratio, and jitter values
- **AND** include them in structured log output

### Requirement: Network Metrics Persistence

The system SHALL maintain network metrics between check cycles.

#### Scenario: Metrics persistence
- **WHEN** a check cycle completes
- **THEN** the system SHALL store current network metrics (RTT, packet loss, jitter)
- **AND** make them available for the next cycle

#### Scenario: Metrics reset on failure
- **WHEN** a stream check fails
- **THEN** the system SHALL reset network metrics (RTT, packet loss, jitter) to 0
- **AND** reset reconnect count
- **AND** recalculate on successful reconnection

