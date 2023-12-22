package service

import (
	"context"

	pb "github.com/go-eagle/eagle-layout/api/helloworld/greeter/v1"
	"github.com/go-eagle/eagle-layout/internal/ecode"
)

var (
	_ pb.GreeterServer = (*GreeterService)(nil)
)

type GreeterService struct {
	pb.UnimplementedGreeterServer
}

func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	err := req.Validate()
	if err != nil {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}
	return &pb.HelloReply{}, nil
}
