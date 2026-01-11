package model

import "time"

type Like struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null;uniqueIndex:uniq_like" json:"userId"`
	VideoID   uint64    `gorm:"not null;uniqueIndex:uniq_like" json:"videoId"`
	CreatedAt time.Time `json:"createdAt"`
}
