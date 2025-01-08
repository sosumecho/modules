package pic

import (
	"context"
	"fmt"
	"github.com/sosumecho/modules/captcha/captcha"
	"github.com/sosumecho/modules/captcha/conf"
	"github.com/sosumecho/modules/captcha/pics"
	"github.com/sosumecho/modules/drivers/nsq"
	"github.com/sosumecho/modules/drivers/redis"

	"time"
)

const (
	// Name 名称
	Name captcha.Type = "pic"
)

// Pic 图片验证码
type Pic struct {
	captcha.BaseCaptcha
	Config      *conf.Config
	QueueConfig *nsq.NsqConfig
	RedisConf   *redis.Conf
}

// Send 发送验证码
func (p *Pic) Send(ip, account, subject, content string) error {
	return nil
}

// Validate 验证
func (p *Pic) Validate(account string, code string) bool {
	ctx := context.TODO()
	key := p.GetCaptchaCacheKey(account)
	id := redis.New(p.RedisConf).Get(ctx, key).Val()
	if pics.Verify(id, code) {
		redis.New(p.RedisConf).Set(ctx, p.GetValidateKey(account), "1", time.Minute*captcha.CacheMinute)
		return true
	}
	p.LogErrorNum(ctx, account)
	return false
}

// IsValidate 是否通过验证
func (p *Pic) IsValidate(account string, clear bool) bool {
	ctx := context.TODO()
	key := p.GetValidateKey(account)
	if redis.New(p.RedisConf).Get(ctx, key).Val() == "1" {
		if clear {
			redis.New(p.RedisConf).Del(ctx, key)
			redis.New(p.RedisConf).Del(ctx, p.GetValidateKey(account))
		}
		return true
	}
	return false
}

// GenerateContent 生成内容
func (p *Pic) GenerateContent(account, subject, content string) (interface{}, error) {
	// 这里生成这次需要的图片验证码的内容
	id, content, _ := pics.Generate(redis.New(p.RedisConf).Client)
	key := p.GetCaptchaCacheKey(account)
	redis.New(p.RedisConf).Set(context.TODO(), key, id, time.Minute*captcha.CacheMinute)
	return content, nil
}

// GetCaptchaCacheKey 得到验证码缓存键
func (p *Pic) GetCaptchaCacheKey(account string) string {
	return fmt.Sprintf("pic_captcha:%s", account)
}

// GetValidateKey 得到验证键
func (p *Pic) GetValidateKey(account string) string {
	return fmt.Sprintf("pic_captcha_validate:%s", account)
}
