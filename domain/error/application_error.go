package error

type ApplicationError struct {
	Code ErrorCode    `json:"code"`
	Msg  string `json:"msg"`
}

var _ error = ApplicationError{}

func (a ApplicationError) Error() string {
	return a.Msg
}