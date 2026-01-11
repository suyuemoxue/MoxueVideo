package model

import "time"

type Video struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	AuthorID    uint64    `gorm:"index;not null" json:"authorId"`
	PlayURL     string    `gorm:"size:1024;not null" json:"playUrl"`
	CoverURL    string    `gorm:"size:1024" json:"coverUrl"`
	Title       string    `gorm:"size:128;not null" json:"title"`
	Description string    `gorm:"size:1024" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
