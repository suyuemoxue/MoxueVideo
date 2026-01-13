package oss

import (
	"context"
	"errors"

	"moxuevideo/core/internal/domain"
)

var ErrUnavailable = errors.New("oss sts unavailable")

type Provider interface {
	GetUploadToken(ctx context.Context, purpose string, userID uint64) (domain.UploadToken, error)
}

type Service struct {
	provider Provider
}

func New(provider Provider) *Service {
	return &Service{provider: provider}
}

func (s *Service) GetUploadToken(ctx context.Context, purpose string, userID uint64) (domain.UploadToken, error) {
	if s == nil || s.provider == nil {
		return domain.UploadToken{}, ErrUnavailable
	}
	return s.provider.GetUploadToken(ctx, purpose, userID)
}
