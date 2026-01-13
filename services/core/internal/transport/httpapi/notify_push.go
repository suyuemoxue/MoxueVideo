package httpapi

func (h *Handler) PushChatMessageCreated(userID uint64, evt any) int {
	return h.notifyHub.broadcastToUser(userID, newWSMessage("notify", evt))
}
