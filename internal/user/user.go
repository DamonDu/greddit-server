package user

import "time"

type User struct {
	ID        int64     `gorm:"column:id" json:"id"`
	Uid       int64     `gorm:"column:uid" json:"uid"`
	Username  string    `gorm:"column:username" json:"username"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	Version   int64     `gorm:"column:version" json:"version"`
	Extra     string    `gorm:"column:extra" json:"extra"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (u *User) GetUserUid() int64 {
	return u.Uid
}
