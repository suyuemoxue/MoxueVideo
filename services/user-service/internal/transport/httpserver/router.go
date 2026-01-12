package httpserver

import (
	"example.com/MoxueVideo/user-service/internal/config"
	"example.com/MoxueVideo/user-service/internal/handler"
	"example.com/MoxueVideo/user-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, auth *handler.AuthHandler, me *handler.MeHandler) *gin.Engine {
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

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWT.Secret))
		protected.GET("/me", me.Me)
	}

	return r
}
