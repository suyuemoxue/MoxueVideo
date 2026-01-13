package model

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey"`
	Username  string    `gorm:"size:64;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Follow struct {
	ID         uint64    `gorm:"primaryKey"`
	FollowerID uint64    `gorm:"index;not null"`
	FolloweeID uint64    `gorm:"index;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type Video struct {
	ID          uint64    `gorm:"primaryKey"`
	AuthorID    uint64    `gorm:"index;not null"`
	Title       string    `gorm:"size:200;not null"`
	CoverURL    string    `gorm:"size:512"`
	PlayURL     string    `gorm:"size:512;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	PublishedAt time.Time `gorm:"index"`
}

type Like struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `gorm:"index;not null"`
	VideoID   uint64    `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Favorite struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `gorm:"index;not null"`
	VideoID   uint64    `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type WatchHistory struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `gorm:"index;not null"`
	VideoID   uint64    `gorm:"index;not null"`
	WatchedAt time.Time `gorm:"index;not null"`
}
