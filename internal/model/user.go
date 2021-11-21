package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Uid      int64  `gorm:"index" json:"uid"`
	Username string `gorm:"type:varchar(50);index" json:"username"`
	Email    string `gorm:"index" json:"email"`
	Password string `gorm:"type:varchar(50)" json:"password"`
}

func (u *User) GetUserUid() int64 {
	return u.Uid
}

type Users []User

func (s *Users) GroupByInt64(fc func(something *User) int64) map[int64]User {
	results := make(map[int64]User, len(*s))
	for _, something := range *s {
		results[fc(&something)] = something
	}
	return results
}
