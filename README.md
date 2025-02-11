# go-project-layout

eagle project template

eagle: https://github.com/go-eagle/eagle

## 目录结构

```bash
.
├── api                          # proto协议定义目录
├── cmd                          # 项目主要的入口文件目录
│   ├── consumer                 # 消费者服务入口
│   ├── gen                      # 代码生成工具入口 
│   └── server                   # HTTP/GRPC 服务主入口
├── config                       # 配置文件目录
│   ├── dev                      # 开发环境配置
│   ├── test                     # 测试环境配置
│   ├── prod                     # 生产环境配置
│   └── docker                   # Docker 环境配置
├── deploy                       # 部署相关配置
│   ├── docker                   # Docker 部署配置
│   └── k8s                      # Kubernetes 部署配置
├── internal                     # 内部应用程序代码
│   ├── dal                      # 数据访问层
│   │   ├── cache               # 缓存操作
│   │   ├── db                  # 数据库操作
│   │   └── rpc                 # RPC 调用
│   ├── ecode                   # 错误码定义
│   ├── event                   # 事件处理
│   │   └── subscribe           # 消息订阅处理
│   ├── handler                 # HTTP 请求处理器
│   ├── repository             # 数据仓库层
│   ├── routers                # 路由定义
│   └── service                # 业务逻辑层
├── scripts                    # 存放shell脚本
└── third_party                # 三方proto文件
```

## 开发流程

1、修改 proto  
2、重新生成 pb及grpc: `make grpc`  
3、运行服务 `make run`  
4、确认可运行后，补充业务逻辑
