package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/MoxueVideo/user-service/internal/config"
	"example.com/MoxueVideo/user-service/internal/handler"
	"example.com/MoxueVideo/user-service/internal/infra/mysql"
	"example.com/MoxueVideo/user-service/internal/logger"
	"example.com/MoxueVideo/user-service/internal/model"
	"example.com/MoxueVideo/user-service/internal/repo"
	"example.com/MoxueVideo/user-service/internal/service"
	"example.com/MoxueVideo/user-service/internal/transport/httpserver"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Env)

	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "dev-secret"
		log.Warn("JWT_SECRET is empty; using dev-secret")
	}

	db, err := mysql.Open(cfg.MySQL.DSN())
	if err != nil {
		log.Error("open mysql failed", slog.Any("err", err))
		os.Exit(1)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.Follow{},
		&model.Like{},
		&model.Favorite{},
	); err != nil {
		log.Error("auto migrate failed", slog.Any("err", err))
		os.Exit(1)
	}

	userRepo := repo.NewUserRepo(db)
	videoRepo := repo.NewVideoRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	favoriteRepo := repo.NewFavoriteRepo(db)
	followRepo := repo.NewFollowRepo(db)

	authSvc := service.NewAuthService(userRepo, cfg.JWT.Secret)
	videoSvc := service.NewVideoService(videoRepo, userRepo, likeRepo, favoriteRepo, followRepo)
	interactSvc := service.NewInteractionService(videoRepo, likeRepo, favoriteRepo)
	followSvc := service.NewFollowService(userRepo, followRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	meHandler := handler.NewMeHandler(userRepo)
	videoHandler := handler.NewVideoHandler(videoSvc, cfg.JWT.Secret)
	interactionHandler := handler.NewInteractionHandler(interactSvc, videoSvc, likeRepo, favoriteRepo)
	followHandler := handler.NewFollowHandler(followSvc, userRepo, followRepo)

	r := httpserver.NewRouter(cfg, authHandler, meHandler, videoHandler, interactionHandler, followHandler)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("http server listening", slog.String("addr", cfg.HTTPAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server failed", slog.Any("err", err))
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown failed", slog.Any("err", err))
	}
}
