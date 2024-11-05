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

	// user grpc errors
	ErrUserIsExist           = errcode.New(20100, "The user already exists.")
	ErrUserNotFound          = errcode.New(20101, "The user was not found.")
	ErrPasswordIncorrect     = errcode.New(20102, "账号或密码错误")
	ErrAreaCodeEmpty         = errcode.New(20103, "手机区号不能为空")
	ErrPhoneEmpty            = errcode.New(20104, "手机号不能为空")
	ErrGenVCode              = errcode.New(20105, "生成验证码错误")
	ErrSendSMS               = errcode.New(20106, "发送短信错误")
	ErrSendSMSTooMany        = errcode.New(20107, "已超出当日限制，请明天再试")
	ErrVerifyCode            = errcode.New(20108, "验证码错误")
	ErrEmailOrPassword       = errcode.New(20109, "邮箱或密码错误")
	ErrTwicePasswordNotMatch = errcode.New(20110, "两次密码输入不一致")
	ErrRegisterFailed        = errcode.New(20111, "注册失败")
	ErrToken                 = errcode.New(20112, "Gen token error")
	ErrEncrypt               = errcode.New(20113, "Encrypting the user password error")
)
