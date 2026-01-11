package model

import "time"

type Favorite struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null;uniqueIndex:uniq_fav" json:"userId"`
	VideoID   uint64    `gorm:"not null;uniqueIndex:uniq_fav" json:"videoId"`
	CreatedAt time.Time `json:"createdAt"`
}
