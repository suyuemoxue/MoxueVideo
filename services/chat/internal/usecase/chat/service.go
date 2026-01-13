package chat

import (
	"context"
	"errors"
	"time"

	"moxuevideo/chat/internal/domain"
	"moxuevideo/chat/internal/infra/persistence/model"
)

type Repository interface {
	EnsureThread(ctx context.Context, userA, userB uint64) (*model.DMThread, error)
	CreateMessage(ctx context.Context, threadID, senderID, receiverID uint64, content string, createdAt time.Time) (*model.DMMessage, error)
	UpdateThreadLast(ctx context.Context, threadID, messageID uint64, at time.Time) error
}

type Publisher interface {
	PublishChatMessageCreated(evt domain.ChatMessageCreated) error
}

type Service struct {
	repo Repository
	pub  Publisher
}

func New(repo Repository, pub Publisher) *Service {
	return &Service{repo: repo, pub: pub}
}

func (s *Service) Send(ctx context.Context, senderID, receiverID uint64, threadID uint64, content string) (domain.ChatMessageCreated, error) {
	if senderID == 0 || receiverID == 0 || senderID == receiverID {
		return domain.ChatMessageCreated{}, errors.New("invalid user ids")
	}
	if content == "" {
		return domain.ChatMessageCreated{}, errors.New("empty content")
	}

	var t *model.DMThread
	var err error
	if threadID == 0 {
		t, err = s.repo.EnsureThread(ctx, senderID, receiverID)
		if err != nil {
			return domain.ChatMessageCreated{}, err
		}
		threadID = t.ID
	}

	now := time.Now()
	m, err := s.repo.CreateMessage(ctx, threadID, senderID, receiverID, content, now)
	if err != nil {
		return domain.ChatMessageCreated{}, err
	}
	_ = s.repo.UpdateThreadLast(ctx, threadID, m.ID, now)

	evt := domain.ChatMessageCreated{
		MessageID:  m.ID,
		ThreadID:   threadID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  now.UnixMilli(),
	}
	if s.pub != nil {
		_ = s.pub.PublishChatMessageCreated(evt)
	}
	return evt, nil
}
