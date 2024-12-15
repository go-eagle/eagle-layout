package server

import (
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/grpc"

	v1 "github.com/go-eagle/eagle-layout/api/user/v1"
	"github.com/go-eagle/eagle-layout/internal/service"
)

// NewGRPCServer creates a gRPC server
// if open grpc, then add second param: svc *service.UserServiceServer
func NewGRPCServer(cfg *app.Config, svc *service.UserServiceServer) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.Network("tcp"),
		grpc.Address(cfg.GRPC.Addr),
		grpc.Timeout(cfg.GRPC.ReadTimeout),
		grpc.EnableLog(),
	)

	// register biz service
	v1.RegisterUserServiceServer(grpcServer, svc)

	return grpcServer
}
