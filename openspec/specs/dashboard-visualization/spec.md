# Dashboard Visualization

## Purpose

Dashboard Visualization 能力通过 Grafana 仪表板提供流监控数据的可视化，包括实时指标展示、历史趋势分析和多维度过滤功能。

## Requirements

### Requirement: Grafana Dashboard Provisioning

The system SHALL automatically provision Grafana dashboards on startup.

#### Scenario: Dashboard auto-provisioning
- **WHEN** Grafana starts with provisioning configuration
- **THEN** the system SHALL load dashboard JSON from `grafana-provisioning/dashboards/`
- **AND** make it available in Grafana without manual import

#### Scenario: Dashboard configuration
- **WHEN** dashboard provisioning is configured
- **THEN** the system SHALL use `dashboard.yml` to configure:
  - Dashboard provider name
  - Dashboard file path
  - Update interval
  - UI update permissions

### Requirement: Prometheus Datasource Provisioning

The system SHALL automatically configure Prometheus as a datasource in Grafana.

#### Scenario: Datasource auto-configuration
- **WHEN** Grafana starts with datasource provisioning
- **THEN** the system SHALL configure Prometheus datasource from `grafana-provisioning/datasources/prometheus.yml`
- **AND** set it as the default datasource
- **AND** configure connection to Prometheus service

#### Scenario: Datasource UID consistency
- **WHEN** datasource is provisioned
- **THEN** the system SHALL set a consistent UID (e.g., "prometheus")
- **AND** ensure dashboard JSON references the same UID
- **AND** prevent "datasource not found" errors

### Requirement: Dashboard Filtering

The system SHALL provide multi-level filtering in the dashboard.

#### Scenario: Project filter
- **WHEN** dashboard loads
- **THEN** the system SHALL provide a "项目筛选" (project) dropdown
- **AND** populate it with values from `label_values(video_stream_up, project)`
- **AND** allow "All" selection to show all projects

#### Scenario: Table ID filter
- **WHEN** a project is selected
- **THEN** the system SHALL provide a "桌台ID" (id) dropdown
- **AND** populate it with values from `label_values(video_stream_up{project=~"$project"}, id)`
- **AND** allow "All" selection

#### Scenario: Stream name filter
- **WHEN** project and table ID are selected
- **THEN** the system SHALL provide a "流名称" (name) dropdown
- **AND** populate it with values from `label_values(video_stream_up{project=~"$project", id=~"$id"}, name)`
- **AND** allow "All" selection

#### Scenario: Cascading filter behavior
- **WHEN** project filter changes
- **THEN** table ID and stream name filters SHALL update automatically
- **AND** show only values matching the selected project
- **WHEN** table ID filter changes
- **THEN** stream name filter SHALL update automatically

### Requirement: Dashboard Panels

The system SHALL provide comprehensive monitoring panels.

#### Scenario: Summary stat cards
- **WHEN** dashboard loads
- **THEN** the system SHALL display stat cards showing:
  - Total streams count
  - Online streams count
  - Offline streams count
  - Success rate percentage
  - Healthy streams count
  - Playable streams count
  - Average RTT
  - Average packet loss ratio
  - Average network jitter
  - Total reconnects

#### Scenario: Time-series graphs
- **WHEN** dashboard loads
- **THEN** the system SHALL display time-series panels for:
  - Stream status trends (by project)
  - Response time trends
  - Bitrate trends (real-time and average)
  - Framerate trends
  - Packet statistics (total, video, audio)
  - Keyframe trends
  - GOP size trends
  - Quality score trends
  - Stability score trends
  - Network metrics (RTT, packet loss, jitter, reconnects)

#### Scenario: Table view
- **WHEN** dashboard loads
- **THEN** the system SHALL provide a table panel showing:
  - All streams with current status
  - Key metrics per stream
  - Sortable columns

### Requirement: Panel Query Expressions

The system SHALL use PromQL queries that respect filter variables.

#### Scenario: Query with project filter
- **WHEN** a panel queries metrics
- **THEN** the query SHALL include `project=~"$project"` filter
- **AND** work correctly when project is "All" (regex match all)

#### Scenario: Query with table ID filter
- **WHEN** a panel queries metrics
- **THEN** the query SHALL include `id=~"$id"` filter
- **AND** work correctly when id is "All"

#### Scenario: Query with stream name filter
- **WHEN** a panel queries metrics
- **THEN** the query SHALL include `name=~"$name"` filter
- **AND** work correctly when name is "All"

#### Scenario: Counter metric handling
- **WHEN** querying `video_stream_reconnect_total`
- **THEN** the query SHALL use `or vector(0)` to handle cases where counter doesn't exist
- **AND** display 0 for streams without reconnects

### Requirement: Dashboard Layout

The system SHALL organize panels in a logical, readable layout.

#### Scenario: Panel positioning
- **WHEN** dashboard is rendered
- **THEN** stat cards SHALL be positioned at the top
- **AND** time-series graphs SHALL be organized in rows
- **AND** related metrics SHALL be grouped together

#### Scenario: Panel sizing
- **WHEN** panels are displayed
- **THEN** stat cards SHALL use appropriate width (e.g., 6 units per card in 24-unit grid)
- **AND** time-series graphs SHALL span full width or half width as appropriate

### Requirement: Dashboard Refresh

The system SHALL support automatic and manual dashboard refresh.

#### Scenario: Auto-refresh configuration
- **WHEN** dashboard is configured
- **THEN** the system SHALL set an appropriate auto-refresh interval (e.g., 30s)
- **AND** allow users to change it

#### Scenario: Manual refresh
- **WHEN** user clicks refresh button
- **THEN** the system SHALL immediately update all panel data
- **AND** respect current filter selections

### Requirement: Dashboard Time Range

The system SHALL support configurable time ranges.

#### Scenario: Default time range
- **WHEN** dashboard loads
- **THEN** the system SHALL set a default time range (e.g., "Last 1 hour")
- **AND** allow users to change it

#### Scenario: Time range persistence
- **WHEN** user selects a time range
- **THEN** the system SHALL persist the selection
- **AND** restore it on next visit

### Requirement: Dashboard Export/Import

The system SHALL support dashboard JSON export and import.

#### Scenario: Dashboard JSON format
- **WHEN** dashboard is defined
- **THEN** it SHALL be stored as JSON in `grafana-provisioning/dashboards/video-stream-dashboard.json`
- **AND** follow Grafana dashboard JSON schema

#### Scenario: Dashboard versioning
- **WHEN** dashboard is updated
- **THEN** the JSON SHALL include version information
- **AND** maintain backward compatibility where possible

