package application

import (
	"github.com/damondu/greddit/domain/entity"
	"github.com/damondu/greddit/infrastructure/persistence"
	"math/rand"
)

type PostApp struct {
	postRepo persistence.PostRepo
	userRepo persistence.UserRepo
}

func NewPostApp(repositories *persistence.Repositories) *PostApp {
	return &PostApp{
		postRepo: repositories.Post,
		userRepo: repositories.User,
	}
}

func (p *PostApp) PageQueryPost(page int, pageSize int) (entity.Posts, error) {
	posts, err := p.postRepo.PageQuery(page, pageSize)
	if err != nil {
		return nil, err
	}
	return posts, err
}

func (p *PostApp) PageQueryPostUser(page int, pageSize int) (entity.PostUsers, error) {
	posts, err := p.PageQueryPost(page, pageSize)
	if err != nil {
		return nil, err
	}

	creatorUidList := posts.MapInt64(entity.GetCreatorUid)
	users, err := p.userRepo.BatchGetByUid(creatorUidList)
	if err != nil {
		return nil, err
	}

	userMap := users.GroupByInt64(entity.GetUserUid)
	postUsers := posts.MapPostUser(func(p *entity.Post) entity.PostUser {
		return entity.PostUser{
			Post: *p,
			User: userMap[p.CreatorUid],
		}
	})
	return postUsers, err
}

func (p *PostApp) Create(creatorUid int64, title, text string) (entity.Post, error) {
	return p.postRepo.Create(int64(rand.Int31()), creatorUid, title, text)
}
