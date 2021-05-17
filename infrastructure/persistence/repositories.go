package persistence

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Repositories struct {
	User UserRepo
	Post PostRepo
	db   *gorm.DB
}

func NewRepositories() (*Repositories, error) {
	db, err := gorm.Open(mysql.Open("root:123456@/greddit?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	return &Repositories{
		User: &userRepository{db: db},
		Post: &postRepository{db: db},
		db:   db,
	}, nil
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		realPage := page
		if realPage == 0 {
			realPage = 1
		}
		offset := (realPage - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
