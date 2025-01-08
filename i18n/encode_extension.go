package i18n

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/sosumecho/modules/logger"
	"golang.org/x/text/language"
)

type EncoderExtension struct {
	encoder     map[reflect2.Type]jsoniter.ValEncoder
	languageTag []language.Tag
}

func NewEncoderExtension(logger *logger.Logger, tags ...language.Tag) *EncoderExtension {
	return &EncoderExtension{
		encoder: map[reflect2.Type]jsoniter.ValEncoder{
			reflect2.TypeOf(String("")): NewValueEncoder(logger, tags),
		},
		languageTag: tags,
	}
}

// UpdateStructDescriptor No-op
func (extension EncoderExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
}

// CreateDecoder No-op
func (extension EncoderExtension) CreateDecoder(typ reflect2.Type) jsoniter.ValDecoder {
	return nil
}

// CreateEncoder get encoder from map
func (extension EncoderExtension) CreateEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	return extension.encoder[typ]
}

// CreateMapKeyDecoder No-op
func (extension EncoderExtension) CreateMapKeyDecoder(typ reflect2.Type) jsoniter.ValDecoder {
	return nil
}

// CreateMapKeyEncoder No-op
func (extension EncoderExtension) CreateMapKeyEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	return nil
}

// DecorateDecoder No-op
func (extension EncoderExtension) DecorateDecoder(typ reflect2.Type, decoder jsoniter.ValDecoder) jsoniter.ValDecoder {
	return decoder
}

// DecorateEncoder No-op
func (extension EncoderExtension) DecorateEncoder(typ reflect2.Type, encoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	return encoder
}
