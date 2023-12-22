package ecode

import (
	"github.com/go-eagle/eagle/pkg/errcode"
	"google.golang.org/grpc/codes"
)

// nolint: golint
var (
	// common errors
	ErrInternalError   = errcode.New(10000, "Internal error")
	ErrInvalidArgument = errcode.New(10001, "Invalid argument")
	ErrNotFound        = errcode.New(10003, "Not found")
	ErrAccessDenied    = errcode.New(10006, "Access denied")
	ErrCanceled        = errcode.New(codes.Canceled, "RPC request is canceled")

	// biz grpc errors
	// example
	ErrUserIsExist = errcode.New(20100, "The user already exists.")
)
