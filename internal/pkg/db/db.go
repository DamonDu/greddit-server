package db

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDb() (*gorm.DB, error) {
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB_NAME")

	err := createDBIfNotExist(username, password, host, port, dbName)
	if err != nil {
		return nil, err
	}
	dbConn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbName,
	)
	return gorm.Open(mysql.Open(dbConn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
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

func createDBIfNotExist(username string, password string, host string, port string, dbName string) error {
	hostConn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port,
	)
	db, err := gorm.Open(mysql.Open(hostConn), &gorm.Config{})
	if err != nil {
		return err
	}
	createSQL := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4;",
		dbName,
	)
	err = db.Exec(createSQL).Error
	if err != nil {
		return err
	}
	return nil
}
