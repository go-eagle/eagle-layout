# Service

## 业务逻辑层说明

当业务逻辑较多并且复杂时可以在 service 目录下划分二级目录，例如:

- `service/order`
- `service/payment`

## 文件说明

为了加以区分，`http` 的 `service` 的实现方式和 `grpc` 的 `server` 实现方式根据文件后缀进行区分

- `_svc.go` 结尾的为 `http` 服务的 `service` 实现
- `_grpc_svc.go` 结尾的为 `grpc` 的 `server` 实现
