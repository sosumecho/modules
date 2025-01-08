package iap

import (
	"context"
)

// 注册的iap容器
var iaps map[string]IAP

// Request 请求参数
type Request struct {
	Data                interface{} // 这个是具体请求内容
	PaymentID           string      // 这个是我们本地订单ID
	ThirdID             []string    // 这个是三方订单ID
	OriginTransactionID string      // 原始订单ID
}

// Response 内购响应
type Response struct {
	PaymentID           string // 订单ID
	ProductID           string //端口ID
	OriginTransactionID string // 原始操作ID
	TransactionID       string // 操作ID
	PurchaseDate        int64  // 购买时间
	ExpireDate          int64  // 过期时间
	Price               int    // 价格
	MemberID            string // 用户ID
	PayType             string // 支付类型
	Receipt             string // 收据
	IsRefund            *bool  // 是否是退款
	AgreementNO         string // 签约号
	AgreementType       string // 签约类型
	Platform            string // 平台
	CountryPrice        int
	CountryCode         string
	ExtraData           string // 额外的参数
	Error               error
}

// IAP is an interface to call validation API in App Store
type IAP interface {
	Verify(ctx context.Context, reqBody Request, price int) (bool, *[]Response)
	VerifyNotify(data []byte) (bool, *[]Response)
}

// RegIAP 注册
func RegIAP(name string, iap IAP) {
	iaps[name] = iap
}

// New 新建
func New(name string) IAP {
	if iap, ok := iaps[name]; ok {
		return iap
	}
	return nil
}

func init() {
	iaps = make(map[string]IAP)
}
