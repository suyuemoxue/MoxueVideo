package httpserver

import (
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(deps Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	h := NewHandlers(deps)

	r.GET("/healthz", h.Healthz)

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/register", h.Register)
			users.POST("/login", h.Login)
			users.GET("/:id", h.GetUser)
			users.GET("/:id/following", h.ListFollowing)
			users.GET("/:id/followers", h.ListFollowers)
			users.POST("/:id/follow", h.FollowUser)
			users.POST("/:id/unfollow", h.UnfollowUser)
		}

		videos := v1.Group("/videos")
		{
			videos.POST("/upload", h.UploadVideo)
			videos.GET("/:id", h.GetVideo)
			videos.POST("/:id/like", h.LikeVideo)
			videos.POST("/:id/unlike", h.UnlikeVideo)
			videos.POST("/:id/favorite", h.FavoriteVideo)
			videos.POST("/:id/unfavorite", h.UnfavoriteVideo)
			videos.POST("/:id/watch", h.RecordWatch)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "NOT_FOUND", "message": "route not found", "ts": time.Now().Unix()})
	})

	return r
}
