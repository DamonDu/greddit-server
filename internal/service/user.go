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
}

type userServiceImpl struct {
}

func NewUserService() UserService {
	return (&userServiceImpl{}).init()
}

func (u userServiceImpl) QueryByUid(uid int64) (model.User, error) {
	return repository.User.GetByUid(uid)
}

func (u userServiceImpl) BatchGetByUid(uidList []int64) (model.Users, error) {
	return repository.User.BatchGetByUid(uidList)
}

func (u userServiceImpl) Register(username, email, password string) (model.User, error) {
	usernameExists, err := repository.User.UsernameExists(username)
	if err != nil {
		return model.User{}, err
	}
	if usernameExists {
		return model.User{}, errors.RegisterError.SetMsg("duplicated username")
	}

	emailExists, err := repository.User.EmailExists(email)
	if err != nil {
		return model.User{}, err
	}
	if emailExists {
		return model.User{}, errors.RegisterError.SetMsg("duplicated email")
	}
	return repository.User.Create(int64(rand.Int31()), username, email, password)
}

func (u userServiceImpl) LoginByUsername(username, password string) (model.User, error) {
	user, userErr := repository.User.GetByUsername(username)
	return u.login(&user, userErr, password)
}

func (u userServiceImpl) LoginByEmail(email, password string) (model.User, error) {
	user, userErr := repository.User.GetByEmail(email)
	return u.login(&user, userErr, password)
}

func (u userServiceImpl) login(user *model.User, err error, password string) (model.User, error) {
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

func (u userServiceImpl) init() userServiceImpl {
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
	err = repository.User.Upsert(users)
	if err != nil {
		panic(err)
	}
	return u
}
