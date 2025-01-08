package encrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const (
	// RSAName 类型
	RSAName EncrypterType = "rsa"
)

// RSAKey Rsa加密内容
type RSAKey struct {
	PubKey     []byte // 公钥内容
	PrivateKey []byte // 私钥内容
	IsPKCS1    bool
}

// RSA rsa
type RSA struct {
	RSAKey
}

// Encrypt 加密
func (r RSA) Encrypt(plainText []byte) []byte {
	//pem解码
	block, _ := pem.Decode(r.PubKey)
	//x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	partLen := publicKey.N.BitLen()/8 - 11
	chunks := split(plainText, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunks := range chunks {
		//对明文进行加密
		b, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, chunks)
		if err != nil {
			panic(err)
		}
		buffer.Write(b)
	}
	//返回密文
	return buffer.Bytes()
}

// Decrypt 解密
// cipherText 需要解密的byte数据
// path 私钥文件路径
func (r RSA) Decrypt(cipherText []byte) []byte {
	//pem解码
	block, _ := pem.Decode(r.PrivateKey)
	//X509解码
	var (
		privateKey *rsa.PrivateKey
		err        error
	)
	if r.IsPKCS1 {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else {
		var rs interface{}
		rs, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err == nil {
			privateKey = rs.(*rsa.PrivateKey)
		}
	}
	if err != nil {
		panic(err)
	}
	partLen := privateKey.PublicKey.N.BitLen() / 8
	//log.New().Debug(partLen)
	//log.New().Debug(cipherText)
	chunks := split(cipherText, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		//对密文进行解密
		//log.New().Debug(hex.EncodeToString(chunk))
		//log.New().Debug(chunk)
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, chunk)
		if err != nil {
			return []byte{}
		}
		buffer.Write(decrypted)
	}
	//返回明文
	return buffer.Bytes()
}

// SetKey 设置key
func (r RSA) SetKey(key interface{}) Encrypter {
	pk := key.(RSAKey)
	r.PrivateKey = pk.PrivateKey
	r.PubKey = pk.PubKey
	r.IsPKCS1 = pk.IsPKCS1
	return r
}

// NewRSA 新建加密rsa
func NewRSA() Encrypter {
	return &RSA{}
}
