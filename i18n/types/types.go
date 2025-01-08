package types

import "golang.org/x/text/language"

type Loader interface {
	Load() (map[language.Tag]map[string]string, error)
}

// Unmarshal 解析
type Unmarshal func(data []byte, v interface{}) error
