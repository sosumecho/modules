package exception

type UAError struct {
	*Err
}

func (e *UAError) Code() Code {
	return CodeUAError
}

func NewUAError(e error) Exception {
	return &UAError{
		Err: NewError(e),
	}
}
