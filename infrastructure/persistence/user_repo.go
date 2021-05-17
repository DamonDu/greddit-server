package persistence

import (
	"github.com/damondu/greddit/domain/entity"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(uid int64, username, email, password string) (entity.User, error)
	GetByUid(uid int64) (entity.User, error)
	BatchGetByUid(uidList []int64) (entity.Users, error)
	GetByUsername(username string) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	UsernameExists(username string) (bool, error)
	EmailExists(email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

var _ UserRepo = &userRepository{}

func (u userRepository) Create(uid int64, username, password, email string) (entity.User, error) {
	user := entity.User{
		Uid:      uid,
		Username: username,
		Password: password,
		Email:    email,
	}
	err := u.db.Select("Uid", "Username", "Password", "Email").Create(user).Error
	return user, err
}

func (u userRepository) GetByUid(uid int64) (entity.User, error) {
	var user entity.User
	err := u.db.Where(&entity.User{Uid: uid}).First(&user).Error
	return user, err
}

func (u userRepository) BatchGetByUid(batchUid []int64) (entity.Users, error) {
	var batchUser entity.Users
	err := u.db.Where("uid IN ?", batchUid).Find(&batchUser).Error
	return batchUser, err
}

func (u userRepository) GetByUsername(username string) (entity.User, error) {
	var user entity.User
	err := u.db.Where(&entity.User{Username: username}).First(&user).Error
	return user, err
}

func (u userRepository) GetByEmail(email string) (entity.User, error) {
	var user entity.User
	err := u.db.Where(&entity.User{Email: email}).First(&user).Error
	return user, err
}

func (u userRepository) UsernameExists(username string) (bool, error) {
	var count int64
	err := u.db.Model(&entity.User{}).Where(&entity.User{Username: username}).Count(&count).Error
	return count > 0, err
}

func (u userRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := u.db.Model(&entity.User{}).Where(&entity.User{Email: email}).Count(&count).Error
	return count > 0, err
}
