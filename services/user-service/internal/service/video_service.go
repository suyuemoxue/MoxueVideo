package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"example.com/MoxueVideo/user-service/internal/model"
	"example.com/MoxueVideo/user-service/internal/repo"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type VideoService struct {
	videos    repo.VideoRepo
	users     repo.UserRepo
	likes     repo.LikeRepo
	favorites repo.FavoriteRepo
	follows   repo.FollowRepo
}

func NewVideoService(videos repo.VideoRepo, users repo.UserRepo, likes repo.LikeRepo, favorites repo.FavoriteRepo, follows repo.FollowRepo) *VideoService {
	return &VideoService{videos: videos, users: users, likes: likes, favorites: favorites, follows: follows}
}

type PublishInput struct {
	PlayURL     string
	CoverURL    string
	Title       string
	Description string
}

func (s *VideoService) Publish(ctx context.Context, authorID uint64, in PublishInput) (*model.Video, error) {
	v := &model.Video{
		AuthorID:    authorID,
		PlayURL:     strings.TrimSpace(in.PlayURL),
		CoverURL:    strings.TrimSpace(in.CoverURL),
		Title:       strings.TrimSpace(in.Title),
		Description: strings.TrimSpace(in.Description),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.videos.Create(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *VideoService) Feed(ctx context.Context, viewerID uint64, cursor uint64, limit int) ([]VideoDTO, error) {
	videos, err := s.videos.ListFeed(ctx, cursor, limit)
	if err != nil {
		return nil, err
	}
	return s.buildVideoDTOs(ctx, viewerID, videos)
}

func (s *VideoService) Get(ctx context.Context, viewerID uint64, id uint64) (*VideoDTO, error) {
	v, err := s.videos.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	res, err := s.buildVideoDTOs(ctx, viewerID, []model.Video{*v})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}
	return &res[0], nil
}

func (s *VideoService) ListByAuthor(ctx context.Context, viewerID uint64, authorID uint64, page int, size int) ([]VideoDTO, error) {
	videos, err := s.videos.ListByAuthor(ctx, authorID, page, size)
	if err != nil {
		return nil, err
	}
	return s.buildVideoDTOs(ctx, viewerID, videos)
}

type UserDTO struct {
	ID             uint64 `json:"id"`
	Username       string `json:"username"`
	DisplayName    string `json:"displayName"`
	AvatarURL      string `json:"avatarUrl"`
	FollowingCount int64  `json:"followingCount"`
	FollowerCount  int64  `json:"followerCount"`
	IsFollowing    bool   `json:"isFollowing"`
}

type VideoDTO struct {
	ID            uint64    `json:"id"`
	Author        UserDTO   `json:"author"`
	PlayURL       string    `json:"playUrl"`
	CoverURL      string    `json:"coverUrl"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	LikeCount     int64     `json:"likeCount"`
	FavoriteCount int64     `json:"favoriteCount"`
	IsLiked       bool      `json:"isLiked"`
	IsFavorited   bool      `json:"isFavorited"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (s *VideoService) buildVideoDTOs(ctx context.Context, viewerID uint64, videos []model.Video) ([]VideoDTO, error) {
	if len(videos) == 0 {
		return []VideoDTO{}, nil
	}

	videoIDs := make([]uint64, 0, len(videos))
	authorIDsSet := make(map[uint64]struct{}, len(videos))
	for _, v := range videos {
		videoIDs = append(videoIDs, v.ID)
		authorIDsSet[v.AuthorID] = struct{}{}
	}

	authorIDs := make([]uint64, 0, len(authorIDsSet))
	for id := range authorIDsSet {
		authorIDs = append(authorIDs, id)
	}

	authors, err := s.users.FindByIDs(ctx, authorIDs)
	if err != nil {
		return nil, err
	}
	authorMap := make(map[uint64]model.User, len(authors))
	for _, u := range authors {
		authorMap[u.ID] = u
	}

	followersCount, err := s.follows.CountFollowersByUserIDs(ctx, authorIDs)
	if err != nil {
		return nil, err
	}
	followingCount, err := s.follows.CountFollowingByUserIDs(ctx, authorIDs)
	if err != nil {
		return nil, err
	}
	isFollowingMap, err := s.follows.IsFollowingMap(ctx, viewerID, authorIDs)
	if err != nil {
		return nil, err
	}

	likeCount, err := s.likes.CountByVideoIDs(ctx, videoIDs)
	if err != nil {
		return nil, err
	}
	favCount, err := s.favorites.CountByVideoIDs(ctx, videoIDs)
	if err != nil {
		return nil, err
	}
	likedMap, err := s.likes.LikedMap(ctx, viewerID, videoIDs)
	if err != nil {
		return nil, err
	}
	favoredMap, err := s.favorites.FavoredMap(ctx, viewerID, videoIDs)
	if err != nil {
		return nil, err
	}

	out := make([]VideoDTO, 0, len(videos))
	for _, v := range videos {
		author := authorMap[v.AuthorID]
		out = append(out, VideoDTO{
			ID: v.ID,
			Author: UserDTO{
				ID:             author.ID,
				Username:       author.Username,
				DisplayName:    author.DisplayName,
				AvatarURL:      author.AvatarURL,
				FollowerCount:  followersCount[author.ID],
				FollowingCount: followingCount[author.ID],
				IsFollowing:    isFollowingMap[author.ID],
			},
			PlayURL:       v.PlayURL,
			CoverURL:      v.CoverURL,
			Title:         v.Title,
			Description:   v.Description,
			LikeCount:     likeCount[v.ID],
			FavoriteCount: favCount[v.ID],
			IsLiked:       likedMap[v.ID],
			IsFavorited:   favoredMap[v.ID],
			CreatedAt:     v.CreatedAt,
		})
	}

	return out, nil
}
