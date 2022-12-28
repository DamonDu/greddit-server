package db

import (
	"os"

	"gorm.io/gorm/schema"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDb() (*gorm.DB, error) {
	dsn := "greddit.db"
	if os.Getenv("DEPLOYMENT") == "PRODUCTION" {
		dsn = "file::memory:?cache=shared"
	}
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//Logger: logger.Default.LogMode(logger.Info),
	})
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
