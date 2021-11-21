package model

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model

	PostId     int64  `gorm:"uniqueIndex" json:"post_id"`
	CreatorUid int64  `json:"creator_uid"`
	Title      string `gorm:"type:text" json:"title"`
	Text       string `gorm:"type:mediumtext" json:"text"`
	VoteCount  int64  `json:"vote_count"`
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
