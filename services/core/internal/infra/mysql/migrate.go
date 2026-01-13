package mysql

import (
	"moxuevideo/core/internal/domain/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.Video{},
		&model.Like{},
		&model.Favorite{},
		&model.WatchHistory{},
	)
}
