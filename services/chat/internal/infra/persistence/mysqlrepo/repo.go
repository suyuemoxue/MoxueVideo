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

func (r *Repo) EnsureThread(ctx context.Context, userA, userB uint64) (*model.DMThread, error) {
	if userA == 0 || userB == 0 || userA == userB {
		return nil, errors.New("invalid user ids")
	}
	low, high := minmax(userA, userB)

	var t model.DMThread
	err := r.db.WithContext(ctx).
		Where("user_low_id = ? AND user_high_id = ?", low, high).
		First(&t).Error
	if err == nil {
		return &t, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	t = model.DMThread{UserLowID: low, UserHighID: high}
	if err := r.db.WithContext(ctx).Create(&t).Error; err != nil {
		var t2 model.DMThread
		if e := r.db.WithContext(ctx).
			Where("user_low_id = ? AND user_high_id = ?", low, high).
			First(&t2).Error; e == nil {
			return &t2, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repo) CreateMessage(ctx context.Context, threadID, senderID, receiverID uint64, content string, createdAt time.Time) (*model.DMMessage, error) {
	m := model.DMMessage{
		ThreadID:   threadID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  createdAt,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repo) UpdateThreadLast(ctx context.Context, threadID, messageID uint64, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.DMThread{}).
		Where("id = ?", threadID).
		Updates(map[string]any{
			"last_message_id": messageID,
			"last_message_at": at,
		}).Error
}

func minmax(a, b uint64) (uint64, uint64) {
	if a < b {
		return a, b
	}
	return b, a
}
