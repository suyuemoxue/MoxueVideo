package model

import "time"

type Chat struct {
	ID         uint64     `gorm:"primaryKey;column:id"`
	FromUserID uint64     `gorm:"column:from_user_id;not null;index:idx_chat_from_to_ct,priority:1"`
	ToUserID   uint64     `gorm:"column:to_user_id;not null;index:idx_chat_from_to_ct,priority:2;index:idx_chat_to_isread_ct,priority:1;index:idx_chat_to_from_ct,priority:1"`
	MsgType    string     `gorm:"column:msg_type;type:enum('text','picture','audio');not null;default:'text'"`
	Content    string     `gorm:"column:content;type:text;not null"`
	IsRead     uint8      `gorm:"column:is_read;not null;default:0;index:idx_chat_to_isread_ct,priority:2"`
	Uniqued    string     `gorm:"column:uniqued;size:64;not null;uniqueIndex:uk_chat_uniqued"`
	CreateTime time.Time  `gorm:"column:create_time;autoCreateTime:milli"`
	ReadTime   *time.Time `gorm:"column:read_time"`
	IsDel      uint8      `gorm:"column:is_del;not null;default:0"`
}

func (Chat) TableName() string {
	return "chat"
}
