package httpapi

import (
	"time"

	"github.com/gin-gonic/gin"

	"moxuevideo/core/internal/middleware"
)

func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", h.Healthz)

	v1 := r.Group("/api/v1")
	{
		ws := v1.Group("/ws")
		{
			ws.GET("/notify", h.NotifyWS)
		}

		oss := v1.Group("/oss")
		{
			oss.GET("/sts", middleware.RequireUser(), h.GetOSSTSToken)
		}

		users := v1.Group("/users")
		{
			users.POST("/register", h.Register)
			users.POST("/login", h.Login)
			users.GET("/:id", h.GetUser)
			users.GET("/:id/following", h.ListFollowing)
			users.GET("/:id/followers", h.ListFollowers)
			users.POST("/:id/follow", middleware.RequireUser(), h.FollowUser)
			users.POST("/:id/unfollow", middleware.RequireUser(), h.UnfollowUser)
		}

		videos := v1.Group("/videos")
		{
			videos.POST("/upload", middleware.RequireUser(), h.UploadVideo)
			videos.GET("/:id", h.GetVideo)
			videos.POST("/:id/like", middleware.RequireUser(), h.LikeVideo)
			videos.POST("/:id/unlike", middleware.RequireUser(), h.UnlikeVideo)
			videos.POST("/:id/favorite", middleware.RequireUser(), h.FavoriteVideo)
			videos.POST("/:id/unfavorite", middleware.RequireUser(), h.UnfavoriteVideo)
			videos.POST("/:id/comment", middleware.RequireUser(), h.CommentVideo)
			videos.POST("/:id/watch", middleware.RequireUser(), h.RecordWatch)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "NOT_FOUND", "message": "route not found", "ts": time.Now().Unix()})
	})

	return r
}
