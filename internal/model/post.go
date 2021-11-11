package model

import (
	"time"
)

type Post struct {
	ID         int64     `gorm:"column:id" json:"id"`
	PostId     int64     `gorm:"column:post_id" json:"post_id"`
	CreatorUid int64     `gorm:"column:creator_uid" json:"creator_uid"`
	Title      string    `gorm:"column:title" json:"title"`
	Text       string    `gorm:"column:text" json:"text"`
	VoteCount  int64     `gorm:"column:vote_count" json:"vote_count"`
	Version    int64     `gorm:"column:version" json:"version"`
	Extra      string    `gorm:"column:extra" json:"extra"`
	DeletedAt  time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type WithUser struct {
	Post
	User
}

func (p *Post) GetCreatorUid() int64 {
	return p.CreatorUid
}

type Posts []Post

func (s *Posts) MapInt64(fc func(something *Post) int64) []int64 {
	results := make([]int64, len(*s))
	for i, something := range *s {
		results[i] = fc(&something)
	}
	return results
}

func (s *Posts) MapWithUser(fc func(something *Post) WithUser) []WithUser {
	results := make([]WithUser, len(*s))
	for i, something := range *s {
		results[i] = fc(&something)
	}
	return results
}

type WithUsers []WithUser

func (s *WithUsers) MapInterface(fc func(something *WithUser) interface{}) []interface{} {
	results := make([]interface{}, len(*s))
	for i, something := range *s {
		results[i] = fc(&something)
	}
	return results
}
