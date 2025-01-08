package captcha

import (
	"github.com/sosumecho/modules/captcha/captcha"
	"github.com/sosumecho/modules/captcha/conf"
	"github.com/sosumecho/modules/captcha/email"
	"github.com/sosumecho/modules/captcha/mobile"
	"github.com/sosumecho/modules/captcha/pic"
	"github.com/sosumecho/modules/drivers/nsq"
	"github.com/sosumecho/modules/drivers/redis"
)

type ICaptcha interface {
	Send(ip, account, subject, content string) error
	Validate(account string, code string) bool
	IsValidate(account string, clear bool) bool
	GenerateContent(account, subject, content string) (interface{}, error)
	GetCaptchaCacheKey(account string) string
	GetValidateKey(account string) string
}

// New 新建
func New(name captcha.Type, config *conf.Config, nsqConfig *nsq.NsqConfig, redisConf *redis.Conf) ICaptcha {
	switch name {
	case email.Name:
		return &email.Email{Config: config, QueueConfig: nsqConfig, RedisConf: redisConf}
	case mobile.Name:
		return &mobile.Mobile{Config: config, QueueConfig: nsqConfig, RedisConf: redisConf}
	case pic.Name:
		return &pic.Pic{Config: config, QueueConfig: nsqConfig, RedisConf: redisConf}
	default:
		return &mobile.Mobile{Config: config, QueueConfig: nsqConfig, RedisConf: redisConf}
	}
}
