package errors

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
	CommonError                = &BizError{Code: 10001, Msg: "Unknown error", Status: fiber.StatusBadRequest}
	RegisterError              = &BizError{Code: 10002, Msg: "Register error", Status: fiber.StatusBadRequest}
	LoginPasswordError         = &BizError{Code: 10003, Msg: "Login password error", Status: fiber.StatusBadRequest}
	LoginAccountError          = &BizError{Code: 10004, Msg: "Login account error", Status: fiber.StatusBadRequest}
	MissingOrMalformedJwtError = &BizError{Code: 10005, Msg: "Missing or malformed JWT", Status: fiber.StatusBadRequest}
	InvalidOrExpiredJwtError   = &BizError{Code: 10006, Msg: "Invalid or expired JWT", Status: fiber.StatusUnauthorized}
	UserNotExistsError         = &BizError{Code: 10007, Msg: "User not exists", Status: fiber.StatusBadRequest}
	InvalidParams              = &BizError{Code: 10008, Msg: "Invalid params", Status: fiber.StatusBadRequest}
)
