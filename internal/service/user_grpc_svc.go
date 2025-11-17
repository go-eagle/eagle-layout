package service

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	pb "github.com/go-eagle/eagle-layout/api/user/v1"
	"github.com/go-eagle/eagle-layout/internal/ecode"
	"github.com/go-eagle/eagle-layout/internal/repository"
	"github.com/go-eagle/eagle-layout/internal/types"
	"github.com/go-eagle/eagle/pkg/errcode"
)

var (
	_ pb.UserServiceServer = (*UserServiceServer)(nil)
)

// UserServiceServer gRPC 服务端
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer

	userSvc *UserService
}

// NewUserServiceServer 创建 gRPC 服务端
func NewUserServiceServer(repo repository.UserRepo) *UserServiceServer {
	return &UserServiceServer{
		userSvc: NewUserService(repo),
	}
}

// Register 注册
func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	err := req.Validate()
	if err != nil {
		return nil, ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	// 协议转换：pb.RegisterRequest → types.RegisterInput
	input := types.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// 调用业务逻辑
	result, err := s.userSvc.Register(ctx, input)
	if err != nil {
		return nil, s.convertToGrpcError(err)
	}

	// 协议转换：types.RegisterOutput → pb.RegisterReply
	return &pb.RegisterReply{
		Id:       result.ID,
		Username: result.Username,
	}, nil
}

// Login 登录
func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	if len(req.Email) == 0 && len(req.Username) == 0 {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}

	// 协议转换
	input := types.LoginInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	// 调用业务逻辑
	result, err := s.userSvc.Login(ctx, input)
	if err != nil {
		return nil, s.convertToGrpcError(err)
	}

	return &pb.LoginReply{
		Id:          result.ID,
		AccessToken: result.AccessToken,
	}, nil
}

// Logout 登出
func (s *UserServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutReply, error) {
	// 协议转换
	input := types.LogoutInput{
		ID:          req.Id,
		AccessToken: req.AccessToken,
	}

	// 调用业务逻辑
	_, err := s.userSvc.Logout(ctx, input)
	if err != nil {
		return nil, s.convertToGrpcError(err)
	}

	return &pb.LogoutReply{}, nil
}

// CreateUser 创建用户
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	// 协议转换
	input := types.CreateUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// 调用业务逻辑
	result, err := s.userSvc.CreateUser(ctx, input)
	if err != nil {
		return nil, s.convertToGrpcError(err)
	}

	return &pb.CreateUserReply{
		Id:       result.ID,
		Username: result.Username,
	}, nil
}

// UpdateUser 更新用户
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	if req.UserId == 0 {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}

	// 协议转换
	input := types.UpdateUserInput{
		UserId:   req.UserId,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
		Birthday: req.Birthday,
		Bio:      req.Bio,
		Status:   req.Status,
	}

	// 调用业务逻辑
	result, err := s.userSvc.UpdateUser(ctx, input)
	if err != nil {
		return nil, s.convertToGrpcError(err)
	}

	return &pb.UpdateUserReply{
		UserId:    result.UserId,
		Nickname:  result.Nickname,
		Phone:     result.Phone,
		Email:     result.Email,
		Avatar:    result.Avatar,
		Gender:    result.Gender,
		Birthday:  result.Birthday,
		Bio:       result.Bio,
		Status:    result.Status,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

// UpdatePassword 更新密码
func (s *UserServiceServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordReply, error) {
	if len(req.Id) == 0 {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}
	if len(req.Password) == 0 || len(req.NewPassword) == 0 || len(req.ConfirmPassword) == 0 {
		return nil, ecode.ErrInvalidArgument.Status(req).Err()
	}
	if req.NewPassword != req.ConfirmPassword {
		return nil, ecode.ErrTwicePasswordNotMatch.Status(req).Err()
	}

	// 协议转换
	input := types.UpdatePasswordInput{
		ID:              req.Id,
		Password:        req.Password,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}

	// 调用业务逻辑
	_, err := s.userSvc.UpdatePassword(ctx, input)
	if err != nil {
		return nil, s.convertToGrpcError(err)
	}

	return &pb.UpdatePasswordReply{}, nil
}

// GetUser 获取用户
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	// 协议转换
	input := types.GetUserInput{
		ID: req.Id,
	}

	// 调用业务逻辑
	result, err := s.userSvc.GetUser(ctx, input)
	if err != nil {
		return nil, s.convertToGrpcError(err)
	}

	return &pb.GetUserReply{
		User: s.convertUser(result.User),
	}, nil
}

// BatchGetUsers 批量获取用户
func (s *UserServiceServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersReply, error) {
	// 检查 RPC 请求是否取消
	if ctx.Err() == context.Canceled {
		return nil, ecode.ErrCanceled.Status(req).Err()
	}

	if len(req.GetIds()) == 0 {
		return nil, errors.New("ids is empty")
	}

	// 协议转换
	input := types.BatchGetUsersInput{
		IDs: req.GetIds(),
	}

	// 调用业务逻辑
	result, err := s.userSvc.BatchGetUsers(ctx, input)
	if err != nil {
		return nil, s.convertToGrpcError(err)
	}

	// 转换为 pb.User 列表
	var pbUsers []*pb.User
	for _, user := range result.Users {
		pbUsers = append(pbUsers, s.convertUser(user))
	}

	return &pb.BatchGetUsersReply{
		Users: pbUsers,
	}, nil
}

// convertToGrpcError 转换错误为 gRPC 错误
func (s *UserServiceServer) convertToGrpcError(err error) error {
	// 检查是否是 gorm.ErrRecordNotFound
	if err == gorm.ErrRecordNotFound {
		return ecode.ErrUserNotFound.Status(nil).Err()
	}

	// 其他错误直接返回，业务逻辑层已经处理过
	return err
}

// convertUser 转换用户模型
func (s *UserServiceServer) convertUser(u *types.User) *pb.User {
	if u == nil {
		return nil
	}
	pbUser := &pb.User{}
	_ = copier.Copy(pbUser, u)
	return pbUser
}
