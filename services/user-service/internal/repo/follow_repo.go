package repo

import (
	"context"

	"example.com/MoxueVideo/user-service/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FollowRepo interface {
	Set(ctx context.Context, followerID uint64, followeeID uint64, following bool) error
	CountFollowersByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64]int64, error)
	CountFollowingByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64]int64, error)
	IsFollowingMap(ctx context.Context, followerID uint64, followeeIDs []uint64) (map[uint64]bool, error)
	ListFollowingIDs(ctx context.Context, userID uint64, page int, size int) ([]uint64, error)
	ListFollowerIDs(ctx context.Context, userID uint64, page int, size int) ([]uint64, error)
}

type followRepo struct {
	db *gorm.DB
}

func NewFollowRepo(db *gorm.DB) FollowRepo {
	return &followRepo{db: db}
}

func (r *followRepo) Set(ctx context.Context, followerID uint64, followeeID uint64, following bool) error {
	if following {
		f := &model.Follow{FollowerID: followerID, FolloweeID: followeeID}
		return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(f).Error
	}
	return r.db.WithContext(ctx).Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&model.Follow{}).Error
}

func (r *followRepo) CountFollowersByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64]int64, error) {
	res := make(map[uint64]int64, len(userIDs))
	if len(userIDs) == 0 {
		return res, nil
	}

	type row struct {
		UserID uint64
		Cnt    int64
	}

	var rows []row
	if err := r.db.WithContext(ctx).
		Model(&model.Follow{}).
		Select("followee_id as user_id, count(*) as cnt").
		Where("followee_id IN ?", userIDs).
		Group("followee_id").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, rw := range rows {
		res[rw.UserID] = rw.Cnt
	}
	return res, nil
}

func (r *followRepo) CountFollowingByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64]int64, error) {
	res := make(map[uint64]int64, len(userIDs))
	if len(userIDs) == 0 {
		return res, nil
	}

	type row struct {
		UserID uint64
		Cnt    int64
	}

	var rows []row
	if err := r.db.WithContext(ctx).
		Model(&model.Follow{}).
		Select("follower_id as user_id, count(*) as cnt").
		Where("follower_id IN ?", userIDs).
		Group("follower_id").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, rw := range rows {
		res[rw.UserID] = rw.Cnt
	}
	return res, nil
}

func (r *followRepo) IsFollowingMap(ctx context.Context, followerID uint64, followeeIDs []uint64) (map[uint64]bool, error) {
	res := make(map[uint64]bool, len(followeeIDs))
	if followerID == 0 || len(followeeIDs) == 0 {
		return res, nil
	}

	var ids []uint64
	if err := r.db.WithContext(ctx).
		Model(&model.Follow{}).
		Select("followee_id").
		Where("follower_id = ? AND followee_id IN ?", followerID, followeeIDs).
		Find(&ids).Error; err != nil {
		return nil, err
	}
	for _, id := range ids {
		res[id] = true
	}
	return res, nil
}

func (r *followRepo) ListFollowingIDs(ctx context.Context, userID uint64, page int, size int) ([]uint64, error) {
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
		Model(&model.Follow{}).
		Select("followee_id").
		Where("follower_id = ?", userID).
		Order("id DESC").
		Offset(offset).
		Limit(size).
		Find(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *followRepo) ListFollowerIDs(ctx context.Context, userID uint64, page int, size int) ([]uint64, error) {
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
		Model(&model.Follow{}).
		Select("follower_id").
		Where("followee_id = ?", userID).
		Order("id DESC").
		Offset(offset).
		Limit(size).
		Find(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
