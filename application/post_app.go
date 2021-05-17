package application

import (
	"github.com/damondu/greddit/domain/entity"
	"github.com/damondu/greddit/infrastructure/persistence"
)

type PostApp interface {
	PageQueryPost(page int, pageSize int) (entity.Posts, error)
	PageQueryPostUser(page int, pageSize int) (entity.PostUsers, error)
}

type postApplication struct {
	postRepo persistence.PostRepo
	userRepo persistence.UserRepo
}

var _ PostApp = &postApplication{}

func NewPostApplication(repositories *persistence.Repositories) *postApplication {
	return &postApplication{
		postRepo: repositories.Post,
		userRepo: repositories.User,
	}
}

func (p postApplication) PageQueryPost(page int, pageSize int) (entity.Posts, error) {
	posts, err := p.postRepo.PageQuery(page, pageSize)
	if err != nil {
		return nil, err
	}
	return posts, err
}

func (p postApplication) PageQueryPostUser(page int, pageSize int) (entity.PostUsers, error) {
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
