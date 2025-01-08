package encrypt

const (
	// XORName 类型
	XORName EncrypterType = "xor"
)

// XOR xor
type XOR struct {
	Key  string `mapstructure:"key"`
	bkey []byte
}

// Encrypt 加密
func (x *XOR) Encrypt(plainText []byte) []byte {
	return x.Decrypt(plainText)
}

// Decrypt 解密
func (x *XOR) Decrypt(cipherText []byte) []byte {
	keyLen := len(x.bkey)
	byteLen := len(cipherText)
	rs := make([]byte, 0)
	for i := 0; i < byteLen; i++ {
		rs = append(rs, cipherText[i]^x.bkey[i%keyLen])
	}
	return rs
}

// SetKey 设置key
func (x *XOR) SetKey(key interface{}) Encrypter {
	k := key.([]byte)
	x.bkey = k
	return x
}

// NewXOR 新建加密rsa
func NewXOR() Encrypter {
	var xor XOR
	return &xor
}
