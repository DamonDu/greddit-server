package entity

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

type Posts []Post

func GetCreatorUid(post *Post) int64 {
	return post.CreatorUid
}
