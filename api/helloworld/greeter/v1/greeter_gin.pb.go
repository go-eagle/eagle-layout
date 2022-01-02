// Code generated protoc-gen-go-gin. DO NOT EDIT.
// versions:
// protoc-gen-go-gin 0.0.1

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
// context.metadata.
// gin.errors.

// app.

// errcode.

type GreeterHTTPServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterHTTPServer(r gin.IRouter, srv GreeterHTTPServer) {
	s := Greeter{
		server: srv,
		router: r,
	}
	s.RegisterService()
}

type Greeter struct {
	server GreeterHTTPServer
	router gin.IRouter
}

func (s *Greeter) SayHello_0(ctx *gin.Context) {
	var in HelloRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		app.Error(ctx, errcode.ErrBind.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(GreeterHTTPServer).SayHello(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *Greeter) RegisterService() {

	s.router.Handle("POST", "/v1/helloworld", s.SayHello_0)

}
