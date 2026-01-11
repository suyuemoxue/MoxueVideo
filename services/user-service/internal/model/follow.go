package model

import "time"

type Follow struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID uint64    `gorm:"not null;uniqueIndex:uniq_follow" json:"followerId"`
	FolloweeID uint64    `gorm:"not null;uniqueIndex:uniq_follow" json:"followeeId"`
	CreatedAt  time.Time `json:"createdAt"`
}
