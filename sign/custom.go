package sign

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

const (
	// CustomSignerType 自定义普通签名方式
	CustomSignerType SignerType = "custom"
)

// Custom 自定义普通签名
type Custom struct {
	Key     string `mapstructure:"key"`
	KeyName string `mapstructure:"key_name"`
}

// Sign 对数据进行签名
func (c *Custom) Sign(data map[string]interface{}) string {
	//STEP 1, 对key进行升序排序.
	sortedKeys := make([]string, 0)
	for k := range data {
		if k == c.KeyName {
			continue
		}
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sortedKeys {
		if k == c.KeyName {
			continue
		}
		value := fmt.Sprintf("%v", data[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	//STEP3, 在键值对的最后加上key=API_KEY
	if c.Key != "" {
		signStrings += fmt.Sprintf("%s=%s", c.KeyName, c.Key)
	}

	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

// Validate 校验
func (c *Custom) Validate(data map[string]interface{}) bool {
	key := data[c.KeyName]
	d := c.Sign(data)
	return d == key
}

// SetKey 设置key
func (c *Custom) SetKey(key interface{}) Signer {
	c.Key = key.(string)
	return c
}

// SetKeyName 设置keyName
func (c *Custom) SetKeyName(keyName string) Signer {
	c.KeyName = keyName
	return c
}

// NewCustomSigner 得到自定义普通签名器
func NewCustomSigner() Signer {
	return &Custom{}
}

func init() {
	RegSigner(CustomSignerType, NewCustomSigner())
}
