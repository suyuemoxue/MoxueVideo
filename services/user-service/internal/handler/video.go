package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"example.com/MoxueVideo/user-service/internal/service"
	"example.com/MoxueVideo/user-service/internal/transport/httpx"
	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videos *service.VideoService
	secret string
}

func NewVideoHandler(videos *service.VideoService, jwtSecret string) *VideoHandler {
	return &VideoHandler{videos: videos, secret: jwtSecret}
}

type publishRequest struct {
	PlayURL     string `json:"playUrl" binding:"required"`
	CoverURL    string `json:"coverUrl"`
	Title       string `json:"title" binding:"required,max=128"`
	Description string `json:"description"`
}

func (h *VideoHandler) Publish(c *gin.Context) {
	userID, ok := httpx.UserID(c)
	if !ok {
		httpx.Fail(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "unauthorized")
		return
	}

	var req publishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "invalid request")
		return
	}

	v, err := h.videos.Publish(c.Request.Context(), userID, service.PublishInput{
		PlayURL:     req.PlayURL,
		CoverURL:    req.CoverURL,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	dto, err := h.videos.Get(c.Request.Context(), userID, v.ID)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	httpx.OK(c, http.StatusCreated, dto)
}

func (h *VideoHandler) Feed(c *gin.Context) {
	viewerID := h.optionalViewerID(c)
	cursor := httpx.QueryUint64(c, "cursor", 0)
	limit := httpx.QueryInt(c, "limit", 20, 1, 50)

	res, err := h.videos.Feed(c.Request.Context(), viewerID, cursor, limit)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	var nextCursor uint64
	if len(res) > 0 {
		nextCursor = res[len(res)-1].ID
	}

	httpx.OK(c, http.StatusOK, gin.H{"items": res, "nextCursor": nextCursor})
}

func (h *VideoHandler) Get(c *gin.Context) {
	viewerID := h.optionalViewerID(c)
	id := parseUint64Param(c, "id")
	if id == 0 {
		id = httpx.QueryUint64(c, "id", 0)
	}
	if id == 0 {
		id = httpx.QueryUint64(c, "videoId", 0)
	}
	if id == 0 {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "missing id")
		return
	}

	res, err := h.videos.Get(c.Request.Context(), viewerID, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httpx.Fail(c, http.StatusNotFound, httpx.CodeNotFound, "not found")
			return
		}
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}

	httpx.OK(c, http.StatusOK, res)
}

func (h *VideoHandler) ListByUser(c *gin.Context) {
	viewerID := h.optionalViewerID(c)
	authorID := parseUint64Param(c, "id")
	if authorID == 0 {
		authorID = httpx.QueryUint64(c, "userId", 0)
	}
	if authorID == 0 {
		httpx.Fail(c, http.StatusBadRequest, httpx.CodeBadRequest, "missing userId")
		return
	}

	page := httpx.QueryInt(c, "page", 1, 1, 1000000)
	size := httpx.QueryInt(c, "size", 20, 1, 50)

	res, err := h.videos.ListByAuthor(c.Request.Context(), viewerID, authorID, page, size)
	if err != nil {
		httpx.Fail(c, http.StatusInternalServerError, httpx.CodeInternal, "internal error")
		return
	}
	if len(res) == 0 {
		httpx.OK(c, http.StatusOK, gin.H{"items": []service.VideoDTO{}})
		return
	}

	httpx.OK(c, http.StatusOK, gin.H{"items": res})
}

func (h *VideoHandler) optionalViewerID(c *gin.Context) uint64 {
	hdr := strings.TrimSpace(c.GetHeader("Authorization"))
	if hdr == "" {
		return 0
	}
	token := strings.TrimSpace(strings.TrimPrefix(hdr, "Bearer"))
	if token == "" {
		return 0
	}
	id, err := service.ParseUserIDFromToken(token, h.secret)
	if err != nil {
		return 0
	}
	return id
}

func parseUint64Param(c *gin.Context, key string) uint64 {
	v := strings.TrimSpace(c.Param(key))
	if v == "" {
		return 0
	}
	u, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0
	}
	return uint64(u)
}
