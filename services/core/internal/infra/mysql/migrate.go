package mysql

import (
	"os"
	"strings"

	"moxuevideo/core/internal/infra/persistence/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	s := strings.TrimSpace(os.Getenv("MYSQL_AUTOMIGRATE"))
	if s == "" {
		return nil
	}
	s = strings.ToLower(s)
	if s != "1" && s != "true" && s != "yes" {
		return nil
	}

	return db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.Video{},
		&model.Like{},
		&model.Favorite{},
		&model.Comment{},
		&model.WatchHistory{},
	)
}
