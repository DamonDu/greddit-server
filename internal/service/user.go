package service

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/duyike/greddit/internal/model"
	"github.com/duyike/greddit/internal/repository"
	"github.com/duyike/greddit/pkg/errors"
)

type UserService interface {
	QueryByUid(uid int64) (model.User, error)
	BatchGetByUid(uidList []int64) (model.Users, error)
	Register(username, email, password string) (model.User, error)
	LoginByUsername(username, password string) (model.User, error)
	LoginByEmail(email, password string) (model.User, error)
	Init()
}

type userService struct {
	repository repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repository: repo}
}

func (u userService) QueryByUid(uid int64) (model.User, error) {
	return u.repository.GetByUid(uid)
}

func (u userService) BatchGetByUid(uidList []int64) (model.Users, error) {
	return u.repository.BatchGetByUid(uidList)
}

func (u userService) Register(username, email, password string) (model.User, error) {
	usernameExists, err := u.repository.UsernameExists(username)
	if err != nil {
		return model.User{}, err
	}
	if usernameExists {
		return model.User{}, errors.RegisterError.SetMsg("duplicated username")
	}

	emailExists, err := u.repository.EmailExists(email)
	if err != nil {
		return model.User{}, err
	}
	if emailExists {
		return model.User{}, errors.RegisterError.SetMsg("duplicated email")
	}
	return u.repository.Create(int64(rand.Int31()), username, email, password)
}

func (u userService) LoginByUsername(username, password string) (model.User, error) {
	user, userErr := u.repository.GetByUsername(username)
	return u.login(&user, userErr, password)
}

func (u userService) LoginByEmail(email, password string) (model.User, error) {
	user, userErr := u.repository.GetByEmail(email)
	return u.login(&user, userErr, password)
}

func (u userService) login(user *model.User, err error, password string) (model.User, error) {
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return model.User{}, errors.LoginAccountError
		}
		return model.User{}, err
	}
	if user.Password != password {
		return model.User{}, errors.LoginPasswordError
	}
	return *user, nil
}

func (u userService) Init() {
	jsonFile, err := os.Open("./assets/users.json")
	if err != nil {
		panic(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			panic(err)
		}
	}(jsonFile)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users []model.User
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		panic(err)
	}
	err = u.repository.Upsert(users)
	if err != nil {
		panic(err)
	}
}
