package ecode

import (
	"github.com/go-eagle/eagle/pkg/errcode"
)

// ============================================
// Eagle 框架标准错误码（统一通过 ecode 包使用）
// ============================================

// nolint: golint
var (
	// 通用错误
	// Common errors
	ErrInternalServer     = errcode.ErrInternalServer
	ErrInvalidArgument    = errcode.ErrInvalidParam
	ErrUnauthorized       = errcode.ErrUnauthorized
	ErrNotFound           = errcode.ErrNotFound
	ErrUnknown            = errcode.ErrUnknown
	ErrDeadlineExceeded   = errcode.ErrDeadlineExceeded
	ErrAccessDenied       = errcode.ErrAccessDenied
	ErrLimitExceed        = errcode.ErrLimitExceed
	ErrMethodNotAllowed   = errcode.ErrMethodNotAllowed
	ErrSignParam          = errcode.ErrSignParam
	ErrValidation         = errcode.ErrValidation
	ErrDatabase           = errcode.ErrDatabase
	ErrToken              = errcode.ErrToken
	ErrInvalidToken       = errcode.ErrInvalidToken
	ErrTokenTimeout       = errcode.ErrTokenTimeout
	ErrTooManyRequests    = errcode.ErrTooManyRequests
	ErrInvalidTransaction = errcode.ErrInvalidTransaction
	ErrEncrypt            = errcode.ErrEncrypt
	ErrServiceUnavailable = errcode.ErrServiceUnavailable

	// user 业务错误码（项目特定）
	// User business errors (project-specific)
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
	ErrInternalError         = errcode.ErrInternalServer
)
