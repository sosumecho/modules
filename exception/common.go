package exception

import (
	"errors"
	"net/http"
)

type Err struct {
	e error
}

func (e *Err) Error() string {
	return e.Msg()
}

func (e *Err) Msg() string {
	return e.e.Error()
}

func (e *Err) Code() int {
	return -1
}

func (e *Err) HttpStatus() int {
	return http.StatusOK
}

func NewError(e error) *Err {
	if e == nil {
		e = errors.New("unknown")
	}
	return &Err{
		e: e,
	}
}

type ParamsError struct {
	*Err
}

func (e *ParamsError) Code() Code {
	return CodeParamsError
}

func NewParamsError(e error) Exception {
	return &ParamsError{
		Err: NewError(e),
	}
}

type SystemError struct {
	*Err
}

func (e *SystemError) Code() Code {
	return CodeSystemError
}

func NewSystemError(e error) Exception {
	return &SystemError{
		Err: NewError(e),
	}
}
