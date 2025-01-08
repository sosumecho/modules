package i18n

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/sosumecho/modules/logger"
	"golang.org/x/text/language"
	"unsafe"
)

type ValueEncoder struct {
	logger *logger.Logger
	tags   []language.Tag
}

func NewValueEncoder(logger *logger.Logger, tags []language.Tag) *ValueEncoder {
	return &ValueEncoder{logger: logger, tags: tags}
}

func (ValueEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (v ValueEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	rw.RLock()
	defer rw.RUnlock()
	val := (*string)(ptr)
	stream.WriteString(New(v.logger).Tr(v.tags, *val, nil))
}
