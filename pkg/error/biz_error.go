package error

import "github.com/gofiber/fiber/v2"

type BizError struct {
	Code   int32  `json:"code"`
	Msg    string `json:"msg"`
	Status int    `json:"-"`
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
	CommonError        = &BizError{Code: 10001, Msg: "unknown error", Status: fiber.StatusBadRequest}
	RegisterError      = &BizError{Code: 10002, Msg: "register error", Status: fiber.StatusBadRequest}
	LoginPasswordError = &BizError{Code: 10003, Msg: "login password error", Status: fiber.StatusBadRequest}
	LoginAccountError  = &BizError{Code: 10004, Msg: "login account error", Status: fiber.StatusBadRequest}
	NoLoginError       = &BizError{Code: 10005, Msg: "no login error", Status: fiber.StatusForbidden}
)
