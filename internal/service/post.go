package service

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/duyike/greddit/internal/model"
	"github.com/duyike/greddit/internal/repository"
)

var _ PostService = &postServiceImpl{}

type PostService interface {
	PageQueryPost(page int, pageSize int) (model.Posts, error)
	PageQueryPostUser(page int, pageSize int) (model.WithUsers, error)
	Create(creatorUid int64, title, text string) (model.Post, error)
	Init()
}

type postServiceImpl struct {
	repository repository.PostRepo
	userApp    UserService
}

func NewPostService(repo repository.PostRepo, app2 UserService) PostService {
	return &postServiceImpl{
		repository: repo,
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

func (p postServiceImpl) Init() {
	jsonFile, err := os.Open("./assets/posts.json")
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

	var posts []model.Post
	err = json.Unmarshal(byteValue, &posts)
	if err != nil {
		panic(err)
	}
	err = p.repository.Upsert(posts)
	if err != nil {
		panic(err)
	}
}
