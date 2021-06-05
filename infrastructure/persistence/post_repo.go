package persistence

import (
	"github.com/damondu/greddit/domain/entity"
	"gorm.io/gorm"
)

type PostRepo interface {
	PageQuery(page int, pageSize int) (entity.Posts, error)
	Create(postId, creatorUid int64, title, text string) (entity.Post, error)
}

var _ PostRepo = &postRepository{}

type postRepository struct {
	db *gorm.DB
}

func (p *postRepository) PageQuery(page int, pageSize int) (entity.Posts, error) {
	var posts entity.Posts
	err := p.db.Scopes(Paginate(page, pageSize)).Order("updated_at desc, id desc").Find(&posts).Error
	return posts, err
}

func (p *postRepository) Create(postId, creatorUid int64, title, text string) (entity.Post, error) {
	post := entity.Post{
		PostId:     postId,
		CreatorUid: creatorUid,
		Title:      title,
		Text:       text,
	}
	err := p.db.Select("PostId", "CreatorUid", "Title", "Text").Create(&post).Error
	return post, err
}
