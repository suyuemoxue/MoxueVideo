# 墨雪影视后端（MoxueVideo）

一个使用 Go 构建的多服务后端示例工程，当前包含两个服务：

- **core**：HTTP API 服务（Gin + GORM + Redis + RabbitMQ + gRPC 客户端 + 可选 OSS STS）
- **chat**：聊天服务（gRPC + WebSocket + GORM + RabbitMQ）

仓库采用多 module 组织：每个服务目录下都有自己的 `go.mod`，根目录使用 `go.work` 统一管理本地开发。

## 目录结构与作用

```
MoxueVideo/
  infra/                         # 本地依赖：MySQL/Redis/RabbitMQ + 初始化 SQL
    docker-compose.yaml
    schema.sql                   # core 使用的数据库/表结构
    chat_schema.sql              # chat 使用的数据库/表结构
  services/
    core/                        # 核心服务（HTTP API）
      cmd/api/main.go            # 入口：启动 HTTP、连接依赖、消费 MQ 事件
      internal/
        config/                  # 配置加载（环境变量 + config.local.yaml）
        domain/                  # 领域实体/事件定义
        infra/                   # MySQL/Redis/MQ/OSS/gRPC 等基础设施适配
        middleware/              # HTTP 中间件（如鉴权）
        transport/httpapi/       # Gin 路由、handler、WS 通知
        usecase/                 # 用例层（业务编排）
    chat/                        # 聊天服务（gRPC + WS）
      cmd/grpc/main.go           # 入口：启动 gRPC 服务 + WS 服务
      internal/
        config/                  # 配置加载（环境变量）
        domain/                  # 领域事件/对象
        infra/                   # MySQL/MQ 等基础设施适配
        transport/               # gRPC/WS 传输层
        usecase/                 # 用例层（发消息、建会话等）
  go.work                        # 本地多 module 工作区
  README.md
```

## 本地依赖（Docker Compose）

仓库提供 `infra/docker-compose.yaml` 用于一键启动本地依赖：

- **MySQL**：`127.0.0.1:3307`
- **Redis**：`127.0.0.1:6379`（带密码）
- **RabbitMQ**：`127.0.0.1:5672`
- **RabbitMQ 管理台**：`http://127.0.0.1:15672`

初始化 SQL：

- `infra/schema.sql`：创建 `moxuevideo` 并初始化 core 相关表
- `infra/chat_schema.sql`：创建 `moxuevideo_chat` 并初始化 chat 相关表

## 启动方法（推荐：WSL）

### 1) 启动依赖（MySQL/Redis/RabbitMQ）

在仓库根目录执行：

```bash
cd infra
docker compose up -d
```

### 2) 启动 chat 服务

```bash
cd services/chat
go run ./cmd/grpc
```

默认监听：

- gRPC：`:50051`
- WebSocket：`:50052`
  - WS 路径：`/ws/chat`
  - WS 健康检查：`http://127.0.0.1:50052/healthz`

### 3) 启动 core 服务

```bash
cd services/core
go run ./cmd/api
```

默认监听：

- HTTP：`:8080`
- 健康检查：`http://127.0.0.1:8080/healthz`

## 常用配置（环境变量）

### core

core 支持环境变量覆盖默认配置，也会尝试读取 `services/core/config.local.yaml`（该文件被 `.gitignore` 忽略，适合放本地开发配置）。

- `HTTP_ADDR`：默认 `:8080`
- `MYSQL_DSN`：默认空（建议通过环境变量或 `config.local.yaml` 配置）
- `REDIS_ADDR`：默认 `127.0.0.1:6379`
- `REDIS_PASSWORD`：默认空（如需密码请自行配置）
- `REDIS_DB`：默认 `0`
- `RABBITMQ_URL`：默认空（如需连接请自行配置）
- `CHAT_GRPC_ADDR`：默认 `127.0.0.1:50051`

数据库迁移：

- `MYSQL_AUTOMIGRATE`：默认不执行；设为 `1 / true / yes` 才会触发 GORM AutoMigrate

OSS STS（可选）：当以下四项都配置时才会初始化，否则接口会返回 `OSS_STS_UNAVAILABLE`。

- `OSS_ACCESS_KEY_ID`
- `OSS_ACCESS_KEY_SECRET`
- `OSS_ROLE_ARN`
- `OSS_BUCKET`

### chat

- `GRPC_ADDR`：默认 `:50051`
- `WS_ADDR`：默认 `:50052`
- `MYSQL_DSN`：默认空（请自行配置）
- `RABBITMQ_URL`：默认空（如需连接请自行配置）

## 验证服务是否正常

core：

```bash
curl -sS http://127.0.0.1:8080/healthz
```

chat：

```bash
curl -sS -o /dev/null -w "%{http_code}\n" http://127.0.0.1:50052/healthz
```

## 备注

- 当前部分业务接口 handler 仍可能返回 `NOT_IMPLEMENTED`，但服务启动、依赖连接与健康检查已可用。
