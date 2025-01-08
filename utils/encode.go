package utils

import "encoding/base32"

// Base32Encode base32编码
func Base32Encode(key []byte) string {
	return base32.StdEncoding.EncodeToString(key)
}

// Base32Decode base32解码
func Base32Decode(key string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(key)
}
