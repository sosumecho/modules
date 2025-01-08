package utils

import (
	"fmt"
	"github.com/spf13/cast"
	"math"
	"math/rand"
	"strings"
	"time"
)

// RandNumber 随便数字
func RandNumber(len int) string {
	rs := rand.Intn(int(math.Pow(10, float64(len))))
	//rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	//vcode := fmt.Sprintf("%0"+len+"v", rnd.Int63n(1000000000000000000))
	return fmt.Sprintf("%0"+cast.ToString(len)+"d", rs)
	//return vcode
}

//func RandNumber(len string) string {
//	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
//	vcode := fmt.Sprintf("%0"+len+"v", rnd.Int31n(1000000))
//	return vcode
//}

// ShadowMobile 隐藏中间四位手机号
func ShadowMobile(mobile string) string {
	return mobile[:3] + "****" + mobile[7:]
}

// ShadowEmail  隐藏邮箱地址信息
func ShadowEmail(email string) string {
	emailArr := strings.Split(email, "@")
	emailContent := emailArr[0]
	if len(emailContent) > 3 {
		return fmt.Sprintf("%s****@%s", emailContent[:3], emailArr[1])
	}
	return fmt.Sprintf("%s****@%s", emailContent[:], emailArr[1])
}

// GetDateFromNow 在当前的时间的基础上得到添加天数后的当前开始时间和结束时间
func GetDateFromNow(year, month, day int) (time.Time, time.Time) {
	now := time.Now().AddDate(year, month, day)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, 1)
	return start, end
}

// GetRandomString 生成随机符串
func GetRandomString(l int) string {
	rand.Seed(time.Now().UnixNano())
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.~!@#$%^&*()_+{}/?>.<,|"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// IsStrInArr 字符串是否在数组中
func IsStrInArr(arr []string, s string) bool {
	for _, item := range arr {
		if strings.Trim(s, " ") == item {
			return true
		}
	}
	return false
}

func FirstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}
