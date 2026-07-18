# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 技术栈 (Tech Stack)

- **语言 (Language)**: Go 1.22
- **应用框架 (Framework)**: [Eagle](https://github.com/go-eagle/eagle) 微服务框架 (`go-eagle/eagle`)
- **Web/HTTP**: Gin (`gin-gonic/gin`) + Swagger 文档 (`swaggo/gin-swagger`)
- **RPC / 协议 (RPC / Protocol)**: gRPC (`google.golang.org/grpc`) + Protocol Buffers，带 `protoc-gen-validate` 参数校验
- **依赖注入 (DI)**: Google Wire (`google/wire`)，编译时注入
- **数据库 ORM (Database)**: GORM (`gorm.io/gorm`) + `gorm/gen` 代码生成 + `dbresolver` 读写分离；驱动含 SQLite，另支持 ClickHouse
- **缓存 (Cache)**: Redis (`redis/go-redis/v9`)
- **消息队列 / 任务 (MQ / Task Queue)**: Asynq (`hibiken/asynq`，基于 Redis) + RabbitMQ (`rabbitmq/amqp091-go`)
- **可观测性 (Observability)**: OpenTelemetry (`go.opentelemetry.io/otel`) 链路追踪 + Prometheus (`prometheus/client_golang`) 指标
- **工具库 (Utilities)**: `jinzhu/copier`（对象拷贝）、`pkg/errors`、`golang.org/x/sync`
- **架构 (Architecture)**: Clean Architecture（Service → Repository → DAL）

## 目录结构 (Directory Structure)

```
.
├── api/                # Proto 接口定义与生成代码（gRPC/HTTP/校验）
│   ├── helloworld/     # 示例服务的 proto 定义
│   └── user/           # 用户服务的 proto 定义
├── cmd/                # 程序入口
│   ├── server/         # HTTP/gRPC 主服务入口（含 wire 依赖注入）
│   ├── consumer/       # 后台任务消费者入口
│   └── gen/            # GORM 模型代码生成工具
├── config/             # 配置文件，按环境分目录
│   ├── dev/            # 开发环境配置
│   ├── test/           # 测试环境配置
│   ├── prod/           # 生产环境配置
│   └── docker/         # Docker 环境配置
├── internal/           # 私有业务代码（不对外暴露）
│   ├── service/        # 业务逻辑层（*_svc.go 业务、*_grpc.go 协议转换）
│   ├── repository/     # 仓储层，统一数据访问接口
│   ├── dal/            # 数据访问层（db 数据库 / cache 缓存 / rpc 外部调用）
│   ├── handler/        # HTTP 请求处理器
│   ├── routers/        # 路由注册
│   ├── server/         # 服务器启动装配（HTTP/gRPC）
│   ├── tasks/          # 异步任务定义（Asynq）
│   ├── event/          # 事件处理（消息队列）
│   ├── types/          # 请求/响应等类型定义
│   ├── ecode/          # 业务错误码定义
│   └── mocks/          # 测试用 mock 代码
├── deploy/             # 部署相关
│   ├── docker/         # Dockerfile 等镜像构建文件
│   ├── docker-compose/ # docker-compose 编排文件
│   └── k8s/            # Kubernetes 部署清单
├── third_party/        # 第三方 proto 依赖（google/gogo/validate 等）
├── scripts/            # 构建与运维脚本
├── docs/               # Swagger 生成的接口文档
└── web/                # 前端/静态资源
```

## 编码规范 (Coding Standards)

### 语言与注释 (Language & Comments)
- 与用户交流一律用**中文**，代码注释一律用**英文**
- 命名清晰达意，遵循 Go 官方命名惯例（导出用大驼峰、非导出用小驼峰）

### 设计原则 (Design Principles)
- **不要过度设计**：保证代码简洁易懂、简单实用
- **最小化改动**：改动时尽量不影响其他模块，控制圈复杂度
- **模块化与复用**：注意模块边界，代码尽可能复用，合理使用设计模式
- **分层依赖**：严格遵循 Service → Repository → DAL 的依赖方向，业务逻辑依赖抽象接口而非具体存储
- **不可变性 (Immutability)**：优先返回新对象，避免原地修改，减少隐藏副作用

### 文件组织 (File Organization)
- 多个小文件优于少数大文件，按功能/领域组织而非按类型
- 单文件一般 200-400 行，最多不超过 800 行
- 函数保持短小（建议 <50 行），嵌套层级不超过 4 层，善用 early return

### 错误处理 (Error Handling)
- 每一层都显式处理错误，禁止静默吞掉错误
- 服务端记录详细错误上下文，对外返回友好且不泄露敏感信息的错误
- 使用 `internal/ecode` 统一管理业务错误码

### 输入校验 (Input Validation)
- 在系统边界（API 入口）校验所有外部输入，借助 proto `validate` 规则
- 快速失败并给出清晰的错误信息，不信任任何外部数据

### 安全 (Security)
- 严禁硬编码密钥/密码/Token，统一用环境变量或配置管理
- 数据库使用参数化查询，防止 SQL 注入

### 提交规范 (Commit Convention)
- 遵循 Conventional Commits：`<type>: <description>`
- type 可选：`feat` / `fix` / `refactor` / `docs` / `test` / `chore` / `perf` / `ci`

## 常用开发命令 (Common Development Commands)

### 构建与运行 (Building and Running)
- `make run` - 运行带有 wire 依赖注入的服务器
- `make build` - 构建二进制文件到 `bin/eagle-service`，包含版本信息和竞态检测
- `make wire` - 使用 Google Wire 生成依赖注入代码（生成 `wire_gen.go`）

### 代码生成 (Code Generation)
- `make grpc` - 从 `.proto` 文件生成 gRPC 和 Protocol Buffer 代码
- `make proto` - 生成带有验证的协议缓冲区结构体
- `make gorm-gen` - 使用 `cmd/gen/generate.go` 生成 GORM 模型文件

### 测试与质量 (Testing and Quality)
- `make test` - 运行带有竞态检测的测试
- `make lint` - 运行 golangci-lint 进行代码质量检查
- `make cover` - 生成测试覆盖率报告到 `coverage.txt`
- `make view-cover` - 生成并查看 HTML 覆盖率报告

### 文档 (Documentation)
- `make docs` - 生成 Swagger 文档（可通过 http://localhost:8080/swagger/index.html 访问）

## 架构概述 (Architecture Overview)

这是使用 Clean Architecture 原则构建的 Go 微服务，基于 Eagle 框架：

### 层级结构 (Layer Structure)
```
业务逻辑层 (Service Layer)
    ↓
数据访问抽象层 (Repository Layer)
    ↓
数据访问层 (DAL)
   ├── 数据库操作 (DB)
   ├── Redis 缓存 (Cache)
   └── 外部服务调用 (RPC)
```

### 关键组件 (Key Components)

**依赖注入 (Dependency Injection)**: 使用 Google Wire 进行编译时依赖注入
- `cmd/server/wire.go` - Wire 提供程序定义
- `cmd/server/wire_gen.go` - 生成的依赖注入代码（运行 `make wire` 重新生成）

**服务层 (Service Layer)** (`internal/service/`):
- `*_svc.go` 文件 - 业务逻辑处理
- `*_grpc.go` 文件 - gRPC 协议转换

**仓储层 (Repository Layer)** (`internal/repository/`):
- 提供统一的数据访问接口
- 抽象底层数据存储（数据库、缓存、RPC）

**数据访问层 (DAL Layer)** (`internal/dal/`):
- `db/` - 使用 GORM 的数据库操作
- `cache/` - Redis 缓存操作
- `rpc/` - 外部服务通信

**配置 (Configuration)**: 使用 Eagle 框架的配置系统
- `config/` 目录中的配置文件，按环境组织（`dev/`、`prod/`、`test/`）
- 通过 main.go 中的 `config.New()` 加载

### 协议缓冲区和 gRPC (Protocol Buffers and gRPC)
- `api/` 目录中的 Proto 定义
- 生成的文件包含验证、gRPC 和 HTTP 网关绑定
- 修改 `.proto` 文件后使用 `make grpc`

### 入口点 (Entry Points)
- `cmd/server/main.go` - 主要 HTTP/gRPC 服务器
- `cmd/consumer/main.go` - 后台作业消费者
- 两者都使用 Wire 进行依赖注入

### 开发工作流 (Development Workflow)
1. 根据需要修改 proto 文件 → `make grpc`
2. 运行服务器 → `make run`
3. 在 service/repository 层添加业务逻辑
4. 生成 wire 代码 → `make wire`（如果依赖关系发生变化）
5. 测试 → `make test`
6. 构建 → `make build`

## Never 规则 (Never Rules)

以下为**绝对禁止**事项，任何情况下都不得违反：

### 交流与注释 (Communication & Comments)
- **Never** 用中文以外的语言回复用户（代码注释除外，注释一律用英文）
- **Never** 用中文写代码注释

### 代码设计 (Code Design)
- **Never** 过度设计；保证代码简洁易懂、简单实用
- **Never** 忽视圈复杂度；重复代码应尽量复用
- **Never** 在改动时波及无关模块，坚持最小化修改
- **Never** 原地修改传入对象；优先返回新对象（保持不可变）
- **Never** 让单文件超过 800 行、函数超过 50 行、嵌套超过 4 层

### 分层与架构 (Layering & Architecture)
- **Never** 跨层反向依赖；严格遵循 Service → Repository → DAL 的依赖方向
- **Never** 在 Service 层直接操作数据库/缓存，必须经由 Repository 抽象
- **Never** 手改由 `make wire` / `make grpc` / `make proto` / `make gorm-gen` 生成的文件（如 `wire_gen.go`、`*.pb.go`）

### 错误与校验 (Errors & Validation)
- **Never** 静默吞掉错误；每一层都必须显式处理
- **Never** 信任外部输入；必须在系统边界完成校验
- **Never** 在对外错误信息中泄露敏感数据

### 安全 (Security)
- **Never** 硬编码密钥/密码/Token，统一用环境变量或配置管理
- **Never** 使用字符串拼接 SQL；必须参数化查询
- **Never** 提交调试语句、临时打印或注释掉的死代码

### 提交 (Commit)
- **Never** 使用不符合 Conventional Commits 规范的提交信息
- **Never** 在未通过 `make lint` 和 `make test` 前提交代码

## Commit 规范

- 格式：type(scope): description
- 类型：feat / fix / docs / style / refactor / test / chore

## 多 Agent 并发

- 禁止 git stash
- 禁止切换分支
- 只 commit 自己修改的文件
