# 项目结构说明

本文档详细说明 Video Exporter 项目的目录结构和文件组织方式。

## 目录结构

```
video-exporter/
├── main.go                       # 程序入口
├── config.go                     # 配置模块
├── logger.go                     # 日志模块
├── exporter.go                   # Prometheus 导出器
├── scheduler.go                  # 任务调度器
├── stream.go                     # 流检查核心逻辑
│
├── config.yml                    # 配置文件（不提交）
├── config.example.yaml           # 配置示例（提交）
│
├── deployments/                  # 部署配置
│   └── grafana/                  # Grafana 配置
│       └── grafana-provisioning/ # 自动配置
│           ├── datasources/      # 数据源配置
│           │   └── prometheus.yml
│           └── dashboards/       # 仪表板配置
│               ├── dashboard.yml
│               └── video-stream-dashboard.json
│
├── docker-compose.yml            # Docker Compose 编排配置
├── Dockerfile                    # Docker 镜像构建
├── prometheus.yml                # Prometheus 配置
│
├── scripts/                      # 脚本工具
│   ├── start.sh                 # 启动脚本
│   ├── stop.sh                  # 停止脚本
│   └── logs.sh                  # 日志查看脚本
│
├── docs/                         # 文档
│   ├── API.md                   # API 文档
│   ├── PROJECT-STRUCTURE.md     # 项目结构说明（本文件）
│   ├── DOCKER-COMPOSE-README.md # Docker Compose 使用文档
│   └── DEPLOYMENT-CHECKLIST.md  # 部署检查清单
│
├── README.md                     # 项目说明
├── CONTRIBUTING.md               # 贡献指南
├── CHANGELOG.md                  # 更新日志
├── LICENSE                       # 许可证
├── .gitignore                    # Git 忽略规则
│
├── Makefile                      # 构建脚本
├── Rakefile                      # Ruby 任务脚本
├── push.rb                       # 推送脚本
│
├── go.mod                        # Go 模块定义
└── go.sum                        # Go 模块校验
```

## 文件说明

### 核心代码

| 文件 | 说明 | 职责 |
|------|------|------|
| `main.go` | 程序入口 | 初始化配置、日志、启动服务 |
| `config.go` | 配置管理 | 配置文件加载、结构体定义 |
| `logger.go` | 日志系统 | 结构化日志、日志级别控制 |
| `exporter.go` | Prometheus 导出器 | 指标定义、采集、导出 |
| `scheduler.go` | 任务调度器 | 并发控制、定时任务 |
| `stream.go` | 流检查逻辑 | 流连接、数据采集、质量分析 |

### 配置文件

| 文件 | 说明 | 版本控制 |
|------|------|----------|
| `config.yml` | 运行时配置 | ❌ 不提交（包含敏感信息） |
| `config.example.yaml` | 配置示例 | ✅ 提交（用于参考） |

### 部署目录 (`deployments/`)

#### Docker 部署 (`deployments/docker/`)

| 文件 | 说明 |
|------|------|
| `Dockerfile` | 多阶段构建，优化镜像大小 |
| `docker-compose.yml` | 编排三个服务：exporter、prometheus、grafana |
| `prometheus.yml` | Prometheus 抓取配置 |

#### Grafana 配置 (`deployments/grafana/`)

| 文件/目录 | 说明 |
|-----------|------|
| `grafana-provisioning/dashboards/video-stream-dashboard.json` | 仪表板定义（面板、查询、样式） |
| `grafana-provisioning/datasources/` | 数据源自动配置 |
| `grafana-provisioning/dashboards/` | 仪表板自动导入 |

### 脚本目录 (`scripts/`)

| 脚本 | 用途 | 使用方式 |
|------|------|----------|
| `start.sh` | 启动所有服务 | `./scripts/start.sh` |
| `stop.sh` | 停止所有服务 | `./scripts/stop.sh` |
| `logs.sh` | 查看服务日志 | `./scripts/logs.sh [service]` |

### 文档目录 (`docs/`)

| 文档 | 内容 |
|------|------|
| `API.md` | Prometheus 指标定义、PromQL 查询示例 |
| `PROJECT-STRUCTURE.md` | 项目结构说明（本文件） |
| `DOCKER-COMPOSE-README.md` | Docker Compose 详细使用说明 |
| `DEPLOYMENT-CHECKLIST.md` | 部署前检查清单 |

### 根目录文档

| 文档 | 说明 |
|------|------|
| `README.md` | 项目概述、快速开始 |
| `CONTRIBUTING.md` | 贡献指南、代码规范 |
| `CHANGELOG.md` | 版本更新记录 |
| `LICENSE` | 开源许可证 |

## 设计原则

### 1. 简洁性

- 核心代码放在根目录，便于快速定位
- 避免过度嵌套，保持结构扁平
- 单文件单职责，模块划分清晰

### 2. 标准化

- 遵循 Go 项目标准布局
- 符合 Cursor 开发规范
- 采用社区最佳实践

### 3. 可维护性

- 配置与代码分离
- 部署脚本集中管理
- 文档完善，便于上手

### 4. 扩展性

预留扩展目录（未来可添加）：
- `internal/` - 私有包
- `pkg/` - 公共库
- `test/` - 测试代码
- `api/` - API 定义

## 工作流程

### 开发流程

```
1. 修改代码 (*.go)
   ↓
2. 本地测试 (go run *.go)
   ↓
3. 构建镜像 (docker build)
   ↓
4. 测试部署 (docker-compose up)
   ↓
5. 提交代码 (git commit)
```

### 部署流程

```
1. 准备配置 (config.yml)
   ↓
2. 启动服务 (./scripts/start.sh)
   ↓
3. 验证服务 (访问 metrics)
   ↓
4. 配置监控 (Grafana 仪表板)
   ↓
5. 监控运行
```

## 配置管理

### 配置优先级

1. 环境变量（最高优先级）
2. 配置文件 `config.yml`
3. 默认值（代码中定义）

### 敏感信息处理

- ❌ 不提交 `config.yml`
- ✅ 提交 `config.example.yaml`
- ✅ 使用 `.gitignore` 排除

## 依赖管理

### Go 模块

- `go.mod` - 模块定义
- `go.sum` - 依赖校验

### 外部依赖

- joy5 - FLV 解析
- prometheus/client_golang - Prometheus 客户端
- slog - 结构化日志

## 最佳实践

### 1. 文件命名

- 使用小写字母
- 单词间用下划线分隔（如 `docker-compose.yml`）
- Go 文件使用下划线（如 `stream_test.go`）

### 2. 目录命名

- 使用小写字母
- 复数形式（如 `scripts/`、`docs/`）
- 语义清晰（如 `deployments/`）

### 3. 文档维护

- 代码变更同步更新文档
- 重要变更记录在 `CHANGELOG.md`
- API 变更更新 `docs/API.md`

### 4. 版本管理

- 遵循语义化版本（SemVer）
- 使用 Git tags 标记版本
- 保持 `CHANGELOG.md` 更新

## 常见问题

### Q: 为什么不使用 `internal/` 目录？

A: 当前项目规模较小，不需要复杂的包结构。如果项目扩大，可以考虑引入。

### Q: 配置文件为什么不放在 `configs/` 目录？

A: `config.yml` 是运行时配置，放在根目录更方便。示例配置 `config.example.yaml` 也在根目录便于参考。

### Q: 测试文件应该放在哪里？

A: 测试文件放在对应的代码文件旁边，如 `stream_test.go` 测试 `stream.go`。

### Q: 如何添加新的部署方式（如 Kubernetes）？

A: 在 `deployments/` 下创建新目录，如 `deployments/kubernetes/`。

## 参考资料

- [Go 项目标准布局](https://github.com/golang-standards/project-layout)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)

