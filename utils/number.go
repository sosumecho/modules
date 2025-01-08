package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomNumber(prefix string, length int) string {
	var rs = prefix
	for i := 0; i < length; i++ {
		rs += fmt.Sprintf("%d", rand.Intn(10))
	}
	return rs
}
