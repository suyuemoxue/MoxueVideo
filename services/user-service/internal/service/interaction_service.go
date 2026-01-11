package service

import (
	"context"
	"errors"

	"example.com/MoxueVideo/user-service/internal/repo"
	"gorm.io/gorm"
)

type InteractionService struct {
	videos    repo.VideoRepo
	likes     repo.LikeRepo
	favorites repo.FavoriteRepo
}

func NewInteractionService(videos repo.VideoRepo, likes repo.LikeRepo, favorites repo.FavoriteRepo) *InteractionService {
	return &InteractionService{videos: videos, likes: likes, favorites: favorites}
}

func (s *InteractionService) SetLike(ctx context.Context, userID uint64, videoID uint64, liked bool) error {
	if _, err := s.videos.FindByID(ctx, videoID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.likes.Set(ctx, userID, videoID, liked)
}

func (s *InteractionService) SetFavorite(ctx context.Context, userID uint64, videoID uint64, favored bool) error {
	if _, err := s.videos.FindByID(ctx, videoID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.favorites.Set(ctx, userID, videoID, favored)
}
