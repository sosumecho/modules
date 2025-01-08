package exception

import "net/http"

type AuthError struct {
	*Err
}

func (s *AuthError) Code() Code {
	return http.StatusUnauthorized
}

func (s *AuthError) HttpStatus() int {
	return http.StatusUnauthorized
}

func NewAuthError(err error) Exception {
	return &AuthError{
		Err: NewError(err),
	}
}

type BuildTokenError struct {
	*Err
}

func (s *BuildTokenError) Code() Code {
	return CodeBuildTokenError
}

func NewBuildTokenError(err error) Exception {
	return &BuildTokenError{
		Err: NewError(err),
	}
}
