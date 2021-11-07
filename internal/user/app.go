package user

import (
	"math/rand"
	"strings"

	error2 "github.com/duyike/greddit/pkg/errors"
)

type App interface {
	QueryByUid(uid int64) (User, error)
	BatchGetByUid(uidList []int64) (Users, error)
	Register(username, email, password string) (User, error)
	LoginByUsername(username, password string) (User, error)
	LoginByEmail(email, password string) (User, error)
}

type app struct {
	repository Repository
}

func NewApp(repository2 Repository) App {
	return &app{repository: repository2}
}

func (a app) QueryByUid(uid int64) (User, error) {
	return a.repository.GetByUid(uid)
}

func (a app) BatchGetByUid(uidList []int64) (Users, error) {
	return a.repository.BatchGetByUid(uidList)
}

func (a app) Register(username, email, password string) (User, error) {
	usernameExists, err := a.repository.UsernameExists(username)
	if err != nil {
		return User{}, err
	}
	if usernameExists {
		return User{}, error2.RegisterError.SetMsg("duplicated username")
	}

	emailExists, err := a.repository.EmailExists(email)
	if err != nil {
		return User{}, err
	}
	if emailExists {
		return User{}, error2.RegisterError.SetMsg("duplicated email")
	}
	return a.repository.Create(int64(rand.Int31()), username, email, password)
}

func (a app) LoginByUsername(username, password string) (User, error) {
	user, userErr := a.repository.GetByUsername(username)
	return a.login(&user, userErr, password)
}

func (a app) LoginByEmail(email, password string) (User, error) {
	user, userErr := a.repository.GetByEmail(email)
	return a.login(&user, userErr, password)
}

func (a app) login(user *User, err error, password string) (User, error) {
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return User{}, error2.LoginAccountError
		}
		return User{}, err
	}
	if user.Password != password {
		return User{}, error2.LoginPasswordError
	}
	return *user, nil
}

var _ App = &app{}
