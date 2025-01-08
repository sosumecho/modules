package utils

import (
	md52 "crypto/md5"
	"encoding/hex"
)

func Md5(data []byte) string {
	md5 := md52.New()
	md5.Write(data)
	return hex.EncodeToString(md5.Sum([]byte(nil)))
}
