package service

import (
	"math/rand"

	"github.com/duyike/greddit/internal/model"
	"github.com/duyike/greddit/internal/repository"
)

var _ PostService = &postServiceImpl{}

type PostService interface {
	PageQueryPost(page int, pageSize int) (model.Posts, error)
	PageQueryPostUser(page int, pageSize int) (model.WithUsers, error)
	Create(creatorUid int64, title, text string) (model.Post, error)
}

type postServiceImpl struct {
	repository repository.PostRepo
	userApp    UserService
}

func NewPostService(repository2 repository.PostRepo, app2 UserService) PostService {
	return &postServiceImpl{
		repository: repository2,
		userApp:    app2,
	}
}

func (p postServiceImpl) PageQueryPost(page int, pageSize int) (model.Posts, error) {
	posts, err := p.repository.PageQuery(page, pageSize)
	if err != nil {
		return nil, err
	}
	return posts, err
}

func (p postServiceImpl) PageQueryPostUser(page int, pageSize int) (model.WithUsers, error) {
	posts, err := p.PageQueryPost(page, pageSize)
	if err != nil {
		return nil, err
	}
	creatorUidList := posts.MapInt64((*model.Post).GetCreatorUid)
	users, err := p.userApp.BatchGetByUid(creatorUidList)
	if err != nil {
		return nil, err
	}

	userMap := users.GroupByInt64((*model.User).GetUserUid)
	postUsers := posts.MapWithUser(func(p *model.Post) model.WithUser {
		return model.WithUser{
			Post: *p,
			User: userMap[p.CreatorUid],
		}
	})
	return postUsers, err
}

func (p postServiceImpl) Create(creatorUid int64, title, text string) (model.Post, error) {
	return p.repository.Create(int64(rand.Int31()), creatorUid, title, text)
}
