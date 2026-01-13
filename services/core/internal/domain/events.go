package domain

type ChatMessageCreated struct {
	MessageID  uint64 `json:"message_id"`
	ThreadID   uint64 `json:"thread_id"`
	SenderID   uint64 `json:"sender_id"`
	ReceiverID uint64 `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
}
