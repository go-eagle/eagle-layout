package ecode

import (
	"github.com/go-eagle/eagle/pkg/errcode"
)

// nolint: golint
var (
	// common errors
	ErrInternalError   = errcode.NewError(10000, "Internal error")
	ErrInvalidArgument = errcode.NewError(10001, "Invalid argument")
	ErrNotFound        = errcode.NewError(10003, "Not found")
	ErrAccessDenied    = errcode.NewError(10006, "Access denied")
	// ErrCanceled        = errcode.NewError(codes.Canceled, "RPC request is canceled")

	// user grpc errors
	ErrUserIsExist           = errcode.NewError(20100, "The user already exists.")
	ErrUserNotFound          = errcode.NewError(20101, "The user was not found.")
	ErrPasswordIncorrect     = errcode.NewError(20102, "账号或密码错误")
	ErrAreaCodeEmpty         = errcode.NewError(20103, "手机区号不能为空")
	ErrPhoneEmpty            = errcode.NewError(20104, "手机号不能为空")
	ErrGenVCode              = errcode.NewError(20105, "生成验证码错误")
	ErrSendSMS               = errcode.NewError(20106, "发送短信错误")
	ErrSendSMSTooMany        = errcode.NewError(20107, "已超出当日限制，请明天再试")
	ErrVerifyCode            = errcode.NewError(20108, "验证码错误")
	ErrEmailOrPassword       = errcode.NewError(20109, "邮箱或密码错误")
	ErrTwicePasswordNotMatch = errcode.NewError(20110, "两次密码输入不一致")
	ErrRegisterFailed        = errcode.NewError(20111, "注册失败")
	ErrToken                 = errcode.NewError(20112, "Gen token error")
	ErrEncrypt               = errcode.NewError(20113, "Encrypting the user password error")
)
