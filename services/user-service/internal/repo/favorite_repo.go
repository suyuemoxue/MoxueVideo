package repo

import (
	"context"

	"example.com/MoxueVideo/user-service/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FavoriteRepo interface {
	Set(ctx context.Context, userID uint64, videoID uint64, favored bool) error
	CountByVideoIDs(ctx context.Context, videoIDs []uint64) (map[uint64]int64, error)
	FavoredMap(ctx context.Context, userID uint64, videoIDs []uint64) (map[uint64]bool, error)
	ListVideoIDsByUser(ctx context.Context, userID uint64, page int, size int) ([]uint64, error)
}

type favoriteRepo struct {
	db *gorm.DB
}

func NewFavoriteRepo(db *gorm.DB) FavoriteRepo {
	return &favoriteRepo{db: db}
}

func (r *favoriteRepo) Set(ctx context.Context, userID uint64, videoID uint64, favored bool) error {
	if favored {
		f := &model.Favorite{UserID: userID, VideoID: videoID}
		return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(f).Error
	}
	return r.db.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&model.Favorite{}).Error
}

func (r *favoriteRepo) CountByVideoIDs(ctx context.Context, videoIDs []uint64) (map[uint64]int64, error) {
	res := make(map[uint64]int64, len(videoIDs))
	if len(videoIDs) == 0 {
		return res, nil
	}

	type row struct {
		VideoID uint64
		Cnt     int64
	}

	var rows []row
	if err := r.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Select("video_id as video_id, count(*) as cnt").
		Where("video_id IN ?", videoIDs).
		Group("video_id").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, rw := range rows {
		res[rw.VideoID] = rw.Cnt
	}
	return res, nil
}

func (r *favoriteRepo) FavoredMap(ctx context.Context, userID uint64, videoIDs []uint64) (map[uint64]bool, error) {
	res := make(map[uint64]bool, len(videoIDs))
	if userID == 0 || len(videoIDs) == 0 {
		return res, nil
	}

	var ids []uint64
	if err := r.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Select("video_id").
		Where("user_id = ? AND video_id IN ?", userID, videoIDs).
		Find(&ids).Error; err != nil {
		return nil, err
	}
	for _, id := range ids {
		res[id] = true
	}
	return res, nil
}

func (r *favoriteRepo) ListVideoIDsByUser(ctx context.Context, userID uint64, page int, size int) ([]uint64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 50 {
		size = 50
	}
	offset := (page - 1) * size

	var ids []uint64
	if err := r.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Select("video_id").
		Where("user_id = ?", userID).
		Order("id DESC").
		Offset(offset).
		Limit(size).
		Find(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
