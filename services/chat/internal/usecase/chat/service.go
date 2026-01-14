package chat

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"moxuevideo/chat/internal/domain"
	"moxuevideo/chat/internal/infra/persistence/model"
)

type Repository interface {
	CreateChat(ctx context.Context, fromUserID, toUserID uint64, msgType, content, uniqued string, createdAt time.Time) (*model.Chat, error)
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

func (s *Service) Send(ctx context.Context, senderID, receiverID uint64, msgType, content, uniqued string) (domain.ChatMessageCreated, error) {
	if senderID == 0 || receiverID == 0 || senderID == receiverID {
		return domain.ChatMessageCreated{}, errors.New("invalid user ids")
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return domain.ChatMessageCreated{}, errors.New("empty content")
	}

	msgType = strings.ToLower(strings.TrimSpace(msgType))
	if msgType == "" {
		msgType = "text"
	}
	switch msgType {
	case "text", "picture", "audio":
	default:
		return domain.ChatMessageCreated{}, errors.New("invalid msg_type")
	}

	uniqued = strings.TrimSpace(uniqued)
	if uniqued == "" {
		uniqued = newUniqued()
	}

	now := time.Now()
	m, err := s.repo.CreateChat(ctx, senderID, receiverID, msgType, content, uniqued, now)
	if err != nil {
		return domain.ChatMessageCreated{}, err
	}

	evt := domain.ChatMessageCreated{
		MessageID:  m.ID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		MsgType:    msgType,
		Content:    content,
		Uniqued:    uniqued,
		CreatedAt:  now.UnixMilli(),
	}
	if s.pub != nil {
		_ = s.pub.PublishChatMessageCreated(evt)
	}
	return evt, nil
}

func newUniqued() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
