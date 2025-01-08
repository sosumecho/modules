package captcha

import (
	"context"
	"github.com/sosumecho/modules/drivers/redis"
	"github.com/sosumecho/modules/utils"

	"fmt"
)

type Type string

// BaseCaptcha 基准验证码
type BaseCaptcha struct {
}

// Send 发送
func (b *BaseCaptcha) Send(ip, account, subject, content string) error { return nil }

// Validate 验证
func (b *BaseCaptcha) Validate(account string, code string) bool { return false }

// IsValidate 是否已经验证成功
func (b *BaseCaptcha) IsValidate(account string, clear bool) bool { return false }

// GenerateContent 生成发送的内容
func (b *BaseCaptcha) GenerateContent(account, subject, content string) (interface{}, error) {
	return "", nil
}

// GetCaptchaCacheKey 生成验证码缓存key
func (b *BaseCaptcha) GetCaptchaCacheKey(account string) string { return "" }

// GetValidateKey 得到验证key
func (b *BaseCaptcha) GetValidateKey(account string) string { return "" }

// LogSendNum 记录发送次数
func (b *BaseCaptcha) LogSendNum(ctx context.Context, ip, account string) {
	key := fmt.Sprintf("captcha_send_num:%s", account)
	redis.Redis.Incr(ctx, key)
	start, _ := utils.GetDateFromNow(0, 0, 1)
	redis.Redis.ExpireAt(ctx, key, start)
	key = fmt.Sprintf("captcha_ip_send_num:%s", ip)
	redis.Redis.Incr(ctx, key)
	redis.Redis.ExpireAt(ctx, key, start)
}

// LogErrorNum 记录错误次数
func (b *BaseCaptcha) LogErrorNum(ctx context.Context, account string) {
	key := fmt.Sprintf("captcha_error_num:%s", account)
	redis.Redis.Incr(ctx, key)
	start, _ := utils.GetDateFromNow(0, 0, 1)
	redis.Redis.ExpireAt(ctx, key, start)
}
