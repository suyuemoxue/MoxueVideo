package mysql

import (
	"moxuevideo/chat/internal/infra/persistence/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.DMThread{}, &model.DMMessage{}, &model.DMMessageRead{})
}
