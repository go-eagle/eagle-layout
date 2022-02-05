package service

import (
	"context"

	"github.com/google/wire"

	pb "github.com/go-eagle/eagle-layout/api/helloworld/greeter/v1"
)

var (
	_ pb.GreeterServer = (*GreeterService)(nil)
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService)

type GreeterService struct {
	pb.UnimplementedGreeterServer
}

func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{}, nil
}
