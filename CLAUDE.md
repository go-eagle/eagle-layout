# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

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

## 注意事项（system prompt）

- Always respond in 中文，但注释一律用英文
- 不要过度设计，保证代码简洁易懂，简单实用
- 写代码时，要注意圈复杂度，代码尽可能复用
- 写代码时，注意模块设计，尽量使用设计模式
- 改动时最小化修改，尽量不修改到其他模块代码