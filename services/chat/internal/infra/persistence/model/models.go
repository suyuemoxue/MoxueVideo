package model

import "time"

type DMThread struct {
	ID            uint64     `gorm:"primaryKey"`
	UserLowID     uint64     `gorm:"index:uk_dm_threads_pair,unique;not null"`
	UserHighID    uint64     `gorm:"index:uk_dm_threads_pair,unique;not null"`
	LastMessageID *uint64    `gorm:"index"`
	LastMessageAt *time.Time `gorm:"index"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
}

type DMMessage struct {
	ID         uint64    `gorm:"primaryKey"`
	ThreadID   uint64    `gorm:"index:idx_dm_messages_thread_id_id,priority:1;not null"`
	SenderID   uint64    `gorm:"index:idx_dm_messages_sender_id_created_at,priority:1;not null"`
	ReceiverID uint64    `gorm:"index:idx_dm_messages_receiver_id_id,priority:1;not null"`
	Content    string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type DMMessageRead struct {
	ID        uint64    `gorm:"primaryKey"`
	MessageID uint64    `gorm:"index:uk_dm_message_reads_message_user,unique;not null"`
	UserID    uint64    `gorm:"index:uk_dm_message_reads_message_user,unique;not null"`
	ReadAt    time.Time `gorm:"autoCreateTime"`
}
