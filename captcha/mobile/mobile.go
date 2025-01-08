package mobile

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/sosumecho/modules/captcha/conf"
	"github.com/sosumecho/modules/drivers/nsq"
	"strings"

	"errors"
	"fmt"
	"github.com/sosumecho/modules/captcha/captcha"
	"github.com/sosumecho/modules/drivers/redis"

	"time"
)

const (
	// Name 名称
	Name  captcha.Type = "mobile"
	Topic              = "sms_captcha"
)

// Mobile 邮箱
type Mobile struct {
	captcha.BaseCaptcha
	Config      *conf.Config
	QueueConfig *nsq.NsqConfig
	RedisConf   *redis.Conf
}

type SmsContent struct {
	RawAccount string            `json:"raw_account"`
	Account    string            `json:"account"`
	Code       string            `json:"code"`
	MsgType    string            `json:"msg_type"`
	Params     map[string]string `json:"params"`
}

// Send 发送
func (e *Mobile) Send(ip, account, subject, content string) error {
	ctx := context.TODO()
	if !redis.New(e.RedisConf).SetNX(ctx, "sms_lock_"+account, "1", time.Minute*1).Val() {
		return errors.New("too frequent")
	}
	redis.New(e.RedisConf).Set(ctx, e.GetCaptchaCacheKey(account), content, time.Minute*captcha.CacheMinute)
	rs, err := e.GenerateContent(account, subject, content)
	if err != nil {
		return err
	}
	e.LogSendNum(ctx, ip, account)
	_ = nsq.Produce(e.QueueConfig, Topic, rs.([]byte))
	return nil
}

// GenerateContent 生成内容
func (e *Mobile) GenerateContent(account, subject, content string) (interface{}, error) {
	accountInfo := strings.Split(account, ":")
	if len(accountInfo) != 2 {
		return nil, errors.New("invalid mobile")
	}
	msgType := e.Config.Sms.CNMsg
	mobile := accountInfo[1]
	if accountInfo[0] != "+86" {
		msgType = e.Config.Sms.GlobalMsg
		mobile = strings.Replace(strings.Replace(account, "+", "", 1), ":", "", 1)
	}
	rs := SmsContent{
		RawAccount: account,
		Account:    mobile,
		Code:       content,
		MsgType:    msgType,
		Params: map[string]string{
			"code": content,
		},
	}
	b, err := jsoniter.Marshal(rs)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Validate 验证
func (e *Mobile) Validate(account string, code string) bool {
	ctx := context.TODO()
	key := e.GetCaptchaCacheKey(account)
	if redis.New(e.RedisConf).Get(ctx, key).Val() == code {
		redis.New(e.RedisConf).Del(ctx, key)
		redis.New(e.RedisConf).Set(ctx, e.GetValidateKey(account), "1", time.Minute*captcha.CacheMinute)
		return true
	}
	e.LogErrorNum(ctx, account)
	return false
}

// IsValidate 是否通过验证
func (e *Mobile) IsValidate(account string, clear bool) bool {
	ctx := context.TODO()
	key := e.GetValidateKey(account)
	if redis.New(e.RedisConf).Get(ctx, key).Val() == "1" {
		if clear {
			redis.New(e.RedisConf).Del(ctx, key)
			redis.New(e.RedisConf).Del(ctx, e.GetValidateKey(account))
		}
		return true
	}
	return false
}

// GetValidateKey 得到验证的键
func (e *Mobile) GetValidateKey(account string) string {
	return fmt.Sprintf("sms_validate:%s", account)
}

// GetCaptchaCacheKey 得到验证码缓存键
func (e *Mobile) GetCaptchaCacheKey(account string) string {
	return fmt.Sprintf("sms_captcha:%s", account)
}
