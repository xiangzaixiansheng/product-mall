# 升级文档 (2026-05-25)

## 2026-05-25 安全加固 & 工程基建完善

### 安全修复

| 项目 | 说明 |
|------|------|
| Dockerfile 版本同步 | `golang:1.22-alpine` → `golang:1.25-alpine`，与 go.mod 一致 |
| Session 密钥 | 硬编码 → 读取 `SESSION_*` 环境变量 |
| JWT 认证标准化 | 从 Cookie header 改为标准 `Authorization: Bearer <token>` |
| Dev 环境认证 | 移除 dev 环境跳过认证的逻辑，统一走 token 校验 |
| Handler 安全 | API handler 不再重复解析 token，从 gin.Context 中获取 user_id |

### 工程基建

| 文件 | 说明 |
|------|------|
| `Makefile` | 统一 build/run/test/lint/docker/swagger 命令 |
| `.golangci.yml` | golangci-lint 配置，启用 gosec/errcheck/staticcheck 等 |
| `.github/workflows/ci.yml` | GitHub Actions CI 流水线 (build + test + lint) |
| `docker-compose.yml` | MySQL + Redis + App 一键编排启动 |

### 代码质量

| 改动 | 说明 |
|------|------|
| `LogrusObj` → `Logger` | 全局重命名，消除 logrus 历史遗留命名 |
| `dto.Response` 增强 | 字段 `omitempty` 优化；新增 `BuildPagedResponse` 分页支持 |
| 优雅关闭增强 | shutdown 时正确关闭 MySQL/Redis 连接 |
| Server 超时 | 新增 ReadTimeout/WriteTimeout/IdleTimeout 配置 |

### 功能增强

| 功能 | 说明 |
|------|------|
| Swagger API 文档 | 集成 swaggo，访问 `/swagger/index.html` 查看文档 |
| 分页响应 | `dto.BuildPagedResponse` 支持标准分页返回格式 |

### 变更文件清单

- `Dockerfile` / `Dockerfile.multistage` — Go 版本更新
- `cmd/main.go` — 优雅关闭、Server 超时、Swagger 注解
- `internal/routes/router.go` — Session 环境变量、Swagger 路由
- `internal/middleware/jwt.go` — Authorization header 认证重构
- `internal/api/v1/*.go` — 使用 context 获取用户ID、Swagger 注解
- `internal/dto/common.go` — 分页响应结构
- `pkg/pkg_logger/*.go` — Logger 重命名
- `.env.example` — 新增 SESSION_* 配置项

---

## 2026-05-25 企业级功能增强

### 新增中间件

| 中间件 | 文件 | 说明 |
|-------|------|------|
| 超时控制 | `internal/middleware/timeout.go` | 30秒请求超时，context timeout |
| Gzip压缩 | `internal/middleware/gzip.go` | HTTP 响应压缩 |
| 安全头 | `internal/middleware/security.go` | X-Frame-Options, CSP, HSTS 等 |
| 限流 | `internal/middleware/ratelimit.go` | Redis 滑动窗口限流 (100req/min) |

### 新增功能

| 功能 | 文件 | 说明 |
|------|------|------|
| 健康检查 | `internal/api/v1/health.go` | `/health` 检查 MySQL/Redis 状态 |
| 统一响应 | `pkg/response/response.go` | 标准化的 API 响应格式 |
| 指标收集中间件 | `internal/middleware/metrics.go` | 请求计数和延迟统计 |

### 依赖更新

| 包 | 版本 | 说明 |
|---|------|------|
| gin | v1.9.1 → v1.12.0 | Web 框架升级 |
| gin-contrib/gzip | 新增 v1.2.6 | 压缩中间件 |
| validator | v10 → v11 | 参数校验 |

### 路由更新

- 新增 `GET /health` 健康检查接口
- 所有中间件默认启用 (压缩、安全头、超时、限流)
- 限流使用 Redis 滑动窗口算法

---

## 初始升级 (2026-05-25 上午)

## 升级概述

本次升级将 `product-mall` 项目从 Go 1.17 升级至 Go 1.22+，并更新了所有过时的依赖和代码模式。

## 升级内容

### 1. Go 版本升级

| 项目 | 原来 | 现在 |
|------|------|------|
| Go 版本 | 1.17 | 1.24 (latest) |
| go.mod 声明 | `go 1.17` | `go 1.24` |

### 2. 依赖升级

| 包 | 原来 | 现在 | 说明 |
|----|------|------|------|
| redis | go-redis/redis/v8 v8.11.5 | redis/go-redis/v9 v9.19.0 | 官方维护的 Redis 客户端 |
| redis mock | go-redis/redismock/v8 | go-redis/redismock/v9 v9.2.0 | |
| RabbitMQ | streadway/amqp v1.1.0 | rabbitmq/amqp091-go v1.11.0 | streadway 已归档，amqp091-go 是官方推荐的替代 |
| JWT | golang-jwt/jwt/v4 v4.5.0 | golang-jwt/jwt/v5 v5.3.1 | |
| GORM | jinzhu/gorm v1.9.16 | 移除 | 使用 gorm.io/gorm |
| logrus | sirupsen/logrus v1.9.3 | 移除 | 使用标准库 slog |

### 3. 代码模式现代化

#### 3.1 废弃的 stdlib 使用
- `ioutil.ReadFile` → `os.ReadFile`
- `ioutil.ReadAll` → `io.ReadAll`

#### 3.2 类型声明
- `interface{}` → `any` (Go 1.18+)

#### 3.3 错误处理
- `errors.New(fmt.Sprintf(...))` → `fmt.Errorf("...: %w", err)`
- 使用 `errors.Is()` 进行错误检查

#### 3.4 随机数
- 删除 `rand.Seed()` (Go 1.20+ 自动种子)

#### 3.5 日志系统
- logrus → slog (标准库 `log/slog`)
- 更新所有日志调用为 `slog.Logger` 方法

### 4. JWT 升级 (v4 → v5)

#### 变化
- `jwt.StandardClaims` → `jwt.RegisteredClaims`
- `ExpiresAt: time.Unix()` → `ExpiresAt: jwt.NewNumericDate(time)`

#### JWT 密钥配置
- 从硬编码改为环境变量 `JWT_SECRET_KEY`
- 支持通过环境变量动态配置

### 5. 服务改进

#### 5.1 优雅关停
新增 `os/signal` 信号处理和 `http.Server.Shutdown()`：
- 支持 SIGINT/SIGTERM 信号
- 10秒超时等待处理中的请求
- 正确的资源清理

#### 5.2 日志结构化
使用 slog JSONHandler 输出结构化日志到 `logs/` 目录

### 6. 基础设施更新

#### Dockerfile
- `golang:1.16-alpine` → `golang:1.22-alpine`

### 7. 新增文件

- `.env.example` - 环境变量配置示例
- `docs/CHANGELOG.md` - 本文档

## 变更文件清单

### 配置文件
- `go.mod` / `go.sum` - 依赖更新
- `Dockerfile` / `Dockerfile.multistage` - Go 版本更新
- `README.md` - 文档更新

### 核心代码
- `cmd/main.go` - 添加优雅关停
- `conf/conf.go` - slog 日志配置
- `conf/i18n.go` - 使用 os.ReadFile, any 类型

### Model 层
- `internal/model/*.go` - 更新 gorm import

### Service 层
- `internal/service/*.go` - 更新日志调用

### 工具类
- `internal/tools/jwt.go` - JWT v5 迁移
- `internal/tools/str.go` - 使用 slog, any 类型
- `internal/tools/curl.go` - 使用 io.ReadAll

### 外部包
- `pkg/db/redis_client.go` - Redis v9 迁移
- `pkg/pkg_logger/*.go` - slog 迁移
- `pkg/rabbitMQ/rabbitMQPool.go` - amqp091-go, 错误包装

## 验证

```bash
go build ./...     # 编译通过
go vet ./...       # 代码检查通过
go test ./...      # 测试通过
```

## 回滚注意事项

如果需要回滚，请恢复以下文件：
- `go.mod` / `go.sum` - 恢复旧依赖
- `Dockerfile*` - 恢复旧 Go 版本
- `internal/tools/jwt.go` - 恢复 JWT v4 语法

## 参考链接

- [Go 1.22 Release Notes](https://go.dev/doc/go1.22)
- [Redis Go Client Migration](https://github.com/redis/go-redis)
- [RabbitMQ AMQP091 Go](https://github.com/rabbitmq/amqp091-go)
- [JWT v5 Migration Guide](https://github.com/golang-jwt/jwt/blob/master/MIGRATION_V5.md)
- [Go slog Package](https://pkg.go.dev/log/slog)