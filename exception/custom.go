package exception

type customException struct {
	*Err
	code Code
}

func (s *customException) Code() Code {
	return s.code
}

func NewException(code Code, err error) Exception {
	return &customException{
		Err:  NewError(err),
		code: code,
	}
}
