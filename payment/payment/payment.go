package payment

import (
	"github.com/sosumecho/modules/payment"
	"github.com/sosumecho/modules/payment/google_play"
)

// New 新建一个支付驱动
func New(paymentType payment.Type) payment.IPayment {
	switch paymentType {
	case google_play.Name:
		return google_play.New()
		//case apple.Name:
		//	return apple.New()
	}
	return nil
}
