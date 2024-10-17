package service

import (
	"context"
	"errors"

	pb "github.com/go-eagle/eagle-layout/api/helloworld/greeter/v1"
	"github.com/go-eagle/eagle-layout/internal/ecode"
	"github.com/go-eagle/eagle-layout/internal/repository"
	"github.com/go-eagle/eagle/pkg/cache"
)

var (
	_ pb.GreeterServiceServer = (*GreeterServiceServer)(nil)
)

type GreeterServiceServer struct {
	pb.UnimplementedGreeterServiceServer

	repo repository.UserRepo
}

func NewGreeterServiceServer(repo repository.UserRepo) *GreeterServiceServer {
	return &GreeterServiceServer{
		repo: repo,
	}
}

func (s *GreeterServiceServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "hello " + req.Name,
	}, nil
}
func (s *GreeterServiceServer) GetUserInfo(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	userInfo, err := s.repo.GetUser(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, cache.ErrPlaceholder) {
			return nil, ecode.ErrNotFound.Status().Err()
		}
		return nil, err
	}

	if userInfo == nil {
		return nil, ecode.ErrNotFound.Status().Err()
	}

	return &pb.GetUserReply{
		User: &pb.User{
			Id:       userInfo.ID,
			Username: userInfo.Username,
			Nickname: userInfo.Nickname,
		},
	}, nil
}
