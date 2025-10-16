# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Building and Running
- `make run` - Run the server with wire dependency injection
- `make build` - Build the binary to `bin/eagle-service` with version info and race detection
- `make wire` - Generate dependency injection code using Google Wire (generates `wire_gen.go`)

### Code Generation
- `make grpc` - Generate gRPC and Protocol Buffer code from `.proto` files
- `make proto` - Generate protocol buffer structs with validation
- `make gorm-gen` - Generate GORM model files using `cmd/gen/generate.go`

### Testing and Quality
- `make test` - Run tests with race detection
- `make lint` - Run golangci-lint for code quality checks
- `make cover` - Generate test coverage report to `coverage.txt`
- `make view-cover` - Generate and view HTML coverage report

### Documentation
- `make docs` - Generate Swagger documentation (accessible at http://localhost:8080/swagger/index.html)

## Architecture Overview

This is a Go microservice built with the Eagle framework using Clean Architecture principles:

### Layer Structure
```
Service Layer (business logic)
    ↓
Repository Layer (data access abstraction)
    ↓
DAL (Data Access Layer)
   ├── DB (database operations)
   ├── Cache (Redis caching)
   └── RPC (external service calls)
```

### Key Components

**Dependency Injection**: Uses Google Wire for compile-time dependency injection
- `cmd/server/wire.go` - Wire provider definitions
- `cmd/server/wire_gen.go` - Generated dependency injection code (run `make wire` to regenerate)

**Service Layer** (`internal/service/`):
- `*_svc.go` files - HTTP service implementations
- `*_grpc.go` files - gRPC service implementations

**Repository Layer** (`internal/repository/`):
- Provides unified data access interface
- Abstracts underlying data storage (database, cache, RPC)

**DAL Layer** (`internal/dal/`):
- `db/` - Database operations using GORM
- `cache/` - Redis caching operations  
- `rpc/` - External service communications

**Configuration**: Uses Eagle framework's config system
- Config files in `config/` directory organized by environment (`dev/`, `prod/`, `test/`)
- Loaded via `config.New()` in main.go

### Protocol Buffers and gRPC
- Proto definitions in `api/` directory
- Generated files include validation, gRPC, and HTTP gateway bindings
- Use `make grpc` after modifying `.proto` files

### Entry Points
- `cmd/server/main.go` - Main HTTP/gRPC server
- `cmd/consumer/main.go` - Background job consumer
- Both use Wire for dependency injection

### Development Workflow
1. Modify proto files if needed → `make grpc`
2. Run server → `make run`  
3. Add business logic in service/repository layers
4. Generate wire code → `make wire` (if dependencies change)
5. Test → `make test`
6. Build → `make build`