// Code generated protoc-gen-go-gin. DO NOT EDIT.
// protoc-gen-go-gin 0.0.14

package v1

import (
	context "context"

	gin "github.com/gin-gonic/gin"
	app "github.com/go-eagle/eagle/pkg/app"
	errcode "github.com/go-eagle/eagle/pkg/errcode"
	"github.com/spf13/cast"
	metadata "google.golang.org/grpc/metadata"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the eagle package it is being compiled against.

// context.
// metadata.
// gin.app.errcode.

type UserServiceHTTPServer interface {
	BatchGetUsers(context.Context, *BatchGetUsersRequest) (*BatchGetUsersReply, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserReply, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserReply, error)
	Login(context.Context, *LoginRequest) (*LoginReply, error)
	Logout(context.Context, *LogoutRequest) (*LogoutReply, error)
	Register(context.Context, *RegisterRequest) (*RegisterReply, error)
	UpdatePassword(context.Context, *UpdatePasswordRequest) (*UpdatePasswordReply, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserReply, error)
}

func RegisterUserServiceHTTPServer(r gin.IRouter, srv UserServiceHTTPServer) {
	s := UserService{
		server: srv,
		router: r,
	}
	s.RegisterService()
}

type UserService struct {
	server UserServiceHTTPServer
	router gin.IRouter
}

func (s *UserService) Register_0(ctx *gin.Context) {
	var in RegisterRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(UserServiceHTTPServer).Register(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *UserService) Login_0(ctx *gin.Context) {
	var in LoginRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(UserServiceHTTPServer).Login(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *UserService) Logout_0(ctx *gin.Context) {
	var in LogoutRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(UserServiceHTTPServer).Logout(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *UserService) CreateUser_0(ctx *gin.Context) {
	var in CreateUserRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(UserServiceHTTPServer).CreateUser(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *UserService) GetUser_0(ctx *gin.Context) {
	var in GetUserRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	// make sure the uri include :id
	in.Id = cast.ToInt64(ctx.Param("id"))

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(UserServiceHTTPServer).GetUser(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *UserService) BatchGetUsers_0(ctx *gin.Context) {
	var in BatchGetUsersRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(UserServiceHTTPServer).BatchGetUsers(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *UserService) UpdateUser_0(ctx *gin.Context) {
	var in UpdateUserRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(UserServiceHTTPServer).UpdateUser(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *UserService) UpdatePassword_0(ctx *gin.Context) {
	var in UpdatePasswordRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(UserServiceHTTPServer).UpdatePassword(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *UserService) RegisterService() {
	s.router.Handle("POST", "/v1/auth/register", s.Register_0)
	s.router.Handle("POST", "/v1/auth/login", s.Login_0)
	s.router.Handle("POST", "/v1/auth/logout", s.Logout_0)
	s.router.Handle("POST", "/v1/users/", s.CreateUser_0)
	s.router.Handle("GET", "/v1/users/:id", s.GetUser_0)
	s.router.Handle("GET", "/v1/users/batch", s.BatchGetUsers_0)
	s.router.Handle("PUT", "/v1/users", s.UpdateUser_0)
	s.router.Handle("POST", "/v1/users/password", s.UpdatePassword_0)
}
