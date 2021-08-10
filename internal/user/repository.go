package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(uid int64, username, email, password string) (User, error)
	GetByUid(uid int64) (User, error)
	BatchGetByUid(uidList []int64) (Users, error)
	GetByUsername(username string) (User, error)
	GetByEmail(email string) (User, error)
	UsernameExists(username string) (bool, error)
	EmailExists(email string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r repository) Create(uid int64, username, email, password string) (User, error) {
	user := User{
		Uid:      uid,
		Username: username,
		Password: password,
		Email:    email,
	}
	err := r.db.Select("Uid", "Username", "Password", "Email").Create(&user).Error
	return user, err
}

func (r repository) GetByUid(uid int64) (User, error) {
	var user User
	err := r.db.Where("uid = ?", uid).First(&user).Error
	return user, err
}

func (r repository) BatchGetByUid(uidList []int64) (Users, error) {
	var batchUser Users
	err := r.db.Where("uid IN ?", uidList).Find(&batchUser).Error
	return batchUser, err
}

func (r repository) GetByUsername(username string) (User, error) {
	var user User
	err := r.db.Where(&User{Username: username}).First(&user).Error
	return user, err
}

func (r repository) GetByEmail(email string) (User, error) {
	var user User
	err := r.db.Where(&User{Email: email}).First(&user).Error
	return user, err
}

func (r repository) UsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where(&User{Username: username}).Count(&count).Error
	return count > 0, err
}

func (r repository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where(&User{Email: email}).Count(&count).Error
	return count > 0, err
}

var _ Repository = (*repository)(nil)
