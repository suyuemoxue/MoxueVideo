package domain

import "time"

type User struct {
	ID        uint64
	Username  string
	Password  string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Follow struct {
	ID         uint64
	FollowerID uint64
	FolloweeID uint64
	CreatedAt  time.Time
}

type Video struct {
	ID          uint64
	AuthorID    uint64
	Title       string
	CoverURL    string
	PlayURL     string
	CreatedAt   time.Time
	PublishedAt time.Time
}

type Like struct {
	ID        uint64
	UserID    uint64
	VideoID   uint64
	CreatedAt time.Time
}

type Favorite struct {
	ID        uint64
	UserID    uint64
	VideoID   uint64
	CreatedAt time.Time
}

type Comment struct {
	ID        uint64
	UserID    uint64
	VideoID   uint64
	Content   string
	CreatedAt time.Time
}

type WatchHistory struct {
	ID        uint64
	UserID    uint64
	VideoID   uint64
	WatchedAt time.Time
}
