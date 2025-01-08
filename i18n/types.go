package i18n

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/sosumecho/modules/logger"
	"golang.org/x/text/language"
	"sync"
)

var (
	// bundle 全局国际化组件
	//bundle = New(language.English)
	rw sync.RWMutex
)

type String string

type JSONConfig struct {
	jsoniter.API
}

//func Bundle() *I18N {
//	return bundle
//}

func NewJSON(logger *logger.Logger, tags ...language.Tag) *JSONConfig {
	var j = &JSONConfig{
		API: jsoniter.Config{
			EscapeHTML: true,
		}.Froze(),
	}

	j.API.RegisterExtension(NewEncoderExtension(logger, tags...))

	return j
}

// Store 将字典存储到内存当中
//func Store(tag, key, content string) {
//	rw.Lock()
//	defer rw.Unlock()
//
//	languageTag, err := language.Parse(tag)
//	if err != nil {
//		return
//	}
//
//	bundle.SetMessage(languageTag, key, map[string]string{"default": content})
//}
