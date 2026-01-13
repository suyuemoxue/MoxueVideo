package video

import (
	"context"
	"errors"

	"moxuevideo/core/internal/domain"
)

var ErrNotImplemented = errors.New("not implemented")

type Repository interface {
	CreateVideo(ctx context.Context, v *domain.Video) error
	GetVideoByID(ctx context.Context, id uint64) (*domain.Video, error)
	Like(ctx context.Context, userID, videoID uint64) error
	Unlike(ctx context.Context, userID, videoID uint64) error
	Favorite(ctx context.Context, userID, videoID uint64) error
	Unfavorite(ctx context.Context, userID, videoID uint64) error
	Comment(ctx context.Context, cmt *domain.Comment) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Upload(ctx context.Context, authorID uint64, title, playURL, coverURL string) error {
	_ = ctx
	_ = authorID
	_ = title
	_ = playURL
	_ = coverURL
	return ErrNotImplemented
}

func (s *Service) Like(ctx context.Context, userID, videoID uint64) error {
	_ = ctx
	_ = userID
	_ = videoID
	return ErrNotImplemented
}

func (s *Service) Unlike(ctx context.Context, userID, videoID uint64) error {
	_ = ctx
	_ = userID
	_ = videoID
	return ErrNotImplemented
}

func (s *Service) Favorite(ctx context.Context, userID, videoID uint64) error {
	_ = ctx
	_ = userID
	_ = videoID
	return ErrNotImplemented
}

func (s *Service) Unfavorite(ctx context.Context, userID, videoID uint64) error {
	_ = ctx
	_ = userID
	_ = videoID
	return ErrNotImplemented
}

func (s *Service) Comment(ctx context.Context, userID, videoID uint64, content string) error {
	_ = ctx
	_ = userID
	_ = videoID
	_ = content
	return ErrNotImplemented
}
