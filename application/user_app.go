package application

import (
	"github.com/damondu/greddit/domain/entity"
	. "github.com/damondu/greddit/domain/error"
	"github.com/damondu/greddit/infrastructure/persistence"
	"math/rand"
	"strings"
)

type UserApp struct {
	repo persistence.UserRepo
}

func NewUserApp(repositories *persistence.Repositories) *UserApp {
	return &UserApp{repo: repositories.User}
}

func (u *UserApp) Me(uid int64) (entity.User, error) {
	return u.repo.GetByUid(uid)
}

func (u *UserApp) Register(username, password, email string) (entity.User, error) {
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
	return u.repo.Create(int64(rand.Int31()), username, email, password)
}

func (u *UserApp) LoginByUsername(username, password string) (entity.User, error) {
	user, userErr := u.repo.GetByUsername(username)
	return u.login(&user, userErr, password)
}

func (u *UserApp) LoginByEmail(email, password string) (entity.User, error) {
	user, userErr := u.repo.GetByEmail(email)
	return u.login(&user, userErr, password)
}

func (u *UserApp) login(user *entity.User, userErr error, password string) (entity.User, error) {
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
