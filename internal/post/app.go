package post

import (
	"math/rand"

	"github.com/duyike/greddit/internal/user"
)

var _ App = &app{}

type App interface {
	PageQueryPost(page int, pageSize int) (Posts, error)
	PageQueryPostUser(page int, pageSize int) (WithUsers, error)
	Create(creatorUid int64, title, text string) (Post, error)
}

type app struct {
	repository Repository
	userApp    user.App
}

func NewApp(repository2 Repository, app2 user.App) App {
	return &app{
		repository: repository2,
		userApp:    app2,
	}
}

func (a app) PageQueryPost(page int, pageSize int) (Posts, error) {
	posts, err := a.repository.PageQuery(page, pageSize)
	if err != nil {
		return nil, err
	}
	return posts, err
}

func (a app) PageQueryPostUser(page int, pageSize int) (WithUsers, error) {
	posts, err := a.PageQueryPost(page, pageSize)
	if err != nil {
		return nil, err
	}
	creatorUidList := posts.MapInt64((*Post).GetCreatorUid)
	users, err := a.userApp.BatchGetByUid(creatorUidList)
	if err != nil {
		return nil, err
	}

	userMap := users.GroupByInt64((*user.User).GetUserUid)
	postUsers := posts.MapWithUser(func(p *Post) WithUser {
		return WithUser{
			Post: *p,
			User: userMap[p.CreatorUid],
		}
	})
	return postUsers, err
}

func (a app) Create(creatorUid int64, title, text string) (Post, error) {
	return a.repository.Create(int64(rand.Int31()), creatorUid, title, text)
}
