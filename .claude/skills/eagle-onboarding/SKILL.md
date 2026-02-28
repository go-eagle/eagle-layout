---
name: eagle-onboarding
description: Eagle 框架 Go 微服务项目新手入门指导。适用于刚加入团队的开发者需要了解项目架构、开发命令、工作流程时使用。当用户询问"如何开始"、"怎么添加新功能"、"项目结构是什么"、"开发流程"、"如何添加 API"、"修改数据模型"、"Wire 怎么用"、"gRPC 如何配置"等相关问题时触发。确保在新人提出项目相关问题、需要架构说明或开发指导时主动使用此 skill。
---

# Eagle 框架新手入门指南

欢迎加入 Eagle 微服务项目！这份指南将帮助你快速了解项目结构、掌握开发流程，并学会如何添加新功能。

## 🏗️ 项目架构概览

本项目基于 **Clean Architecture（整洁架构）** 设计，使用 Eagle 框架构建。核心分为三层：

```
┌─────────────────────────────────────────┐
│         Service Layer (业务逻辑层)        │
│   • 处理业务逻辑                          │
│   • 协议转换 (gRPC/HTTP)                  │
│   • 输入输出验证                          │
└─────────────┬───────────────────────────┘
              │
┌─────────────▼───────────────────────────┐
│      Repository Layer (仓储抽象层)        │
│   • 统一数据访问接口                      │
│   • 缓存策略管理                          │
│   • 数据聚合逻辑                          │
└─────────────┬───────────────────────────┘
              │
┌─────────────▼───────────────────────────┐
│         DAL Layer (数据访问层)            │
│   ├─ DB: 数据库操作 (GORM)                │
│   ├─ Cache: Redis 缓存                    │
│   └─ RPC: 外部服务调用                    │
└─────────────────────────────────────────┘
```

### 目录结构

```
eagle-layout/
├── api/                    # Proto 文件定义
│   ├── user/v1/           # 用户服务 API 定义
│   └── helloworld/        # 示例服务
├── cmd/                    # 应用入口
│   ├── server/            # HTTP/gRPC 服务器
│   │   ├── main.go        # 主程序入口
│   │   ├── wire.go        # Wire 依赖注入定义
│   │   └── wire_gen.go    # Wire 生成的代码
│   └── consumer/          # 后台任务消费者
├── config/                 # 配置文件
│   ├── dev/               # 开发环境配置
│   ├── test/              # 测试环境配置
│   └── prod/              # 生产环境配置
├── internal/               # 内部代码（不对外暴露）
│   ├── service/           # Service 层：业务逻辑
│   │   ├── user_svc.go    # 业务逻辑实现
│   │   └── user_grpc.go   # gRPC 协议转换
│   ├── repository/        # Repository 层：数据访问抽象
│   │   └── user_repo.go   # 仓储接口和实现
│   ├── dal/               # DAL 层：底层数据访问
│   │   ├── db/            # 数据库相关
│   │   │   ├── dao/       # GORM Gen 生成的 DAO
│   │   │   └── model/     # 数据模型
│   │   ├── cache/         # Redis 缓存
│   │   └── rpc/           # 外部服务调用
│   ├── server/            # 服务器配置
│   ├── types/             # 内部类型定义
│   └── ecode/             # 错误码定义
└── Makefile               # 开发命令集合
```

## 🔧 常用开发命令

所有命令都定义在 `Makefile` 中，使用 `make <command>` 执行：

### 构建与运行

```bash
# 运行服务器（包含 Wire 依赖注入）
make run

# 构建二进制文件到 bin/eagle-service
# 包含版本信息和竞态检测
make build
```

### 代码生成

```bash
# 生成 Wire 依赖注入代码
# 会生成 cmd/server/wire_gen.go
make wire

# 从 .proto 文件生成 gRPC 和 Protocol Buffer 代码
# 生成文件到 api/ 目录下对应位置
make grpc

# 生成 Protocol Buffer 结构体（带验证）
make proto

# 生成 GORM 模型文件
# 使用 cmd/gen/generate.go 配置
make gorm-gen
```

### 测试与质量

```bash
# 运行测试（带竞态检测）
make test

# 代码检查（golangci-lint）
make lint

# 生成测试覆盖率报告
make cover

# 生成并查看 HTML 覆盖率报告
make view-cover
```

### 文档

```bash
# 生成 Swagger 文档
make docs

# 访问 Swagger UI
# http://localhost:8080/swagger/index.html
```

## 🚀 核心技术栈使用说明

### 1. Google Wire 依赖注入

Wire 是编译时依赖注入工具，避免运行时反射带来的性能损耗。

#### Wire 工作原理

**定义 Provider（提供者）**：

`cmd/server/wire.go`:
```go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
    "github.com/go-eagle/eagle-layout/internal/server"
)

// InitApp 定义依赖关系
func InitApp(cfg *eagle.Config) (*eagle.App, func(), error) {
    wire.Build(
        server.ServerSet,  // Server 层依赖集
        newApp,            // App 构造函数
    )
    return &eagle.App{}, nil, nil
}
```

**Provider Set 示例**：

`internal/server/grpc.go`:
```go
// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
    NewHTTPServer,
    NewGRPCServer,
    service.ProviderSet,
    repository.ProviderSet,
    dal.ProviderSet,
)
```

**生成代码**：
```bash
make wire
# 生成 cmd/server/wire_gen.go
```

#### 添加新依赖的步骤

1. 创建构造函数（返回接口类型）：
```go
// internal/service/order_svc.go
func NewOrderService(repo repository.OrderRepo) *OrderService {
    return &OrderService{repo: repo}
}
```

2. 添加到 ProviderSet：
```go
// internal/service/service.go
var ProviderSet = wire.NewSet(
    NewUserService,
    NewOrderService,  // 新增
)
```

3. 重新生成 Wire 代码：
```bash
make wire
```

### 2. gRPC 和 Protocol Buffers

#### Proto 文件定义

`api/user/v1/user.proto`:
```protobuf
syntax = "proto3";

package user.v1;

import "validate/validate.proto";
import "google/api/annotations.proto";

option go_package = "github.com/go-eagle/eagle-layout/api/user/v1;v1";

// 用户服务
service UserService {
  // 创建用户
  rpc CreateUser(CreateUserRequest) returns(CreateUserReply) {
    option (google.api.http) = {
        post: "/v1/users/"
        body: "*"
    };
  }

  // 获取用户
  rpc GetUser(GetUserRequest) returns (GetUserReply) {
    option (google.api.http) = {
        get: "/v1/users/{id}"
    };
  }
}

// 请求消息（带验证）
message CreateUserRequest {
  string username = 1 [(validate.rules).string.min_len = 6];
  string email = 2 [(validate.rules).string.email = true];
  string password = 3 [(validate.rules).string.min_len = 6];
}

// 响应消息
message CreateUserReply {
  int64 id = 1;
  string username = 2;
  string email = 3;
}
```

#### 生成代码

```bash
# 从 proto 文件生成 Go 代码
make grpc

# 生成的文件：
# api/user/v1/user.pb.go           - Protocol Buffer 定义
# api/user/v1/user_grpc.pb.go      - gRPC 服务代码
# api/user/v1/user.pb.validate.go  - 验证代码
# api/user/v1/user_http.pb.go      - HTTP 网关代码
```

#### 实现 gRPC 服务

`internal/service/user_grpc.go`:
```go
package service

import (
    "context"
    pb "github.com/go-eagle/eagle-layout/api/user/v1"
    "github.com/go-eagle/eagle-layout/internal/types"
)

// 确保实现了 gRPC 接口
var _ pb.UserServiceServer = (*UserService)(nil)

// CreateUser implements gRPC CreateUser method
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
    // 1. 协议转换：gRPC -> 内部类型
    input := types.CreateUserInput{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }

    // 2. 调用业务逻辑
    output, err := s.CreateUser(ctx, input)
    if err != nil {
        return nil, err
    }

    // 3. 协议转换：内部类型 -> gRPC
    return &pb.CreateUserReply{
        Id:       output.ID,
        Username: output.Username,
        Email:    output.Email,
    }, nil
}
```

### 3. GORM 数据库操作

本项目使用 **GORM Gen** 生成类型安全的数据库操作代码。

#### 数据模型定义

`internal/dal/db/model/user.go`:
```go
package model

// UserInfoModel 用户信息表
type UserInfoModel struct {
    ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
    Username  string `gorm:"column:username;type:varchar(50);uniqueIndex;not null"`
    Email     string `gorm:"column:email;type:varchar(100);uniqueIndex;not null"`
    Password  string `gorm:"column:password;type:varchar(255);not null"`
    Nickname  string `gorm:"column:nickname;type:varchar(50)"`
    Avatar    string `gorm:"column:avatar;type:varchar(255)"`
    Status    int32  `gorm:"column:status;default:0"`
    CreatedAt int64  `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserInfoModel) TableName() string {
    return "user_info"
}
```

#### 生成 DAO

`cmd/gen/generate.go`:
```go
package main

import (
    "gorm.io/gen"
    "github.com/go-eagle/eagle-layout/internal/dal/db/model"
)

func main() {
    g := gen.NewGenerator(gen.Config{
        OutPath: "./internal/dal/db/dao",
        Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
    })

    // 使用数据库连接
    g.UseDB(db)

    // 生成模型对应的 DAO
    g.ApplyBasic(model.UserInfoModel{})

    // 自定义查询方法
    g.ApplyInterface(func(method gen.Method) {}, model.UserInfoModel{})

    g.Execute()
}
```

运行生成：
```bash
make gorm-gen
# 生成到 internal/dal/db/dao/
```

#### 使用 DAO 查询

Repository 层使用生成的 DAO：

`internal/repository/user_repo.go`:
```go
package repository

import (
    "context"
    "github.com/go-eagle/eagle-layout/internal/dal/db/dao"
    "github.com/go-eagle/eagle-layout/internal/dal/db/model"
)

type userRepo struct {
    db    *dal.DBClient
    cache cache.UserCache
}

// GetUser 获取单个用户
func (r *userRepo) GetUser(ctx context.Context, id int64) (*model.UserInfoModel, error) {
    // 使用生成的 DAO 查询
    user, err := dao.UserInfoModel.WithContext(ctx).
        Where(dao.UserInfoModel.ID.Eq(id)).
        First()
    if err != nil {
        return nil, err
    }
    return user, nil
}

// CreateUser 创建用户
func (r *userRepo) CreateUser(ctx context.Context, data model.UserInfoModel) (int64, error) {
    err := dao.UserInfoModel.WithContext(ctx).Create(&data)
    if err != nil {
        return 0, err
    }
    return data.ID, nil
}

// UpdateUser 更新用户
func (r *userRepo) UpdateUser(ctx context.Context, id int64, data model.UserInfoModel) error {
    _, err := dao.UserInfoModel.WithContext(ctx).
        Where(dao.UserInfoModel.ID.Eq(id)).
        Updates(data)
    return err
}
```

#### 缓存集成

Repository 层集成多级缓存（本地缓存 + Redis）：

`internal/repository/user_repo.go`:
```go
func (r *userRepo) GetUser(ctx context.Context, id int64) (*model.UserInfoModel, error) {
    // 1. 尝试本地缓存
    var ret *model.UserInfoModel
    err := r.localCache.Get(ctx, cast.ToString(id), &ret)
    if err == nil && ret != nil && ret.ID > 0 {
        return ret, nil
    }

    // 2. 尝试 Redis 缓存
    ret, err = r.cache.GetUserCache(ctx, id)
    if err == nil && ret != nil && ret.ID > 0 {
        return ret, nil
    }

    // 3. 查询数据库（使用 singleflight 防止缓存击穿）
    val, err, _ := r.sg.Do("sg:user:"+cast.ToString(id), func() (interface{}, error) {
        data, err := dao.UserInfoModel.WithContext(ctx).
            Where(dao.UserInfoModel.ID.Eq(id)).
            First()
        if err != nil {
            return nil, err
        }

        // 4. 写入缓存
        _ = r.cache.SetUserCache(ctx, id, data, 5*time.Minute)
        _ = r.localCache.Set(ctx, cast.ToString(id), data, 2*time.Minute)

        return data, nil
    })

    if err != nil {
        return nil, err
    }

    return val.(*model.UserInfoModel), nil
}
```

## 📝 完整开发工作流程

### 场景 1: 添加新的 API 接口

假设你要添加一个"获取用户列表"的接口。

#### 步骤 1: 定义 Proto

编辑 `api/user/v1/user.proto`：

```protobuf
service UserService {
  // 新增：获取用户列表
  rpc ListUsers(ListUsersRequest) returns (ListUsersReply) {
    option (google.api.http) = {
        get: "/v1/users"
    };
  }
}

message ListUsersRequest {
  int32 page = 1 [(validate.rules).int32.gte = 1];
  int32 page_size = 2 [(validate.rules).int32 = {gte: 1, lte: 100}];
}

message ListUsersReply {
  repeated User users = 1;
  int32 total = 2;
}
```

#### 步骤 2: 生成 gRPC 代码

```bash
make grpc
```

#### 步骤 3: 定义内部类型

创建 `internal/types/user.go`:

```go
package types

type ListUsersInput struct {
    Page     int32
    PageSize int32
}

type ListUsersOutput struct {
    Users []*User
    Total int32
}
```

#### 步骤 4: 实现 Repository 层

在 `internal/repository/user_repo.go` 添加接口方法：

```go
type UserRepo interface {
    // ... 现有方法
    ListUsers(ctx context.Context, page, pageSize int32) ([]*model.UserInfoModel, int64, error)
}

func (r *userRepo) ListUsers(ctx context.Context, page, pageSize int32) ([]*model.UserInfoModel, int64, error) {
    offset := (page - 1) * pageSize

    // 查询列表
    users, err := dao.UserInfoModel.WithContext(ctx).
        Limit(int(pageSize)).
        Offset(int(offset)).
        Order(dao.UserInfoModel.CreatedAt.Desc()).
        Find()
    if err != nil {
        return nil, 0, err
    }

    // 查询总数
    total, err := dao.UserInfoModel.WithContext(ctx).Count()
    if err != nil {
        return nil, 0, err
    }

    return users, total, nil
}
```

#### 步骤 5: 实现 Service 层业务逻辑

在 `internal/service/user_svc.go` 添加：

```go
func (s *UserService) ListUsers(ctx context.Context, input types.ListUsersInput) (*types.ListUsersOutput, error) {
    users, total, err := s.repo.ListUsers(ctx, input.Page, input.PageSize)
    if err != nil {
        return nil, fmt.Errorf("[UserService] ListUsers error: %w", err)
    }

    // 转换为内部类型
    var userList []*types.User
    for _, u := range users {
        user, err := s.convertUser(u)
        if err != nil {
            continue
        }
        userList = append(userList, user)
    }

    return &types.ListUsersOutput{
        Users: userList,
        Total: int32(total),
    }, nil
}
```

#### 步骤 6: 实现 gRPC 协议转换

在 `internal/service/user_grpc.go` 添加：

```go
func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
    // gRPC -> 内部类型
    input := types.ListUsersInput{
        Page:     req.Page,
        PageSize: req.PageSize,
    }

    // 调用业务逻辑
    output, err := s.ListUsers(ctx, input)
    if err != nil {
        return nil, err
    }

    // 内部类型 -> gRPC
    var pbUsers []*pb.User
    for _, u := range output.Users {
        pbUsers = append(pbUsers, &pb.User{
            Id:       u.Id,
            Username: u.Username,
            Email:    u.Email,
            // ... 其他字段
        })
    }

    return &pb.ListUsersReply{
        Users: pbUsers,
        Total: output.Total,
    }, nil
}
```

#### 步骤 7: 测试

```bash
# 运行服务
make run

# 测试 HTTP 接口
curl "http://localhost:8080/v1/users?page=1&page_size=10"

# 运行单元测试
make test
```

### 场景 2: 修改数据模型

假设要给用户表添加"最后登录时间"字段。

#### 步骤 1: 修改数据库表结构

```sql
ALTER TABLE user_info
ADD COLUMN last_login_at BIGINT DEFAULT 0 COMMENT '最后登录时间';
```

#### 步骤 2: 更新 Model 定义

编辑 `internal/dal/db/model/user.go`:

```go
type UserInfoModel struct {
    // ... 现有字段
    LastLoginAt int64 `gorm:"column:last_login_at;default:0"`
}
```

#### 步骤 3: 重新生成 DAO

```bash
make gorm-gen
```

#### 步骤 4: 更新 Proto 定义

编辑 `api/user/v1/user.proto`:

```protobuf
message User {
  // ... 现有字段
  int64 last_login_at = 14;
}
```

```bash
make grpc
```

#### 步骤 5: 更新相关业务逻辑

在需要的地方更新字段使用：

```go
// internal/service/user_svc.go
func (s *UserService) Login(ctx context.Context, input types.LoginInput) (*types.LoginOutput, error) {
    // ... 登录逻辑

    // 更新最后登录时间
    err = s.repo.UpdateUser(ctx, user.ID, model.UserInfoModel{
        LastLoginAt: time.Now().Unix(),
    })

    // ...
}
```

### 场景 3: 添加新的业务服务

假设要添加订单服务（Order Service）。

#### 步骤 1: 创建 Proto 定义

创建 `api/order/v1/order.proto`:

```protobuf
syntax = "proto3";

package order.v1;

import "validate/validate.proto";
import "google/api/annotations.proto";

option go_package = "github.com/go-eagle/eagle-layout/api/order/v1;v1";

service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderReply) {
    option (google.api.http) = {
        post: "/v1/orders"
        body: "*"
    };
  }

  rpc GetOrder(GetOrderRequest) returns (GetOrderReply) {
    option (google.api.http) = {
        get: "/v1/orders/{id}"
    };
  }
}

message CreateOrderRequest {
  int64 user_id = 1 [(validate.rules).int64.gte = 1];
  repeated int64 product_ids = 2;
  string address = 3;
}

message CreateOrderReply {
  int64 id = 1;
  string order_no = 2;
}

message GetOrderRequest {
  int64 id = 1;
}

message GetOrderReply {
  Order order = 1;
}

message Order {
  int64 id = 1;
  string order_no = 2;
  int64 user_id = 3;
  int64 total_amount = 4;
  int32 status = 5;
}
```

生成代码：
```bash
make grpc
```

#### 步骤 2: 创建数据模型

创建 `internal/dal/db/model/order.go`:

```go
package model

type OrderModel struct {
    ID          int64  `gorm:"column:id;primaryKey;autoIncrement"`
    OrderNo     string `gorm:"column:order_no;type:varchar(50);uniqueIndex;not null"`
    UserID      int64  `gorm:"column:user_id;index;not null"`
    TotalAmount int64  `gorm:"column:total_amount;not null"`
    Status      int32  `gorm:"column:status;default:0"`
    CreatedAt   int64  `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   int64  `gorm:"column:updated_at;autoUpdateTime"`
}

func (OrderModel) TableName() string {
    return "orders"
}
```

生成 DAO：
```bash
# 更新 cmd/gen/generate.go，添加 OrderModel
# 然后运行
make gorm-gen
```

#### 步骤 3: 创建 Repository 层

创建 `internal/repository/order_repo.go`:

```go
package repository

import (
    "context"
    "github.com/go-eagle/eagle-layout/internal/dal"
    "github.com/go-eagle/eagle-layout/internal/dal/db/dao"
    "github.com/go-eagle/eagle-layout/internal/dal/db/model"
)

type OrderRepo interface {
    CreateOrder(ctx context.Context, data model.OrderModel) (int64, error)
    GetOrder(ctx context.Context, id int64) (*model.OrderModel, error)
}

type orderRepo struct {
    db *dal.DBClient
}

func NewOrderRepo(db *dal.DBClient) OrderRepo {
    return &orderRepo{db: db}
}

func (r *orderRepo) CreateOrder(ctx context.Context, data model.OrderModel) (int64, error) {
    err := dao.OrderModel.WithContext(ctx).Create(&data)
    if err != nil {
        return 0, err
    }
    return data.ID, nil
}

func (r *orderRepo) GetOrder(ctx context.Context, id int64) (*model.OrderModel, error) {
    return dao.OrderModel.WithContext(ctx).
        Where(dao.OrderModel.ID.Eq(id)).
        First()
}
```

更新 `internal/repository/repository.go` 的 ProviderSet：

```go
var ProviderSet = wire.NewSet(
    NewUserRepo,
    NewOrderRepo,  // 新增
)
```

#### 步骤 4: 创建 Service 层

创建 `internal/service/order_svc.go`:

```go
package service

import (
    "context"
    "fmt"
    "time"
    "github.com/go-eagle/eagle-layout/internal/repository"
    "github.com/go-eagle/eagle-layout/internal/dal/db/model"
    "github.com/go-eagle/eagle-layout/internal/types"
)

type OrderService struct {
    repo repository.OrderRepo
}

func NewOrderService(repo repository.OrderRepo) *OrderService {
    return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, input types.CreateOrderInput) (*types.CreateOrderOutput, error) {
    // 生成订单号
    orderNo := fmt.Sprintf("ORD%d", time.Now().UnixNano())

    // 创建订单
    order := model.OrderModel{
        OrderNo:     orderNo,
        UserID:      input.UserID,
        TotalAmount: input.TotalAmount,
        Status:      0,
        CreatedAt:   time.Now().Unix(),
    }

    id, err := s.repo.CreateOrder(ctx, order)
    if err != nil {
        return nil, fmt.Errorf("[OrderService] CreateOrder error: %w", err)
    }

    return &types.CreateOrderOutput{
        ID:      id,
        OrderNo: orderNo,
    }, nil
}
```

创建 `internal/service/order_grpc.go`:

```go
package service

import (
    "context"
    pb "github.com/go-eagle/eagle-layout/api/order/v1"
    "github.com/go-eagle/eagle-layout/internal/types"
)

var _ pb.OrderServiceServer = (*OrderService)(nil)

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
    input := types.CreateOrderInput{
        UserID:      req.UserId,
        ProductIDs:  req.ProductIds,
        TotalAmount: 0, // 实际应计算
    }

    output, err := s.CreateOrder(ctx, input)
    if err != nil {
        return nil, err
    }

    return &pb.CreateOrderReply{
        Id:      output.ID,
        OrderNo: output.OrderNo,
    }, nil
}
```

更新 `internal/service/service.go` 的 ProviderSet：

```go
var ProviderSet = wire.NewSet(
    NewUserService,
    NewOrderService,  // 新增
)
```

#### 步骤 5: 注册 gRPC 服务

编辑 `internal/server/grpc.go`:

```go
package server

import (
    userV1 "github.com/go-eagle/eagle-layout/api/user/v1"
    orderV1 "github.com/go-eagle/eagle-layout/api/order/v1"  // 新增
    "github.com/go-eagle/eagle-layout/internal/service"
)

func NewGRPCServer(
    cfg *eagle.Config,
    userSvc *service.UserService,
    orderSvc *service.OrderService,  // 新增参数
) *grpc.Server {
    // ...

    // 注册服务
    userV1.RegisterUserServiceServer(srv, userSvc)
    orderV1.RegisterOrderServiceServer(srv, orderSvc)  // 新增注册

    return srv
}
```

#### 步骤 6: 重新生成 Wire 代码

```bash
make wire
```

Wire 会自动分析依赖关系并生成初始化代码：
- `OrderService` 依赖 `OrderRepo`
- `OrderRepo` 依赖 `DBClient`
- 所有依赖会自动注入

#### 步骤 7: 测试

```bash
# 运行服务
make run

# 测试创建订单
curl -X POST http://localhost:8080/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "product_ids": [101, 102],
    "address": "上海市浦东新区"
  }'

# 测试获取订单
curl http://localhost:8080/v1/orders/1
```

## 💡 开发最佳实践

### 1. 错误处理

使用项目定义的错误码：

```go
// internal/ecode/user.go
var (
    ErrUserNotFound          = errcode.NewError(20101, "用户不存在")
    ErrUserIsExist           = errcode.NewError(20102, "用户已存在")
    ErrPasswordIncorrect     = errcode.NewError(20103, "密码错误")
)

// 使用示例
func (s *UserService) GetUser(ctx context.Context, id int64) (*types.User, error) {
    user, err := s.repo.GetUser(ctx, id)
    if err != nil {
        return nil, err
    }
    if user == nil || user.ID == 0 {
        return nil, ecode.ErrUserNotFound  // 使用定义的错误码
    }
    return user, nil
}
```

### 2. 分层原则

- **Service 层**: 只包含业务逻辑，不直接操作数据库
- **Repository 层**: 提供数据访问接口，封装缓存策略
- **DAL 层**: 直接操作数据库、缓存、RPC

错误示例（不要这样做）：
```go
// ❌ Service 层直接使用 DAO
func (s *UserService) GetUser(ctx context.Context, id int64) {
    user, _ := dao.UserInfoModel.WithContext(ctx).First()  // 错误！
}
```

正确示例：
```go
// ✅ Service 调用 Repository
func (s *UserService) GetUser(ctx context.Context, id int64) {
    user, err := s.repo.GetUser(ctx, id)  // 正确！
}
```

### 3. 缓存策略

Repository 层统一处理缓存，避免在 Service 层操作缓存：

```go
// ✅ 在 Repository 中处理缓存
func (r *userRepo) GetUser(ctx context.Context, id int64) (*model.UserInfoModel, error) {
    // 1. 查本地缓存
    // 2. 查 Redis
    // 3. 查数据库并回写缓存
    // Service 层无需关心缓存细节
}
```

### 4. 测试编写

为每个层编写单元测试：

```go
// internal/service/user_svc_test.go
func TestUserService_CreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepo(ctrl)
    svc := NewUserService(mockRepo)

    mockRepo.EXPECT().
        CreateUser(gomock.Any(), gomock.Any()).
        Return(int64(1), nil)

    output, err := svc.CreateUser(context.Background(), types.CreateUserInput{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    })

    assert.NoError(t, err)
    assert.Equal(t, int64(1), output.ID)
}
```

运行测试：
```bash
make test
```

## 🔍 故障排查

### Wire 相关问题

**问题**: `wire: no provider found for XXX`

**原因**: 缺少依赖的 Provider

**解决**:
1. 检查对应的 `ProviderSet` 是否包含该依赖
2. 确认构造函数签名正确
3. 运行 `make wire` 重新生成

### gRPC 相关问题

**问题**: `Proto file not found`

**原因**: Proto 路径配置错误

**解决**: 检查 Makefile 中的 `PROTO_PATH` 配置

**问题**: `method XXX not implemented`

**原因**: gRPC 接口方法未实现

**解决**: 在 `*_grpc.go` 文件中实现所有 Proto 定义的方法

### GORM 相关问题

**问题**: `table not found`

**原因**: 数据库表不存在

**解决**:
1. 检查数据库迁移是否执行
2. 确认 `TableName()` 方法返回正确的表名

**问题**: `column not found`

**原因**: Model 字段和数据库列不匹配

**解决**:
1. 检查 gorm tag 是否正确
2. 运行 `make gorm-gen` 重新生成 DAO

## 📚 扩展阅读

- [Eagle 框架文档](https://go-eagle.org)
- [Google Wire 使用指南](https://github.com/google/wire)
- [gRPC Go 快速开始](https://grpc.io/docs/languages/go/quickstart/)
- [GORM 文档](https://gorm.io)
- [Clean Architecture 原则](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 🤝 获取帮助

- 项目内部文档: `CLAUDE.md`
- 代码注释: 查看相关文件的注释
- 团队协作: 向有经验的同事请教

---

**提示**: 本指南生成的 Markdown 文档包含完整的代码示例和详细说明，建议保存为参考文档供新手随时查阅。当你遇到具体问题时，也可以回来查找相应章节。

祝你在 Eagle 项目中开发愉快！🚀
