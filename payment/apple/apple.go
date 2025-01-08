package apple

import (
	"context"
	"github.com/sosumecho/modules/payment/iap"
	"io"
	"net/http"
	"strconv"

	//"github.com/sosumecho/modules/config"
	"github.com/awa/go-iap/appstore"
	jsoniter "github.com/json-iterator/go"
	"github.com/sosumecho/modules/payment"
)

const (
	Name = "apple"
)

type Apple struct {
	password string
}

func (a Apple) Create(params payment.CreatePaymentParam) payment.CreatePaymentResult {
	//TODO implement me
	panic("implement me")
}

func (a Apple) Refund(id string, amount int64) interface{} {
	//TODO implement me
	panic("implement me")
}

func (a Apple) Verify(id string, thirdID []string, extraData interface{}, price int) (bool, *[]iap.Response) {
	ctx := context.Background()
	r := appstore.IAPRequest{
		ReceiptData: extraData.(string),
		Password:    a.password,
	}

	resp := &appstore.IAPResponse{}
	err := appstore.New().Verify(ctx, r, resp)
	if err != nil {
		return false, &[]iap.Response{
			{
				Error: err,
			},
		}
	}
	//respJSON, _ := jsoniter.MarshalToString(resp)
	//log.New().WithFields(logrus.Fields{
	//	"response": respJSON,
	//}).Debug("得到苹果返回的验证信息")
	if resp.Status != 0 {
		return false, nil
	}
	if resp.Environment != appstore.Production {
		return false, nil
	}
	thirdIDs := make(map[string]int)
	for _, item := range thirdID {
		thirdIDs[item] = 1
	}
	rs := make([]iap.Response, 0)
	latestReceipts := make(map[string]appstore.InApp)
	for _, item := range resp.LatestReceiptInfo {
		latestReceipts[item.TransactionID] = item
	}
	if len(latestReceipts) > 0 {
		rs = a.check(id, resp.LatestReceipt, thirdIDs, resp.LatestReceiptInfo, latestReceipts, price, resp)
		if len(rs) == 0 {
			rs = a.checkInApp(id, price, thirdIDs, resp)
		}
	} else {
		rs = a.checkInApp(id, price, thirdIDs, resp)
	}
	if len(rs) == 0 {
		return false, nil
	}
	return true, &rs
}

func (a Apple) checkInApp(id string, price int, thirdIDs map[string]int, resp *appstore.IAPResponse) []iap.Response {
	rs := make([]iap.Response, 0)
	inapps := make(map[string]appstore.InApp)
	for _, item := range resp.Receipt.InApp {
		inapps[item.TransactionID] = item
	}
	rs = a.check(id, resp.LatestReceipt, thirdIDs, resp.Receipt.InApp, inapps, price, resp)
	return rs
}

func (a Apple) check(id string, receipts string, thirdIDs map[string]int, inappItems []appstore.InApp, inapps map[string]appstore.InApp, price int, resp *appstore.IAPResponse) []iap.Response {
	rs := make([]iap.Response, 0)
	for _, item := range inappItems {
		var expireDate int64 = 0
		var purchaseDate int64 = 0
		var err error
		if len(thirdIDs) == 0 {
			if item.ExpiresDateMS != "" {
				expireDate, err = strconv.ParseInt(item.ExpiresDateMS, 10, 64)
				if err != nil {
					//log.New().WithFields(logrus.Fields{
					//	"expire_date": item.ExpiresDateMS,
					//}).Error("转换过期时间失败, ", err.Error())
				}
			}
			if item.PurchaseDateMS != "" {
				purchaseDate, err = strconv.ParseInt(item.PurchaseDateMS, 10, 64)
				if err != nil {
					//log.New().WithFields(logrus.Fields{
					//	"purchase_date": item.PurchaseDateMS,
					//}).Error("转换购买时间失败, ", err.Error())
				}
			}

			isRefund := false
			// 如果取消的时间不为空就表示用户成功退款了
			if item.CancellationDateMS != "" {
				isRefund = true
			}
			rs = append(rs, iap.Response{
				PaymentID:           id,
				ProductID:           item.ProductID,
				OriginTransactionID: item.OriginalTransactionID,
				TransactionID:       item.TransactionID,
				ExpireDate:          expireDate / 1000,
				PurchaseDate:        purchaseDate / 1000,
				Price:               price,
				PayType:             Name,
				Receipt:             receipts,
				IsRefund:            &isRefund,
			})
		} else if _, ok := thirdIDs[item.TransactionID]; ok {

			if inapps[item.TransactionID].ExpiresDateMS != "" {
				expireDate, err = strconv.ParseInt(inapps[item.TransactionID].ExpiresDateMS, 10, 64)
				if err != nil {
					//log.New().WithFields(logrus.Fields{
					//	"expire_date": inapps[item.TransactionID].ExpiresDateMS,
					//}).Error("转换过期时间失败, ", err.Error())
				}
			}
			if inapps[item.TransactionID].PurchaseDateMS != "" {
				purchaseDate, err = strconv.ParseInt(inapps[item.TransactionID].PurchaseDateMS, 10, 64)
				if err != nil {
					//log.New().WithFields(logrus.Fields{
					//	"purchase_date": inapps[item.TransactionID].PurchaseDateMS,
					//}).Error("转换购买时间失败, ", err.Error())
				}
			}

			isRefund := false
			// 如果取消的时间不为空就表示用户成功退款了
			if item.CancellationDateMS != "" {
				isRefund = true
			}
			rs = append(rs, iap.Response{
				PaymentID:           id,
				ProductID:           item.ProductID,
				OriginTransactionID: item.OriginalTransactionID,
				TransactionID:       item.TransactionID,
				ExpireDate:          expireDate / 1000,
				PurchaseDate:        purchaseDate / 1000,
				Price:               price,
				PayType:             Name,
				Receipt:             receipts,
				IsRefund:            &isRefund,
			})
		}
	}
	return rs
}

func (a Apple) IsIAP() bool {
	//TODO implement me
	return true
}

func (a Apple) SetIsAPP(isAPP bool) payment.IPayment {
	//TODO implement me
	panic("implement me")
}

func (a Apple) VerifyNotify(req *http.Request) (bool, *[]iap.Response) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return false, &[]iap.Response{
			{
				Error: err,
			},
		}
	}
	var notificationData appstore.SubscriptionNotification
	err = jsoniter.Unmarshal(data, &notificationData)
	if err != nil {
		//log.New().Error("解析通过数据失败", err.Error())
		return false, nil
	}
	// 取出用户发过来的最新的一条去验证
	transactionIDs := make([]string, 0, len(notificationData.UnifiedReceipt.LatestReceiptInfo))
	if notificationData.NotificationType == appstore.NotificationTypeCancel || notificationData.NotificationType == appstore.NotificationTypeRefund {
		for _, item := range notificationData.UnifiedReceipt.LatestReceiptInfo {
			transactionIDs = append(transactionIDs, item.TransactionID)
		}
	} else {
		transactionIDs = append(transactionIDs, notificationData.UnifiedReceipt.LatestReceiptInfo[0].TransactionID)
	}
	return a.Verify("", transactionIDs, notificationData.UnifiedReceipt.LatestReceipt, 0)
}

func (a Apple) Ack(writer http.ResponseWriter, isOK bool) {
	if isOK {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (a Apple) Sync(req *http.Request) (bool, *[]iap.Response) {
	//TODO implement me
	panic("implement me")
}

func New(password string) payment.IPayment {
	return &Apple{
		password: password,
	}
}
