package entity

import (
	"strconv"
)

type UserCookie struct {
	Uid      int64
	Name     string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

func NewUserCookie(user *User) UserCookie {
	value := ""
	maxAge := -1
	uid := int64(0)
	if user != nil {
		uid = user.Uid
		value = strconv.FormatInt(user.Uid, 10)
		maxAge = 1800
	}
	return UserCookie{
		Uid:      uid,
		Name:     "user_id",
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
}
