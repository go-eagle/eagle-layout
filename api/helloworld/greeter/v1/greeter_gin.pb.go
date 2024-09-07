// Code generated protoc-gen-go-gin. DO NOT EDIT.
// protoc-gen-go-gin 0.0.14

package v1

import (
	context "context"
	gin "github.com/gin-gonic/gin"
	app "github.com/go-eagle/eagle/pkg/app"
	errcode "github.com/go-eagle/eagle/pkg/errcode"
	metadata "google.golang.org/grpc/metadata"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the eagle package it is being compiled against.

// context.
// metadata.
// gin.app.errcode.

type GreeterServiceHTTPServer interface {
	GetUserInfo(context.Context, *GetUserRequest) (*GetUserReply, error)
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterServiceHTTPServer(r gin.IRouter, srv GreeterServiceHTTPServer) {
	s := GreeterService{
		server: srv,
		router: r,
	}
	s.RegisterService()
}

type GreeterService struct {
	server GreeterServiceHTTPServer
	router gin.IRouter
}

func (s *GreeterService) SayHello_0(ctx *gin.Context) {
	var in HelloRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(GreeterServiceHTTPServer).SayHello(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *GreeterService) GetUserInfo_0(ctx *gin.Context) {
	var in GetUserRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(GreeterServiceHTTPServer).GetUserInfo(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *GreeterService) RegisterService() {
	s.router.Handle("GET", "/v1/helloworld", s.SayHello_0)
	s.router.Handle("GET", "/v1/users/info", s.GetUserInfo_0)
}
