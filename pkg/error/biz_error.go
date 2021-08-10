package error

type BizError struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

var _ error = BizError{}

func (e BizError) Error() string {
	return e.Msg
}

func (e *BizError) SetMsg(msg string) *BizError {
	e.Msg = msg
	return e
}

var (
	CommonError        = &BizError{Code: 10001, Msg: "unknown error"}
	RegisterError      = &BizError{Code: 10002, Msg: "register error"}
	LoginPasswordError = &BizError{Code: 10003, Msg: "login password error"}
	LoginAccountError  = &BizError{Code: 10004, Msg: "login account error"}
	NoLoginError       = &BizError{Code: 10005, Msg: "no login error"}
)
