package persistence

import (
	"github.com/damondu/greddit/domain/entity"
	"gorm.io/gorm"
)

type PostRepo interface {
	PageQuery(page int, pageSize int) (entity.Posts, error)
}

type postRepository struct {
	db *gorm.DB
}

var _ PostRepo = &postRepository{}

func (p postRepository) PageQuery(page int, pageSize int) (entity.Posts, error) {
	var posts entity.Posts
	err := p.db.Scopes(Paginate(page, pageSize)).Order("updated_at desc").Find(&posts).Error
	return posts, err
}
