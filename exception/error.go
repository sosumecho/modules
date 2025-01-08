package exception

type Exception interface {
	error
	Code() Code
	Msg() string
	HttpStatus() int
}
