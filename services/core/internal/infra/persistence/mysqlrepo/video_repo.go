package mysqlrepo

import (
	"context"

	"moxuevideo/core/internal/domain"
	"moxuevideo/core/internal/infra/persistence/model"

	"gorm.io/gorm"
)

type VideoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) CreateVideo(ctx context.Context, v *domain.Video) error {
	m := toModelVideo(v)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	if v != nil {
		v.ID = m.ID
		v.CreatedAt = m.CreatedAt
		v.PublishedAt = m.PublishedAt
	}
	return nil
}

func (r *VideoRepository) GetVideoByID(ctx context.Context, id uint64) (*domain.Video, error) {
	var m model.Video
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return fromModelVideo(&m), nil
}

func (r *VideoRepository) Like(ctx context.Context, userID, videoID uint64) error {
	l := model.Like{UserID: userID, VideoID: videoID}
	return r.db.WithContext(ctx).Create(&l).Error
}

func (r *VideoRepository) Unlike(ctx context.Context, userID, videoID uint64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(&model.Like{}).Error
}

func (r *VideoRepository) Favorite(ctx context.Context, userID, videoID uint64) error {
	f := model.Favorite{UserID: userID, VideoID: videoID}
	return r.db.WithContext(ctx).Create(&f).Error
}

func (r *VideoRepository) Unfavorite(ctx context.Context, userID, videoID uint64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(&model.Favorite{}).Error
}

func (r *VideoRepository) Comment(ctx context.Context, cmt *domain.Comment) error {
	return r.db.WithContext(ctx).Create(toModelComment(cmt)).Error
}

func toModelVideo(v *domain.Video) *model.Video {
	if v == nil {
		return &model.Video{}
	}
	return &model.Video{
		ID:          v.ID,
		AuthorID:    v.AuthorID,
		Title:       v.Title,
		CoverURL:    v.CoverURL,
		PlayURL:     v.PlayURL,
		CreatedAt:   v.CreatedAt,
		PublishedAt: v.PublishedAt,
	}
}

func fromModelVideo(m *model.Video) *domain.Video {
	if m == nil {
		return nil
	}
	return &domain.Video{
		ID:          m.ID,
		AuthorID:    m.AuthorID,
		Title:       m.Title,
		CoverURL:    m.CoverURL,
		PlayURL:     m.PlayURL,
		CreatedAt:   m.CreatedAt,
		PublishedAt: m.PublishedAt,
	}
}

func toModelComment(cmt *domain.Comment) *model.Comment {
	if cmt == nil {
		return &model.Comment{}
	}
	return &model.Comment{
		ID:        cmt.ID,
		UserID:    cmt.UserID,
		VideoID:   cmt.VideoID,
		Content:   cmt.Content,
		CreatedAt: cmt.CreatedAt,
	}
}
