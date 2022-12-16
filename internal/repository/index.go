package repository

import "gorm.io/gorm"

var (
	User UserRepo
	Post PostRepo
)

func Init(db *gorm.DB) {
	User = NewUserRepo(db)
	Post = NewPostRepo(db)
}
