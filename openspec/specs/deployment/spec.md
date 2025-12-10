# Deployment

## Purpose

Deployment 能力涵盖容器化部署、Docker Compose 编排、CI/CD 集成和部署脚本，提供一键部署和运维支持。

## Requirements

### Requirement: Docker Image Building

The system SHALL provide a Dockerfile for building container images.

#### Scenario: Multi-stage build
- **WHEN** building Docker image
- **THEN** the system SHALL use multi-stage build:
  - Builder stage: Compile Go binary with dependencies
  - Runtime stage: Minimal Alpine image with FFmpeg and binary

#### Scenario: Build optimization
- **WHEN** building image
- **THEN** the system SHALL:
  - Cache Go modules layer
  - Copy source code after module download
  - Build static binary (CGO_ENABLED=0)
  - Use non-root user in runtime

#### Scenario: Runtime dependencies
- **WHEN** runtime image is created
- **THEN** it SHALL include:
  - FFmpeg (for stream processing)
  - CA certificates (for HTTPS)
  - Video exporter binary

### Requirement: Docker Compose Orchestration

The system SHALL provide Docker Compose configuration for multi-service deployment.

#### Scenario: Service definition
- **WHEN** Docker Compose starts
- **THEN** it SHALL define services:
  - `video-exporter`: Main monitoring service
  - `prometheus`: Metrics collection
  - `grafana`: Visualization

#### Scenario: Network configuration
- **WHEN** services start
- **THEN** they SHALL be on the same Docker network
- **AND** services SHALL communicate using service names
- **AND** network SHALL be named appropriately (e.g., "monitoring")

#### Scenario: Volume mounts
- **WHEN** services start
- **THEN** the system SHALL mount:
  - Config file for video-exporter
  - Prometheus data directory
  - Grafana provisioning directories
  - Grafana data directory

#### Scenario: Environment variables
- **WHEN** services start
- **THEN** the system SHALL configure:
  - Grafana provisioning path
  - Service ports
  - Data persistence paths

### Requirement: Prometheus Configuration

The system SHALL provide Prometheus configuration for scraping metrics.

#### Scenario: Scrape configuration
- **WHEN** Prometheus starts
- **THEN** it SHALL be configured to scrape `video-exporter:8080/metrics`
- **AND** use appropriate scrape interval (e.g., 30s)

#### Scenario: Target discovery
- **WHEN** Prometheus runs
- **THEN** it SHALL use static_configs to discover video-exporter
- **AND** use service name for DNS resolution in Docker network

### Requirement: Grafana Provisioning

The system SHALL automatically provision Grafana on startup.

#### Scenario: Provisioning directory structure
- **WHEN** Grafana starts
- **THEN** provisioning directories SHALL exist:
  - `grafana-provisioning/datasources/`
  - `grafana-provisioning/dashboards/`
  - `grafana-provisioning/plugins/` (with .gitkeep)
  - `grafana-provisioning/alerting/` (with .gitkeep)
  - `grafana-provisioning/notifiers/` (with .gitkeep)

#### Scenario: Provisioning environment variable
- **WHEN** Grafana starts
- **THEN** `GF_PATHS_PROVISIONING` SHALL be set to `/etc/grafana/provisioning`
- **AND** Grafana SHALL load provisioning files from this path

### Requirement: Deployment Scripts

The system SHALL provide scripts for common deployment operations.

#### Scenario: Start script
- **WHEN** `./scripts/start.sh` is executed
- **THEN** it SHALL:
  - Start Docker Compose services
  - Wait for services to be healthy
  - Display service URLs
  - Show status

#### Scenario: Stop script
- **WHEN** `./scripts/stop.sh` is executed
- **THEN** it SHALL:
  - Stop Docker Compose services
  - Optionally remove containers
  - Clean up resources

#### Scenario: Logs script
- **WHEN** `./scripts/logs.sh [service]` is executed
- **THEN** it SHALL:
  - Show logs for specified service or all services
  - Follow log output (tail -f)
  - Support service name filtering

### Requirement: CI/CD Integration

The system SHALL provide GitHub Actions workflow for automated builds.

#### Scenario: Docker image build
- **WHEN** code is pushed or PR is created
- **THEN** GitHub Actions SHALL:
  - Build Docker image
  - Use Dockerfile from project root
  - Enable build cache for faster builds
  - Tag image appropriately

#### Scenario: Image publishing
- **WHEN** build succeeds
- **THEN** GitHub Actions SHALL:
  - Push image to container registry
  - Tag with commit SHA and branch name
  - Use docker/build-push-action
  - Use docker/metadata-action for tags

#### Scenario: Build cache
- **WHEN** building image
- **THEN** the workflow SHALL:
  - Use cache-from for faster builds
  - Use cache-to for cache persistence
  - Cache Go modules layer

### Requirement: Health Checks

The system SHALL provide health check endpoints and configurations.

#### Scenario: Video exporter health check
- **WHEN** Docker container runs
- **THEN** it SHALL include HEALTHCHECK instruction
- **AND** check `/metrics` endpoint
- **AND** use appropriate intervals and timeouts

#### Scenario: Service health verification
- **WHEN** services start
- **THEN** scripts SHALL verify:
  - Video exporter responds on `/metrics`
  - Prometheus scrapes metrics successfully
  - Grafana is accessible and dashboards load

### Requirement: Configuration Management

The system SHALL support configuration via files and environment variables.

#### Scenario: Config file mounting
- **WHEN** Docker Compose starts
- **THEN** `config.yml` SHALL be mounted into container
- **AND** video-exporter SHALL load it on startup

#### Scenario: Config file path
- **WHEN** video-exporter starts
- **THEN** it SHALL look for config at:
  - Path specified in `CONFIG_FILE` environment variable
  - Default path: `/app/config.yml` or `./config.yml`

### Requirement: Data Persistence

The system SHALL persist data across container restarts.

#### Scenario: Prometheus data persistence
- **WHEN** Prometheus runs
- **THEN** data directory SHALL be mounted as volume
- **AND** metrics SHALL persist across restarts

#### Scenario: Grafana data persistence
- **WHEN** Grafana runs
- **THEN** data directory SHALL be mounted as volume
- **AND** dashboards and settings SHALL persist

### Requirement: Port Exposure

The system SHALL expose necessary ports for external access.

#### Scenario: Service ports
- **WHEN** services start
- **THEN** ports SHALL be exposed:
  - Video exporter: 8080 (metrics)
  - Prometheus: 9090 (UI and API)
  - Grafana: 3000 (UI)

#### Scenario: Port mapping
- **WHEN** Docker Compose starts
- **THEN** ports SHALL be mapped to host
- **AND** allow external access to services

### Requirement: Resource Limits

The system SHALL support resource limit configuration.

#### Scenario: Resource constraints
- **WHEN** services are defined
- **THEN** Docker Compose MAY include:
  - Memory limits
  - CPU limits
  - Restart policies

#### Scenario: Resource monitoring
- **WHEN** services run
- **THEN** resource usage SHALL be monitorable
- **AND** logs SHALL indicate resource constraints if hit

