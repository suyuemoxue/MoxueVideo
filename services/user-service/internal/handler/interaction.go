package handler

import (
	"errors"
	"net/http"

	"example.com/MoxueVideo/user-service/internal/repo"
	"example.com/MoxueVideo/user-service/internal/service"
	"example.com/MoxueVideo/user-service/internal/transport/httpx"
	"github.com/gin-gonic/gin"
)

type InteractionHandler struct {
	interact *service.InteractionService
	videoSvc *service.VideoService
	likes    repo.LikeRepo
	fav      repo.FavoriteRepo
}

func NewInteractionHandler(interact *service.InteractionService, videoSvc *service.VideoService, likes repo.LikeRepo, fav repo.FavoriteRepo) *InteractionHandler {
	return &InteractionHandler{interact: interact, videoSvc: videoSvc, likes: likes, fav: fav}
}

type actionRequest struct {
	VideoID uint64 `json:"videoId" binding:"required"`
	Action  int    `json:"action" binding:"required"`
}

func (h *InteractionHandler) LikeAction(c *gin.Context) {
	userID, ok := httpx.UserID(c)
	if !ok {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "unauthorized")
		return
	}

	var req actionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "invalid request")
		return
	}

	liked := req.Action == 1
	if req.Action != 0 && req.Action != 1 {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "invalid action")
		return
	}

	if err := h.interact.SetLike(c.Request.Context(), userID, req.VideoID, liked); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httpx.Fail(c, http.StatusNotFound, httpx.CodeNotFound, "video not found")
			return
		}
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	httpx.OK(c, http.StatusOK, gin.H{"ok": true})
}

func (h *InteractionHandler) FavoriteAction(c *gin.Context) {
	userID, ok := httpx.UserID(c)
	if !ok {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "unauthorized")
		return
	}

	var req actionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "invalid request")
		return
	}

	favored := req.Action == 1
	if req.Action != 0 && req.Action != 1 {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "invalid action")
		return
	}

	if err := h.interact.SetFavorite(c.Request.Context(), userID, req.VideoID, favored); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httpx.Fail(c, http.StatusNotFound, httpx.CodeNotFound, "video not found")
			return
		}
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	httpx.OK(c, http.StatusOK, gin.H{"ok": true})
}

func (h *InteractionHandler) LikedList(c *gin.Context) {
	viewerID, ok := httpx.UserID(c)
	if !ok {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "unauthorized")
		return
	}
	userID := httpx.QueryUint64(c, "userId", viewerID)
	page := httpx.QueryInt(c, "page", 1, 1, 1000000)
	size := httpx.QueryInt(c, "size", 20, 1, 50)

	videoIDs, err := h.likes.ListVideoIDsByUser(c.Request.Context(), userID, page, size)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}
	if len(videoIDs) == 0 {
		httpx.OK(c, http.StatusOK, gin.H{"items": []service.VideoDTO{}})
		return
	}

	videos := make([]service.VideoDTO, 0, len(videoIDs))
	for _, id := range videoIDs {
		dto, err := h.videoSvc.Get(c.Request.Context(), viewerID, id)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				continue
			}
			continue
		}
		videos = append(videos, *dto)
	}

	httpx.OK(c, http.StatusOK, gin.H{"items": videos})
}

func (h *InteractionHandler) FavoriteList(c *gin.Context) {
	viewerID, ok := httpx.UserID(c)
	if !ok {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "unauthorized")
		return
	}
	userID := httpx.QueryUint64(c, "userId", viewerID)
	page := httpx.QueryInt(c, "page", 1, 1, 1000000)
	size := httpx.QueryInt(c, "size", 20, 1, 50)

	videoIDs, err := h.fav.ListVideoIDsByUser(c.Request.Context(), userID, page, size)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}
	if len(videoIDs) == 0 {
		httpx.OK(c, http.StatusOK, gin.H{"items": []service.VideoDTO{}})
		return
	}

	videos := make([]service.VideoDTO, 0, len(videoIDs))
	for _, id := range videoIDs {
		dto, err := h.videoSvc.Get(c.Request.Context(), viewerID, id)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				continue
			}
			continue
		}
		videos = append(videos, *dto)
	}

	httpx.OK(c, http.StatusOK, gin.H{"items": videos})
}
