package mysqlrepo

import (
	"context"
	"time"

	"moxuevideo/core/internal/domain"
	"moxuevideo/core/internal/infra/persistence/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, u *domain.User) error {
	m := toModelUser(u)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	if u != nil {
		u.ID = m.ID
		u.CreatedAt = m.CreatedAt
		u.UpdatedAt = m.UpdatedAt
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return fromModelUser(&m), nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).First(&m, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return fromModelUser(&m), nil
}

func (r *UserRepository) Follow(ctx context.Context, followerID, followeeID uint64) error {
	f := model.Follow{FollowerID: followerID, FolloweeID: followeeID}
	return r.db.WithContext(ctx).Create(&f).Error
}

func (r *UserRepository) Unfollow(ctx context.Context, followerID, followeeID uint64) error {
	return r.db.WithContext(ctx).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&model.Follow{}).Error
}

func (r *UserRepository) AddWatchHistory(ctx context.Context, userID, videoID uint64, watchedAt time.Time) error {
	h := model.WatchHistory{UserID: userID, VideoID: videoID, WatchedAt: watchedAt}
	return r.db.WithContext(ctx).Create(&h).Error
}

func toModelUser(u *domain.User) *model.User {
	if u == nil {
		return &model.User{}
	}
	return &model.User{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.Password,
		AvatarURL:    u.AvatarURL,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func fromModelUser(m *model.User) *domain.User {
	if m == nil {
		return nil
	}
	return &domain.User{
		ID:        m.ID,
		Username:  m.Username,
		Password:  m.PasswordHash,
		AvatarURL: m.AvatarURL,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
