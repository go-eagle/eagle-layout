package service

import (
	"github.com/go-eagle/eagle/pkg/errcode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// convertToGrpcError 转换错误为 gRPC 错误
func convertToGrpcError(err error) error {
	// 检查是否是 gorm.ErrRecordNotFound
	if err == gorm.ErrRecordNotFound {
		return status.Errorf(codes.NotFound, "user not found")
	}

	// 获取错误码
	e, ok := err.(*errcode.Error)
	if !ok {
		// 非标准错误，包装为 Internal 错误
		return status.Errorf(codes.Internal, "internal error: %v", err)
	}

	// 转换错误码
	code := errcode.ToRPCCode(e.Code())
	return status.Errorf(code, e.Msg())
}
