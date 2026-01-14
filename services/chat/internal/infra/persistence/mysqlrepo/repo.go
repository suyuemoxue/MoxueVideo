package mysqlrepo

import (
	"context"
	"errors"
	"time"

	"moxuevideo/chat/internal/infra/persistence/model"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateChat(ctx context.Context, fromUserID, toUserID uint64, msgType, content, uniqued string, createdAt time.Time) (*model.Chat, error) {
	if fromUserID == 0 || toUserID == 0 || fromUserID == toUserID {
		return nil, errors.New("invalid user ids")
	}
	if content == "" {
		return nil, errors.New("empty content")
	}

	m := model.Chat{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		MsgType:    msgType,
		Content:    content,
		IsRead:     0,
		Uniqued:    uniqued,
		CreateTime: createdAt,
		IsDel:      0,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
