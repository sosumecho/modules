package google_play

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/awa/go-iap/playstore"
	jsoniter "github.com/json-iterator/go"
	"github.com/sosumecho/modules/payment"
	"github.com/sosumecho/modules/payment/iap"
	"github.com/sosumecho/modules/utils"

	"io/ioutil"

	"strings"
)

const (
	// Name 名称
	Name = "google_play"
)

type GooglePlayVerifyParam struct {
	SubscriptionID string
	ProductID      string
	Package        string
	PurchaseToken  string
}

type GooglePlay struct {
	client    *playstore.Client
	extraData GooglePlayVerifyParam
}

func (g *GooglePlay) Create(params payment.CreatePaymentParam) payment.CreatePaymentResult {
	return payment.CreatePaymentResult{}
}

func (g *GooglePlay) SetExtraData(extraData interface{}) error {
	var e GooglePlayVerifyParam
	if err := jsoniter.UnmarshalFromString(extraData.(string), &e); err != nil {
		return err
	}
	g.extraData = e
	return nil
}

func (g *GooglePlay) Refund(id string, amount int64) interface{} {
	return g.client.RefundSubscription(context.Background(), g.extraData.Package, g.extraData.SubscriptionID, g.extraData.PurchaseToken) == nil
}

func (g *GooglePlay) Verify(id string, thirdID []string, extraData interface{}, price int) (bool, *[]iap.Response) {
	data := extraData.(GooglePlayVerifyParam)
	var (
		ctx = context.TODO()
	)
	//extraDataStr, _ := jsoniter.MarshalToString(data)
	extraDataStr := data.PurchaseToken
	if data.SubscriptionID != "" {
		resp, err := g.client.VerifySubscription(ctx, data.Package, data.SubscriptionID, data.PurchaseToken)
		if err != nil /*|| (resp.PurchaseType != nil && *resp.PurchaseType == 0 && config.Env != "dev") */ {
			//log.New().WithFields(logrus.Fields{
			//	"extra_data": extraData,
			//	"err":        err,
			//}).Error("验证订阅失败,")
			return false, nil
		}
		//log.New().WithFields(logrus.Fields{
		//	"resp": resp,
		//}).Debug("验证gp数据")
		if resp.PaymentState != nil && *resp.PaymentState == 1 {
			orderIDs := strings.Split(resp.OrderId, "..")
			originOrderID := resp.OrderId
			purchaseDate := resp.StartTimeMillis / 1000
			if len(orderIDs) > 1 {
				originOrderID = orderIDs[0]
				purchaseDate = time.Now().Unix()
			}
			return true, &[]iap.Response{
				{
					ProductID:           data.ProductID,
					OriginTransactionID: originOrderID,
					TransactionID:       resp.OrderId,
					Receipt:             resp.OrderId,
					PayType:             Name,
					CountryCode:         resp.CountryCode,
					CountryPrice:        int(resp.PriceAmountMicros / 10000),
					PurchaseDate:        purchaseDate,
					ExpireDate:          resp.ExpiryTimeMillis / 1000,
					ExtraData:           extraDataStr,
				},
			}
		}
	} else {
		resp, err := g.client.VerifyProduct(ctx, data.Package, data.ProductID, data.PurchaseToken)
		if err != nil || (resp.PurchaseType != nil && *resp.PurchaseType == 0) {
			//log.New().WithFields(logrus.Fields{
			//	"extra_data": extraData,
			//	"resp":       resp,
			//}).Error("验证一次性购买失败,", err.Error())
			return false, nil
		}
		if resp.PurchaseState == 0 && resp.ConsumptionState == 0 {
			return true, &[]iap.Response{
				{
					ProductID:           data.ProductID,
					OriginTransactionID: resp.OrderId,
					TransactionID:       resp.OrderId,
					Receipt:             resp.OrderId,
					CountryCode:         resp.RegionCode,
					PayType:             Name,
					ExtraData:           extraDataStr,
				},
			}
		}
	}

	return false, nil
}

func (g *GooglePlay) IsIAP() bool {
	return true
}

func (g *GooglePlay) SetIsAPP(isAPP bool) payment.IPayment {
	return g
}

func (g *GooglePlay) VerifyNotify(req *http.Request) (bool, *[]iap.Response) {
	return false, nil
}

func (g *GooglePlay) Ack(writer http.ResponseWriter, isOK bool) {
}

func (g *GooglePlay) Sync(req *http.Request) (bool, *[]iap.Response) {
	return false, nil
}

func (g *GooglePlay) GetClient() *playstore.Client {
	return g.client
}

func New() *GooglePlay {
	jsonKey, err := ioutil.ReadFile(fmt.Sprintf("%s/configs/googleplay.json", utils.GetAbsDir()))
	if err != nil {
		panic(err)
	}
	client, err := playstore.New(jsonKey)
	if err != nil {
		panic(err)
	}

	return &GooglePlay{
		client: client,
	}
}
