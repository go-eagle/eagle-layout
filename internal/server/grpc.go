package server

import (
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
)

// NewGRPCServer creates a gRPC server
func NewGRPCServer(cfg *app.Config) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.Network("tcp"),
		grpc.Address(cfg.GRPC.Addr),
		grpc.Timeout(cfg.GRPC.ReadTimeout),
		grpc.EnableLog(),
	)

	// register biz service

	return grpcServer
}
