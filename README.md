# 墨雪影视后端（MoxueVideo）

本仓库是一个使用 Go + Gin + GORM 的后端项目，当前采用多个独立 Go Module 的服务目录组织方式：每个服务目录下都有自己的 `go.mod`，可单独启动与测试。

## 目录结构总览

```
.
 README.md
 services/
    infra/                       # 基础依赖（MySQL/Redis/RabbitMQ）容器编排
       docker-compose.yaml
    user-service/                # 用户服务（注册/登录/查询当前用户）
    video-service/               # 视频服务（投稿/查询视频/Feed）
    relation-service/            # 关系服务（关注/粉丝/关注列表）
    interaction-service/         # 互动服务（点赞/收藏/列表）
 .gitignore
```

各服务目录内部结构基本一致（以 `video-service` 为例）：

```
services/video-service/
 cmd/server/main.go               # 服务启动入口：加载配置、初始化依赖、启动 HTTP Server
 internal/
    config/                      # 配置结构与环境变量读取
    handler/                     # Gin Handler（HTTP 层：参数校验、调用 service、返回响应）
    infra/mysql/                 # MySQL 连接初始化（GORM）
    logger/                      # slog logger 初始化
    middleware/                  # Gin 中间件（JWT 鉴权）
    model/                       # GORM Model（数据库表结构）
    repo/                        # Repo 层（数据库访问封装）
    service/                     # 业务层（领域逻辑）
    transport/
        httpserver/              # 路由注册（Gin Router）
        httpx/                   # 统一响应封装、Query 解析、Context helper
├ go.mod
 go.sum
```

## 基础设施（services/infra）

文件：`services/infra/docker-compose.yaml`

- MySQL: `3306`
  - DB: `shortvideo`
  - user/pass: `app / apppass`
- Redis: `6379`
- RabbitMQ: `5672`（管理台 `15672`）

启动：

```bash
cd services/infra
docker compose up -d
```

## 通用约定（所有服务）

### 统一响应结构（httpx.Response）

所有接口返回统一 envelope：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

- `code == 0` 表示成功
- 失败时 `code` 通常为：
  - `40000` 参数错误
  - `40100` 未登录/Token 无效
  - `40400` 资源不存在
  - `40900` 冲突（如用户名已占用）
  - `50000` 服务端错误

实现位置（各服务一致）：`services/*/internal/transport/httpx/response.go`

### JWT 鉴权

- 需要登录的接口通过 `Authorization` Header 传递 Bearer Token：
  - `Authorization: Bearer <token>`
- 服务器使用 `JWT_SECRET`（HS256）验签，通过后把 `userID` 写入 Gin Context（key 为 `userID`）

实现位置：

- `services/*/internal/middleware/jwt_auth.go`
- `services/*/internal/transport/httpx/context.go`

### 配置（环境变量）

所有服务 `internal/config/config.go` 使用相同环境变量名和默认值（不同服务可能未使用某些依赖，但配置项保留）。

常用项（默认值）：

- `APP_ENV`：`dev`
- `HTTP_ADDR`：`:8080`
- `MYSQL_HOST`：`127.0.0.1`
- `MYSQL_PORT`：`3306`
- `MYSQL_USER`：`app`
- `MYSQL_PASSWORD`：`apppass`
- `MYSQL_DB`：`shortvideo`
- `MYSQL_PARAMS`：`charset=utf8mb4&parseTime=True&loc=Local`
- `REDIS_ADDR`：`127.0.0.1:6379`
- `REDIS_PASSWORD`：空
- `REDIS_DB`：`0`
- `RABBITMQ_ADDR`：`amqp://app:apppass@127.0.0.1:5672/`
- `JWT_SECRET`：默认空（各服务 `main.go` 会在为空时使用 `dev-secret`）

## 如何启动服务

注意：所有服务默认 `HTTP_ADDR=:8080`，如果同时启动多个服务，请分别设置不同端口（例如 `:8081/:8082/...`）。

### 启动拆分服务（user/video/relation/interaction）

```bash
cd services/user-service && go run ./cmd/server
cd services/video-service && go run ./cmd/server
cd services/relation-service && go run ./cmd/server
cd services/interaction-service && go run ./cmd/server
```

## 服务与接口说明

下面按服务列出路由、参数与鉴权要求。所有服务均提供健康检查接口：

- `GET /healthz`：返回 `data={"status":"ok"}`

### user-service（用户服务）

代码入口：`services/user-service/cmd/server/main.go`

路由定义：`services/user-service/internal/transport/httpserver/router.go`

接口：

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/me`（需要登录）

说明：

- 该服务只管理 `User` 表（`AutoMigrate(&model.User{})`）
- token 的签发/解析逻辑在 `services/user-service/internal/service/auth_service.go` 内

### video-service（视频服务）

代码入口：`services/video-service/cmd/server/main.go`

路由定义：`services/video-service/internal/transport/httpserver/router.go`

接口：

- `GET /api/v1/videos/feed`
  - Query：`cursor`、`limit`
  - 返回：`data.items: Video[]`、`data.nextCursor`
- `GET /api/v1/videos/:id`
  - 返回：`Video`
- `GET /api/v1/users/:id/videos`
  - Query：`page`、`size`
  - 返回：`data.items: Video[]`
- `POST /api/v1/videos/publish`（需要登录）
  - Body：`playUrl/coverUrl/title/description`
  - 返回：`Video`

说明：

- video-service 返回的是基础 `model.Video`（非富 DTO）
- 该服务只管理 `Video` 表（`AutoMigrate(&model.Video{})`）

### relation-service（关系服务）

代码入口：`services/relation-service/cmd/server/main.go`

路由定义：`services/relation-service/internal/transport/httpserver/router.go`

接口（全部需要登录）：

- `POST /api/v1/follow/action`
  - Body：`{ "toUserId": 2, "action": 1 }`
- `GET /api/v1/follow/following`
  - Query：`userId`（默认当前用户）、`page`、`size`
- `GET /api/v1/follow/followers`
  - Query：`userId`（默认当前用户）、`page`、`size`

返回：

- action 返回：`data.ok=true`
- 列表返回：`data.items: UserDTO[]`

说明：

- 该服务维护 `User` 与 `Follow` 表（`AutoMigrate(&model.User{}, &model.Follow{})`）
- `FollowService.SetFollow` 会校验不能关注自己、被关注用户必须存在

### interaction-service（互动服务）

代码入口：`services/interaction-service/cmd/server/main.go`

路由定义：`services/interaction-service/internal/transport/httpserver/router.go`

接口（全部需要登录）：

- `POST /api/v1/likes/action`
  - Body：`{ "videoId": 1, "action": 1 }`
  - 返回：`data.status="ok"`
- `GET /api/v1/likes/list`
  - Query：`userId`（默认当前用户）、`page`、`size`
  - 返回：`data.items: uint64[]`（视频 ID 列表）
- `POST /api/v1/favorites/action`
  - Body：`{ "videoId": 1, "action": 1 }`
  - 返回：`data.status="ok"`
- `GET /api/v1/favorites/list`
  - Query：`userId`（默认当前用户）、`page`、`size`
  - 返回：`data.items: uint64[]`

说明：

- 该服务维护 `Video`、`Like`、`Favorite` 表（`AutoMigrate(&model.Video{}, &model.Like{}, &model.Favorite{})`）
- list 接口只返回 id 列表；如果需要视频详情/富信息，建议由上层（如网关/聚合层）再去查询视频服务并组装

## 数据模型（MySQL）

各服务按需维护各自的数据表，迁移逻辑在各服务的 `cmd/server/main.go` 中：

- `user-service`: `User`
- `video-service`: `Video`
- `relation-service`: `User`、`Follow`
- `interaction-service`: `Video`、`Like`、`Favorite`

## 开发与测试

各服务独立执行（示例）：

```bash
cd services/video-service
go test ./...
go vet ./...
```
