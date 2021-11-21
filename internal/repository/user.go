package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/duyike/greddit/internal/model"
)

type UserRepo interface {
	Create(uid int64, username, email, password string) (model.User, error)
	GetByUid(uid int64) (model.User, error)
	BatchGetByUid(uidList []int64) (model.Users, error)
	GetByUsername(username string) (model.User, error)
	GetByEmail(email string) (model.User, error)
	UsernameExists(username string) (bool, error)
	EmailExists(email string) (bool, error)
	Upsert([]model.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r userRepo) Create(uid int64, username, email, password string) (model.User, error) {
	user := model.User{
		Uid:      uid,
		Username: username,
		Password: password,
		Email:    email,
	}
	err := r.db.Select("Uid", "Username", "Password", "Email").Create(&user).Error
	return user, err
}

func (r userRepo) GetByUid(uid int64) (model.User, error) {
	var user model.User
	err := r.db.Where("uid = ?", uid).First(&user).Error
	return user, err
}

func (r userRepo) BatchGetByUid(uidList []int64) (model.Users, error) {
	var batchUser model.Users
	err := r.db.Where("uid IN ?", uidList).Find(&batchUser).Error
	return batchUser, err
}

func (r userRepo) GetByUsername(username string) (model.User, error) {
	var user model.User
	err := r.db.Where(&model.User{Username: username}).First(&user).Error
	return user, err
}

func (r userRepo) GetByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Where(&model.User{Email: email}).First(&user).Error
	return user, err
}

func (r userRepo) UsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where(&model.User{Username: username}).Count(&count).Error
	return count > 0, err
}

func (r userRepo) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where(&model.User{Email: email}).Count(&count).Error
	return count > 0, err
}

func (r userRepo) Upsert(users []model.User) error {
	for _, user := range users {
		err := r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&user).Error
		if err != nil {
			return err
		}
	}
	return nil
}
