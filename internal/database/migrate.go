package database

import (
	"github.com/bagasriatmaputera/novelku/internal/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.RefreshToken{},
		&domain.Genres{},
		&domain.Novel{},
		&domain.Chapter{},
		&domain.Comment{},
		&domain.Library{},
		&domain.Like{},
	)
}
