package service

import (
	"context"

	pb "github.com/go-eagle/eagle-layout/api/helloworld/greeter/v1"
)

type GreeterGRPCService struct {
	pb.UnimplementedGreeterServer
}

func NewGreeterService() *GreeterGRPCService {
	return &GreeterGRPCService{}
}

func (s *GreeterGRPCService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{}, nil
}
