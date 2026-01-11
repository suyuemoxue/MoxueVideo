package repo

import (
	"context"

	"example.com/MoxueVideo/user-service/internal/model"
	"gorm.io/gorm"
)

type VideoRepo interface {
	Create(ctx context.Context, video *model.Video) error
	FindByID(ctx context.Context, id uint64) (*model.Video, error)
	ListFeed(ctx context.Context, cursor uint64, limit int) ([]model.Video, error)
	ListByAuthor(ctx context.Context, authorID uint64, page int, size int) ([]model.Video, error)
}

type videoRepo struct {
	db *gorm.DB
}

func NewVideoRepo(db *gorm.DB) VideoRepo {
	return &videoRepo{db: db}
}

func (r *videoRepo) Create(ctx context.Context, video *model.Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

func (r *videoRepo) FindByID(ctx context.Context, id uint64) (*model.Video, error) {
	var v model.Video
	if err := r.db.WithContext(ctx).First(&v, id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *videoRepo) ListFeed(ctx context.Context, cursor uint64, limit int) ([]model.Video, error) {
	if limit <= 0 {
		limit = 20
	}

	q := r.db.WithContext(ctx).Model(&model.Video{}).Order("id DESC").Limit(limit)
	if cursor > 0 {
		q = q.Where("id < ?", cursor)
	}

	var videos []model.Video
	if err := q.Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func (r *videoRepo) ListByAuthor(ctx context.Context, authorID uint64, page int, size int) ([]model.Video, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 50 {
		size = 50
	}

	var videos []model.Video
	offset := (page - 1) * size
	if err := r.db.WithContext(ctx).
		Model(&model.Video{}).
		Where("author_id = ?", authorID).
		Order("id DESC").
		Offset(offset).
		Limit(size).
		Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}
