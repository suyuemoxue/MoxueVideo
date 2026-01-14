package domain

type ChatMessageCreated struct {
	MessageID  uint64 `json:"message_id"`
	SenderID   uint64 `json:"sender_id"`
	ReceiverID uint64 `json:"receiver_id"`
	MsgType    string `json:"msg_type"`
	Content    string `json:"content"`
	Uniqued    string `json:"uniqued"`
	CreatedAt  int64  `json:"created_at"`
}
