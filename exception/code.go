package exception

type Code int

const (
	CodeParamsError Code = iota + 101
	CodeSystemError
	CodeKickOutError
	CodeVersionTooOldError
	CodeBuildTokenError
	CodeTotpError
	CodeUAError
	CodeExistError
)
