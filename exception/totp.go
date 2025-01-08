package exception

type TotpError struct {
	*Err
}

func (s *TotpError) Code() Code {
	return CodeTotpError
}

func NewTotpError(err error) Exception {
	return &TotpError{
		Err: NewError(err),
	}
}
