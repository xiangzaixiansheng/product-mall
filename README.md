# product-mall

基于 Go 1.22+ 的电商学习项目，使用 Gin 框架构建 RESTful API。

## 技术栈

- **Go**: 1.22+
- **Web 框架**: Gin v1.9.1
- **ORM**: GORM v1.25.5
- **数据库**: MySQL (读写分离)
- **缓存**: Redis v9 (go-redis)
- **消息队列**: RabbitMQ (amqp091-go)
- **对象存储**: 七牛云 SDK
- **日志**: slog (标准库)
- **认证**: JWT v5

## 项目结构

```
product-mall/
├── cmd/                    # 应用入口
│   └── main.go
├── conf/                   # 配置文件
│   ├── app.dev.ini         # 开发环境配置
│   ├── app.prod.ini        # 生产环境配置
│   ├── conf.go             # 配置加载器
│   └── i18n.go            # 国际化
├── internal/               # 内部包
│   ├── api/v1/            # API 处理器
│   ├── dto/               # 数据传输对象
│   ├── middleware/        # 中间件 (JWT, CORS, Logger)
│   ├── model/             # 数据模型
│   ├── repo/mysql/        # MySQL 仓库
│   ├── routes/           # 路由定义
│   └── service/           # 业务逻辑
├── pkg/                   # 公共包
│   ├── db/                # 数据库连接
│   ├── pkg_logger/        # 日志封装
│   └── rabbitMQ/          # RabbitMQ 连接池
├── cache/                 # Redis 缓存封装
└── Dockerfile             # Docker 配置
```

## 快速启动

### 1. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

主要配置项：
- `ENV`: 运行模式 (dev/prod)
- `JWT_SECRET_KEY`: JWT 密钥
- `HTTP_PORT`: 服务端口

### 2. 配置数据库

修改 `conf/app.dev.ini` (开发环境) 或 `conf/app.prod.ini` (生产环境)：

```ini
[mysql]
DbHost = localhost
DbPort = 3306
DbUser = root
DbPassword = your-password
DbName = product_mall

[redis]
RedisAddr = localhost:6379
RedisDbName = 0
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 启动服务

```bash
# 开发环境
go run cmd/main.go

# 生产环境
ENV=prod go run cmd/main.go
```

服务启动后会监听 `conf` 中配置的端口 (默认 3000)。

## API 文档

### 用户接口
- `POST /api/v1/user/register` - 用户注册
- `POST /api/v1/user/login` - 用户登录
- `POST /api/v1/user/update` - 更新用户信息
- `POST /api/v1/user/send_email` - 发送邮箱验证
- `GET /api/v1/user/valid_email` - 验证邮箱

### 商品接口 (需认证)
- `POST /api/v1/product/create` - 创建商品

### 购物车接口 (需认证)
- `GET /api/v1/cart/list/:id` - 获取购物车列表
- `POST /api/v1/cart/create` - 添加商品到购物车
- `POST /api/v1/cart/update` - 更新购物车
- `POST /api/v1/cart/delete` - 删除购物车商品

### 地址接口 (需认证)
- `GET /api/v1/address/list/:id` - 获取地址列表
- `POST /api/v1/address/create` - 创建地址
- `POST /api/v1/address/update` - 更新地址
- `POST /api/v1/address/delete` - 删除地址

## 开发

### 代码格式化

```bash
go fmt ./...
go vet ./...
```

### 运行测试

```bash
go test ./...
```

### Docker 构建

```bash
# 单阶段构建
docker build -t product-mall .

# 多阶段构建 (更小镜像)
docker build -f Dockerfile.multistage -t product-mall:multistage .
```

### Docker 运行

```bash
docker run -d -p 3000:3000 --network host product-mall
```

## 主要依赖

| 包 | 说明 |
|---|---|
| github.com/gin-gonic/gin | Web 框架 |
| gorm.io/gorm | ORM |
| gorm.io/driver/mysql | MySQL 驱动 |
| github.com/redis/go-redis/v9 | Redis 客户端 |
| github.com/rabbitmq/amqp091-go | RabbitMQ 客户端 |
| github.com/golang-jwt/jwt/v5 | JWT 认证 |
| github.com/qiniu/go-sdk/v7 | 七牛云 SDK |

## 升级记录

- Go 版本: 1.17 → 1.22+
- JWT: v4 → v5 (使用 RegisteredClaims)
- Redis: go-redis/v8 → redis/go-redis/v9
- RabbitMQ: streadway/amqp → amqp091-go
- 日志: logrus → slog (标准库)
- 数据库: jinzhu/gorm → gorm.io/gorm
- 添加优雅关停支持