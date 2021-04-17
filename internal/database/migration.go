package database

import (
	"github.com/jinzhu/gorm"
	"github.com/shortdaddy0711/go-rest-api/internal/comment"
)

func MigrateDB(db *gorm.DB) error {
	if result := db.AutoMigrate(&comment.Comment{}); result.Error != nil {
		return result.Error
	}
	return nil
}
