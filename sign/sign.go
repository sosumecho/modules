package sign

import "sync"

var (
	mutex     = sync.Mutex{}
	signerMap = make(Signers)
)

// SignerType 签名器类型
type SignerType string

// Signer 签名器
type Signer interface {
	// Sign 签名
	Sign(map[string]interface{}) string
	// Validate 校验
	Validate(data map[string]interface{}) bool
	// SetKey 设置key
	SetKey(key interface{}) Signer
	// SetKeyName 设置key名称
	SetKeyName(keyName string) Signer
}

// Signers 签名器集合
type Signers map[SignerType]Signer

// RegSigner 注册签名器
func RegSigner(name SignerType, signer Signer) {
	mutex.Lock()
	defer mutex.Unlock()
	signerMap[name] = signer
}

// New 新建
func New(typ SignerType) Signer {
	signer, ok := signerMap[typ]
	if ok {
		return signer
	}
	return nil
}
