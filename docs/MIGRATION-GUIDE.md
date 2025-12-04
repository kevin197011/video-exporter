# 项目结构重构迁移指南

本文档说明项目结构重构后的变化，帮助您快速适应新结构。

## 变更概览

### 文件移动

| 原位置 | 新位置 | 说明 |
|--------|--------|------|
| `./grafana-dashboard.json` | `deployments/grafana/grafana-provisioning/dashboards/video-stream-dashboard.json` | Grafana 仪表板 |

**注意**：Docker 相关文件（`Dockerfile`, `docker-compose.yml`, `prometheus.yml`）保留在根目录，便于访问。
| `./start.sh` | `scripts/start.sh` | 启动脚本 |
| `./stop.sh` | `scripts/stop.sh` | 停止脚本 |
| `./DOCKER-COMPOSE-README.md` | `docs/DOCKER-COMPOSE-README.md` | Docker Compose 文档 |
| `./DEPLOYMENT-CHECKLIST.md` | `docs/DEPLOYMENT-CHECKLIST.md` | 部署检查清单 |

### 新增文件

| 文件 | 说明 |
|------|------|
| `CONTRIBUTING.md` | 贡献指南 |
| `CHANGELOG.md` | 更新日志 |
| `docs/API.md` | API 文档 |
| `docs/PROJECT-STRUCTURE.md` | 项目结构说明 |
| `docs/MIGRATION-GUIDE.md` | 迁移指南（本文件） |
| `scripts/logs.sh` | 日志查看脚本 |

## 使用方式变更

### Docker Compose

**原来的方式**：
```bash
# 启动
docker-compose up -d

# 停止
docker-compose stop

# 查看日志
docker-compose logs -f
```

**现在的方式**：
```bash
# 启动（推荐使用脚本）
./scripts/start.sh
# 或直接使用 docker-compose
docker-compose up -d

# 停止
./scripts/stop.sh
# 或
docker-compose stop

# 查看日志
./scripts/logs.sh
# 或
docker-compose logs -f
```

### Docker Build

**原来的方式**：
```bash
docker build -t video-exporter .
```

**现在的方式**：
```bash
docker build -t video-exporter .
```

### 文档查阅

**原来的方式**：
- 所有文档在根目录

**现在的方式**：
- 主要文档：根目录（`README.md`, `CONTRIBUTING.md`, `CHANGELOG.md`）
- 详细文档：`docs/` 目录
  - API 文档：`docs/API.md`
  - 部署文档：`docs/DOCKER-COMPOSE-README.md`
  - 项目结构：`docs/PROJECT-STRUCTURE.md`

## 迁移步骤

### 1. 更新本地仓库

```bash
# 拉取最新代码
git pull origin main

# 查看新结构
tree -L 2 -I 'vendor|.git'
```

### 2. 更新 Docker Compose 使用方式

如果您之前使用 `docker-compose` 命令：

```bash
# 停止旧的容器
docker-compose down

# 使用 docker-compose 启动
docker-compose up -d

# 或使用脚本（推荐）
./scripts/start.sh
```

### 3. 更新 CI/CD 配置

如果您有 CI/CD 流程，需要更新相关路径：

**Docker Build**：
```yaml
# 都是一样的（文件在根目录）
docker build -t video-exporter .
```

**Docker Compose**：
```yaml
# 都是一样的（文件在根目录）
docker-compose up -d
```

### 4. 更新脚本引用

如果您有自定义脚本引用了项目文件：

```bash
# 原来
./start.sh

# 现在
./scripts/start.sh
```

## 兼容性说明

### 向后兼容

- ✅ Go 代码无变化，编译方式不变
- ✅ 配置文件 `config.yml` 格式不变
- ✅ Prometheus 指标不变
- ✅ Grafana 仪表板功能不变

### 不兼容变更

- ❌ Docker Compose 文件路径变更
- ❌ 脚本路径变更
- ❌ 部分文档路径变更

## 优势

### 1. 更清晰的结构

- 部署文件集中在 `deployments/`
- 脚本集中在 `scripts/`
- 文档集中在 `docs/`

### 2. 更好的可维护性

- 符合 Go 项目标准布局
- 遵循 Cursor 开发规范
- 便于团队协作

### 3. 更完善的文档

- 新增贡献指南
- 新增更新日志
- 新增 API 文档
- 新增项目结构说明

### 4. 更方便的使用

- 统一的启动脚本
- 统一的日志查看
- 统一的停止脚本

## 常见问题

### Q: 为什么要重构项目结构？

A: 为了遵循 Go 项目最佳实践和 Cursor 开发规范，提高代码可维护性和团队协作效率。

### Q: 旧的命令还能用吗？

A: 大部分 Go 相关命令不变。Docker Compose 命令需要指定新路径，或使用新的脚本。

### Q: 需要重新构建 Docker 镜像吗？

A: 是的，因为 Dockerfile 路径变了。但镜像内容和功能完全相同。

### Q: Grafana 仪表板需要重新导入吗？

A: 不需要。如果使用 Docker Compose 启动，仪表板会自动导入。

### Q: 配置文件需要修改吗？

A: 不需要。`config.yml` 格式和内容完全不变。

## 回滚方案

如果遇到问题需要回滚到旧版本：

```bash
# 查看提交历史
git log --oneline

# 回滚到重构前的版本
git checkout <commit-hash>

# 或创建新分支
git checkout -b old-structure <commit-hash>
```

## 获取帮助

如果您在迁移过程中遇到问题：

1. 查看 [项目结构文档](PROJECT-STRUCTURE.md)
2. 查看 [部署检查清单](DEPLOYMENT-CHECKLIST.md)
3. 提交 Issue 询问

## 反馈

如果您对新结构有任何建议或意见，欢迎：

1. 提交 Issue
2. 提交 Pull Request
3. 在讨论区留言

---

感谢您的理解和支持！新结构将为项目带来更好的可维护性和扩展性。

