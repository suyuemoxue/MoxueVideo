package user

import (
	"context"
	"errors"
	"time"

	"moxuevideo/core/internal/domain"
)

var ErrNotImplemented = errors.New("not implemented")

type Repository interface {
	CreateUser(ctx context.Context, u *domain.User) error
	GetUserByID(ctx context.Context, id uint64) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	Follow(ctx context.Context, followerID, followeeID uint64) error
	Unfollow(ctx context.Context, followerID, followeeID uint64) error
	AddWatchHistory(ctx context.Context, userID, videoID uint64, watchedAt time.Time) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, username, password string) error {
	_ = ctx
	_ = username
	_ = password
	return ErrNotImplemented
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	_ = ctx
	_ = username
	_ = password
	return "", ErrNotImplemented
}

func (s *Service) Follow(ctx context.Context, followerID, followeeID uint64) error {
	_ = ctx
	_ = followerID
	_ = followeeID
	return ErrNotImplemented
}

func (s *Service) Unfollow(ctx context.Context, followerID, followeeID uint64) error {
	_ = ctx
	_ = followerID
	_ = followeeID
	return ErrNotImplemented
}

func (s *Service) RecordWatch(ctx context.Context, userID, videoID uint64, watchedAt time.Time) error {
	_ = ctx
	_ = userID
	_ = videoID
	_ = watchedAt
	return ErrNotImplemented
}
