package validator

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	customValidators = make(map[string]validator.Func)
	mutex            = &sync.Mutex{}
)

func reg(name string, f validator.Func) {
	mutex.Lock()
	defer mutex.Unlock()
	customValidators[name] = f
}
