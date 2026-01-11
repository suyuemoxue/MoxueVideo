package httpserver

import (
	"example.com/MoxueVideo/user-service/internal/config"
	"example.com/MoxueVideo/user-service/internal/handler"
	"example.com/MoxueVideo/user-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg config.Config,
	auth *handler.AuthHandler,
	me *handler.MeHandler,
	videos *handler.VideoHandler,
	interactions *handler.InteractionHandler,
	follow *handler.FollowHandler,
) *gin.Engine {
	if cfg.Env != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	health := handler.NewHealthHandler()
	r.GET("/healthz", health.Healthz)

	v1 := r.Group("/api/v1")
	{
		authg := v1.Group("/auth")
		authg.POST("/register", auth.Register)
		authg.POST("/login", auth.Login)

		v1.GET("/videos/feed", videos.Feed)
		v1.GET("/videos/:id", videos.Get)
		v1.GET("/users/:id/videos", videos.ListByUser)

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWT.Secret))
		protected.GET("/me", me.Me)
		protected.POST("/videos/publish", videos.Publish)

		protected.POST("/likes/action", interactions.LikeAction)
		protected.GET("/likes/list", interactions.LikedList)

		protected.POST("/favorites/action", interactions.FavoriteAction)
		protected.GET("/favorites/list", interactions.FavoriteList)

		protected.POST("/follow/action", follow.Action)
		protected.GET("/follow/following", follow.FollowingList)
		protected.GET("/follow/followers", follow.FollowersList)
	}

	return r
}
