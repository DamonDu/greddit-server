package error

type ErrorCode int32

const (
	CommonError        ErrorCode = 10000
	RegisterError      ErrorCode = 10001
	LoginPasswordError ErrorCode = 10002
	LoginAccountError  ErrorCode = 10003
	NoLoginError       ErrorCode = 10004
)
