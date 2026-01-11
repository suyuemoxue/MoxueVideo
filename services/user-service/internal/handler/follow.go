package handler

import (
	"errors"
	"net/http"

	"example.com/MoxueVideo/user-service/internal/repo"
	"example.com/MoxueVideo/user-service/internal/service"
	"example.com/MoxueVideo/user-service/internal/transport/httpx"
	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followSvc *service.FollowService
	users     repo.UserRepo
	follows   repo.FollowRepo
}

func NewFollowHandler(followSvc *service.FollowService, users repo.UserRepo, follows repo.FollowRepo) *FollowHandler {
	return &FollowHandler{followSvc: followSvc, users: users, follows: follows}
}

type followActionRequest struct {
	ToUserID uint64 `json:"toUserId" binding:"required"`
	Action   int    `json:"action" binding:"required"`
}

func (h *FollowHandler) Action(c *gin.Context) {
	userID, ok := httpx.UserID(c)
	if !ok {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "unauthorized")
		return
	}

	var req followActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "invalid request")
		return
	}
	following := req.Action == 1
	if req.Action != 0 && req.Action != 1 {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "invalid action")
		return
	}

	if err := h.followSvc.SetFollow(c.Request.Context(), userID, req.ToUserID, following); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httpx.Fail(c, http.StatusNotFound, httpx.CodeNotFound, "user not found")
			return
		}
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, err.Error())
		return
	}

	httpx.OK(c, http.StatusOK, gin.H{"ok": true})
}

func (h *FollowHandler) FollowingList(c *gin.Context) {
	viewerID, ok := httpx.UserID(c)
	if !ok {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "unauthorized")
		return
	}
	userID := httpx.QueryUint64(c, "userId", viewerID)
	page := httpx.QueryInt(c, "page", 1, 1, 1000000)
	size := httpx.QueryInt(c, "size", 20, 1, 50)

	ids, err := h.follows.ListFollowingIDs(c.Request.Context(), userID, page, size)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}
	if len(ids) == 0 {
		httpx.OK(c, http.StatusOK, gin.H{"items": []service.UserDTO{}})
		return
	}

	users, err := h.users.FindByIDs(c.Request.Context(), ids)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	followersCount, err := h.follows.CountFollowersByUserIDs(c.Request.Context(), ids)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}
	followingCount, err := h.follows.CountFollowingByUserIDs(c.Request.Context(), ids)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}
	isFollowingMap, err := h.follows.IsFollowingMap(c.Request.Context(), viewerID, ids)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	out := make([]service.UserDTO, 0, len(users))
	for _, u := range users {
		out = append(out, service.UserDTO{
			ID:             u.ID,
			Username:       u.Username,
			DisplayName:    u.DisplayName,
			AvatarURL:      u.AvatarURL,
			FollowerCount:  followersCount[u.ID],
			FollowingCount: followingCount[u.ID],
			IsFollowing:    isFollowingMap[u.ID],
		})
	}

	httpx.OK(c, http.StatusOK, gin.H{"items": out})
}

func (h *FollowHandler) FollowersList(c *gin.Context) {
	viewerID, ok := httpx.UserID(c)
	if !ok {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "unauthorized")
		return
	}
	userID := httpx.QueryUint64(c, "userId", viewerID)
	page := httpx.QueryInt(c, "page", 1, 1, 1000000)
	size := httpx.QueryInt(c, "size", 20, 1, 50)

	ids, err := h.follows.ListFollowerIDs(c.Request.Context(), userID, page, size)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}
	if len(ids) == 0 {
		httpx.OK(c, http.StatusOK, gin.H{"items": []service.UserDTO{}})
		return
	}

	users, err := h.users.FindByIDs(c.Request.Context(), ids)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	followersCount, err := h.follows.CountFollowersByUserIDs(c.Request.Context(), ids)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}
	followingCount, err := h.follows.CountFollowingByUserIDs(c.Request.Context(), ids)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}
	isFollowingMap, err := h.follows.IsFollowingMap(c.Request.Context(), viewerID, ids)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	out := make([]service.UserDTO, 0, len(users))
	for _, u := range users {
		out = append(out, service.UserDTO{
			ID:             u.ID,
			Username:       u.Username,
			DisplayName:    u.DisplayName,
			AvatarURL:      u.AvatarURL,
			FollowerCount:  followersCount[u.ID],
			FollowingCount: followingCount[u.ID],
			IsFollowing:    isFollowingMap[u.ID],
		})
	}

	httpx.OK(c, http.StatusOK, gin.H{"items": out})
}
