package post

import (
	"gorm.io/gorm"

	"github.com/duyike/greddit/internal/pkg/db"
)

var _ Repository = (*repository)(nil)

type Repository interface {
	PageQuery(page int, pageSize int) (Posts, error)
	Create(postId, creatorUid int64, title, text string) (Post, error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) PageQuery(page int, pageSize int) (Posts, error) {
	var posts Posts
	err := r.db.Scopes(db.Paginate(page, pageSize)).Order("updated_at desc, id desc").Find(&posts).Error
	return posts, err
}

func (r repository) Create(postId, creatorUid int64, title, text string) (Post, error) {
	post := Post{
		PostId:     postId,
		CreatorUid: creatorUid,
		Title:      title,
		Text:       text,
	}
	err := r.db.Select("PostId", "CreatorUid", "Title", "Text").Create(&post).Error
	return post, err
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
