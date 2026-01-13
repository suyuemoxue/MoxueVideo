package model

import "time"

type User struct {
	ID           uint64    `gorm:"primaryKey;column:id"`
	Username     string    `gorm:"column:username;size:64;not null"`
	PasswordHash string    `gorm:"column:password_hash;size:255;not null"`
	AvatarURL    string    `gorm:"column:avatar_url;size:512"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type Follow struct {
	FollowerID uint64    `gorm:"primaryKey;column:follower_id"`
	FolloweeID uint64    `gorm:"primaryKey;column:followee_id"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

type Video struct {
	ID          uint64     `gorm:"primaryKey;column:id"`
	AuthorID    uint64     `gorm:"column:author_id;not null"`
	Title       string     `gorm:"column:title;size:200;not null"`
	CoverURL    string     `gorm:"column:cover_url;size:512"`
	PlayURL     string     `gorm:"column:play_url;size:512;not null"`
	PublishedAt *time.Time `gorm:"column:published_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

type Like struct {
	UserID    uint64    `gorm:"primaryKey;column:user_id"`
	VideoID   uint64    `gorm:"primaryKey;column:video_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Like) TableName() string {
	return "video_likes"
}

type Favorite struct {
	UserID    uint64    `gorm:"primaryKey;column:user_id"`
	VideoID   uint64    `gorm:"primaryKey;column:video_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Favorite) TableName() string {
	return "video_favorites"
}

type Comment struct {
	ID        uint64    `gorm:"primaryKey;column:id"`
	UserID    uint64    `gorm:"column:user_id;not null"`
	VideoID   uint64    `gorm:"column:video_id;not null"`
	Content   string    `gorm:"column:content;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type WatchHistory struct {
	ID        uint64    `gorm:"primaryKey;column:id"`
	UserID    uint64    `gorm:"column:user_id;not null"`
	VideoID   uint64    `gorm:"column:video_id;not null"`
	WatchedAt time.Time `gorm:"column:watched_at;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}
