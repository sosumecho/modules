package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
	"sync"
)

type customValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &customValidator{}

func (v *customValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

func (v *customValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *customValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("validate")

		for tag, vf := range customValidators {
			v.validate.RegisterValidation(tag, vf)
		}
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func Register() {
	binding.Validator = new(customValidator)
}
