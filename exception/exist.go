package exception

type ExistError struct {
	*Err
}

func (s *ExistError) Code() Code {
	return CodeExistError
}

func NewExistError(err error) Exception {
	return &ExistError{
		Err: NewError(err),
	}
}
