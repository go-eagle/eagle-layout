---
paths:
  - "**/*.go"
  - "**/go.mod"
  - "**/go.sum"
---
# Go Testing

> This file extends [common/testing.md](../common/testing.md) with Go specific content.

## Framework

Use the standard `go test` with **table-driven tests**.

## Race Detection

Always run with the `-race` flag:

```bash
go test -race ./...
```

## Coverage

```bash
go test -cover ./...
```

## Reference

See skill: `golang-testing` for detailed Go testing patterns and helpers.

## Examples in this repository

- Table-driven user service tests with `go.uber.org/mock`: `internal/service/user_svc_auth_test.go`, `internal/service/user_svc_profile_test.go`
- gRPC adapter tests: `internal/service/user_grpc_test.go`
- Handler tests with `net/http/httptest`: `internal/handler/v1/auth_handler_test.go`
- Regenerate committed mocks with `make mockgen` (`go.uber.org/mock`, not the archived `github.com/golang/mock`)
- `make test` runs `-race` and writes coverage. `make cover-check` enforces the `internal/service` coverage floor.
- `make lint` runs the golangci-lint v2 version pinned in the Makefile.
