package application

import (
	"github.com/damondu/greddit/domain/entity"
	. "github.com/damondu/greddit/domain/error"
	"github.com/damondu/greddit/infrastructure/persistence"
	"math/rand"
	"strings"
)

type UserApp interface {
	Me(uid int64) (entity.User, error)
	Register(username, email, password string) (entity.User, error)
	LoginByUsername(username, password string) (entity.User, error)
	LoginByEmail(email, password string) (entity.User, error)
}

type userApplication struct {
	repo persistence.UserRepo
}

var _ UserApp = &userApplication{}

func NewUserApplication(repositories *persistence.Repositories) *userApplication {
	return &userApplication{repo: repositories.User}
}

func (u *userApplication) Me(uid int64) (entity.User, error) {
	return u.repo.GetByUid(uid)
}

func (u *userApplication) Register(username, password, email string) (entity.User, error) {
	usernameExists, err := u.repo.UsernameExists(username)
	if err != nil {
		return entity.User{}, err
	}
	if usernameExists {
		return entity.User{}, &ApplicationError{Code: RegisterError, Msg: "duplicated username"}
	}

	emailExists, err := u.repo.EmailExists(email)
	if err != nil {
		return entity.User{}, err
	}
	if emailExists {
		return entity.User{}, &ApplicationError{Code: RegisterError, Msg: "duplicated email"}
	}
	return u.repo.Create(rand.Int63(), username, email, password)
}

func (u *userApplication) LoginByUsername(username, password string) (entity.User, error) {
	user, userErr := u.repo.GetByUsername(username)
	return u.login(&user, userErr, password)
}

func (u userApplication) LoginByEmail(email, password string) (entity.User, error) {
	user, userErr := u.repo.GetByEmail(email)
	return u.login(&user, userErr, password)
}

func (u userApplication) login(user *entity.User, userErr error, password string) (entity.User, error) {
	if userErr != nil {
		if strings.Contains(userErr.Error(), "record not found") {
			return entity.User{}, &ApplicationError{Code: LoginAccountError, Msg: "Account Not Exists"}
		}
		return entity.User{}, userErr
	}
	if user.Password != password {
		return entity.User{}, &ApplicationError{Code: LoginPasswordError, Msg: "Password Error"}
	}
	return *user, nil
}
