package service

import (
	"context"
	"errors"

	"example.com/MoxueVideo/user-service/internal/repo"
	"gorm.io/gorm"
)

type FollowService struct {
	users   repo.UserRepo
	follows repo.FollowRepo
}

func NewFollowService(users repo.UserRepo, follows repo.FollowRepo) *FollowService {
	return &FollowService{users: users, follows: follows}
}

func (s *FollowService) SetFollow(ctx context.Context, followerID uint64, followeeID uint64, following bool) error {
	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}
	if _, err := s.users.FindByID(ctx, followeeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.follows.Set(ctx, followerID, followeeID, following)
}
