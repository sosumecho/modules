package exceptions

import (
	"errors"
	"github.com/sosumecho/modules/exception"
)

const (
	// CodeVIPExpired 会员已过期
	CodeVIPExpired exception.Code = iota + 201
	// CodeCourseNotPaid 课程未购买
	CodeCourseNotPaid

	CodeExceedMaxSendTime = 301
	CodeNeedCaptcha       = 302
	CodeInvalidCaptcha    = 303

	// CodeNeedUnblock 解锁对话
	CodeNeedUnblock = 101

	CodeSettleExpired = 304
)

var (
	// VIPExpired 会员已过期
	VIPExpired = errors.New("vip expired")
	// InvalidOrderType 错误的订单类型
	InvalidOrderType = errors.New("invalid order type")
	// TooFrequent 请求太频繁
	TooFrequent = errors.New("too frequent")
	// InvalidPaymentState 错误的支付状态
	InvalidPaymentState         = errors.New("invalid payment state")
	InvalidPrice                = errors.New("invalid price")
	InvalidCourse               = errors.New("invalid course")
	InvalidGoods                = errors.New("invalid goods")
	ParamsError                 = errors.New("params error")
	SystemError                 = errors.New("system error")
	ExceedMaxSendTime           = errors.New("exceed max send times")
	NeedCaptcha                 = errors.New("captcha is needed")
	InvalidCaptcha              = errors.New("invalid captcha")
	InvalidToken                = errors.New("invalid token")
	InvalidUser                 = errors.New("invalid user")
	InvalidMobile               = errors.New("invalid mobile or password")
	InvalidEmail                = errors.New("invalid email")
	InvalidUserState            = errors.New("forbidden")
	UsernameExists              = errors.New("username already exists")
	InvalidInviteCode           = errors.New("invalid invite code")
	InsufficientBalance         = errors.New("insufficient balance")
	InvalidPassword             = errors.New("invalid password")
	ExceedMaxPasswordRetryCount = errors.New("exceed max password retry count")
	InvalidDepositAmount        = errors.New("invalid deposit amount")
	InvalidEmailDomain          = errors.New("invalid email domain")
	NeedUnblock                 = errors.New("need unblock")
	NeedInteract                = errors.New("need interact")
	InvalidUploadFileType       = errors.New("invalid upload file type")
	ExceedMaxFileSize           = errors.New("exceed max file size")
	PermissionDeny              = errors.New("permission deny")
	AlreadyPaid                 = errors.New("already paid")
	InvalidI18NLoader           = errors.New("invalid i18n loader")
	NotFound                    = errors.New("record not found")
	ContactAlreadyExists        = errors.New("contact already exists")
	NicknameAlreadyExists       = errors.New("nickname already exists")
	SettleExpired               = errors.New("settle expired")
)
