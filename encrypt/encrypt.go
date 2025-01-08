package encrypt

// EncrypterType 加密类型
type EncrypterType string

// Encrypter 加密器
type Encrypter interface {
	// Encrypt 加密
	Encrypt(data []byte) []byte
	//Decrypt 解密
	Decrypt(cipherText []byte) []byte
	// SetKey 设置加密的key
	SetKey(key interface{}) Encrypter
}

// New 返回一个加密器
func New(typ EncrypterType) Encrypter {
	switch typ {
	case RSAName:
		return NewRSA()
	case XORName:
		return NewXOR()
	}
	return nil
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
