package utils

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func Snake(s string) string {
	var (
		j int
		b strings.Builder
	)
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		// Put '_' if it is not a start or end of a word, current letter is uppercase,
		// and previous is lowercase (cases like: "UserInfo"), or next letter is also
		// a lowercase and previous letter is not "_".
		if i > 0 && i < len(s)-1 && unicode.IsUpper(r) {
			if unicode.IsLower(rune(s[i-1])) ||
				j != i-1 && unicode.IsLower(rune(s[i+1])) && unicode.IsLetter(rune(s[i-1])) {
				j = i
				b.WriteString("_")
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

func Word(s string) string {
	var b strings.Builder
	s = strings.ReplaceAll(s, "_", "-")
	arr := strings.Split(s, "-")
	for _, item := range arr {
		b.WriteString(LeftUpper(strings.ToLower(item)))
	}
	return b.String()
}

// LeftLower 首字母小写
func LeftLower(s string) string {
	if len(s) > 0 {
		return strings.ToLower(string(s[0])) + s[1:]
	}
	return s
}

func LeftUpper(s string) string {
	if len(s) > 0 {
		s = strings.ToUpper(string(s[0])) + s[1:]
	}
	return s
}

func ReverseSnake(s string) string {
	arr := strings.Split(s, "_")
	var rs string
	for _, item := range arr {
		rs += LeftUpper(item)
	}

	return rs
}

var Rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(length int) string {
	str := "ABCDEFGHJKLMNPQRSTUVWYZ"
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < length; i++ {
		result = append(result, bytes[Rng.Intn(len(bytes))])
	}
	return string(result)
}

func LowerRandomString(length int) string {
	return strings.ToLower(RandomString(length))
}

type UInt64Str uint64

func (i UInt64Str) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(i), 10))
}

func (i *UInt64Str) UnmarshalJSON(b []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		if s == "" {
			*i = UInt64Str(0)
		} else {
			value, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			*i = UInt64Str(value)
		}
		return nil
	}

	return json.Unmarshal(b, (*uint64)(i))
}

func UInt64StrSliceToUint64Slice(slice []UInt64Str) []uint64 {
	var uint64slice []uint64
	for _, str := range slice {
		uint64slice = append(uint64slice, uint64(str))
	}
	return uint64slice
}
