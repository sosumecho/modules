package pics

import (
	"context"
	"github.com/redis/go-redis/v9"

	"fmt"
	"image/color"
	"time"

	"github.com/mojocn/base64Captcha"
)

var (
	source = "23456789qwertyuipkjhgfdsazxcvbnm"
	store  = &rdsStore{
		client: nil,
	}
)

// Generate 生成
func Generate(client *redis.Client) (id string, b64s string, err error) {
	var driver = base64Captcha.NewDriverString(
		80,
		240,
		1,
		base64Captcha.OptionShowHollowLine|base64Captcha.OptionShowSineLine|base64Captcha.OptionShowSlimeLine,
		4,
		source,
		&color.RGBA{R: 254, G: 254, B: 254, A: 254},
		nil,
		[]string{"wqy-microhei.ttc"})
	store.client = client
	captcha := base64Captcha.NewCaptcha(driver, store)
	return captcha.Generate()
}

// Verify 验证
func Verify(id, code string) bool {
	return store.Verify(id, code, true)
}

type rdsStore struct {
	client *redis.Client
}

func (r *rdsStore) Client() *redis.Client {
	//if r.client == nil {
	//	//r.client = redis2.New()
	//}
	return r.client
}

// Set 设置
func (r *rdsStore) Set(id string, value string) error {
	return r.Client().Set(context.TODO(), fmt.Sprintf("pic_captcha:%s", id), value, time.Minute*5).Err()
}

// Get 得到数据
func (r *rdsStore) Get(id string, clear bool) string {
	ctx := context.TODO()
	key := fmt.Sprintf("pic_captcha:%s", id)
	val := r.Client().Get(ctx, key).Val()
	if clear && val != "" {
		r.client.Del(ctx, key)
	}
	return val
}

// Verify 验证
func (r *rdsStore) Verify(id, answer string, clear bool) bool {
	return r.Get(id, clear) == answer
}
