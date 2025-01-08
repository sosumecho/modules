package sms

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	jsoniter "github.com/json-iterator/go"
)

type SmsConfig struct {
	GlobalMsg string            `mapstructure:"global-msg"`
	CNMsg     string            `mapstructure:"cn-msg"`
	Sign      string            `mapstructure:"sign"`
	AppKey    string            `mapstructure:"app-key"`
	AppSecret string            `mapstructure:"app-secret"`
	Product   bool              `mapstructure:"product"`
	Template  map[string]string `mapstructure:"template"`
}

// Send 发送验证码
func Send(conf *SmsConfig, mobile, code, msgType string) error {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", conf.AppKey, conf.AppSecret)
	if err != nil {
		return err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = mobile
	request.SignName = conf.Sign
	request.TemplateCode = msgType
	codeStr, _ := jsoniter.Marshal(map[string]string{"code": code})
	request.TemplateParam = string(codeStr)
	res, err := client.SendSms(request)
	if err != nil {
		return err
	}
	if res.Code == "OK" {
		return nil
	} else {
		return errors.New(res.Message)
	}
}

// SendByParams 发送验证码
func SendByParams(conf *SmsConfig, mobile, msgType string, params map[string]string) error {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", conf.AppKey, conf.AppSecret)
	if err != nil {
		return err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = mobile
	request.SignName = conf.Sign
	request.TemplateCode = msgType
	codeStr, _ := jsoniter.Marshal(params)
	request.TemplateParam = string(codeStr)
	res, err := client.SendSms(request)
	if err != nil {
		return err
	}
	if res.Code == "OK" {
		return nil
	} else {
		return errors.New(res.Message)
	}
}
