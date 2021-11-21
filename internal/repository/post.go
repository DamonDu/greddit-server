package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/duyike/greddit/internal/model"
	"github.com/duyike/greddit/internal/pkg/db"
)

type PostRepo interface {
	PageQuery(page int, pageSize int) (model.Posts, error)
	Create(postId, creatorUid int64, title, text string) (model.Post, error)
	Upsert([]model.Post) error
}

type postRepoImpl struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) PostRepo {
	return &postRepoImpl{db: db}
}

func (r postRepoImpl) PageQuery(page int, pageSize int) (model.Posts, error) {
	var posts model.Posts
	err := r.db.Scopes(db.Paginate(page, pageSize)).Order("updated_at desc, id desc").Find(&posts).Error
	return posts, err
}

func (r postRepoImpl) Create(postId, creatorUid int64, title, text string) (model.Post, error) {
	post := model.Post{
		PostId:     postId,
		CreatorUid: creatorUid,
		Title:      title,
		Text:       text,
	}
	err := r.db.Select("PostId", "CreatorUid", "Title", "Text").Create(&post).Error
	return post, err
}

func (r postRepoImpl) Upsert(posts []model.Post) error {
	for _, post := range posts {
		err := r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&post).Error
		if err != nil {
			return err
		}
	}
	return nil
}
