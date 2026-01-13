package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"moxuevideo/core/internal/config"
	"moxuevideo/core/internal/infra/grpcchat"
	"moxuevideo/core/internal/infra/mq"
	"moxuevideo/core/internal/infra/mysql"
	"moxuevideo/core/internal/infra/redisc"
	"moxuevideo/core/internal/transport/httpserver"
)

func main() {
	cfg := config.Load()

	db, err := mysql.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("mysql open: %v", err)
	}

	if err := mysql.Migrate(db); err != nil {
		log.Fatalf("mysql migrate: %v", err)
	}

	redisClient := redisc.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)

	rmq, err := mq.Open(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("rabbitmq open: %v", err)
	}
	defer rmq.Close()

	chatClient, err := grpcchat.Dial(cfg.ChatGRPCAddr)
	if err != nil {
		log.Fatalf("grpc chat dial: %v", err)
	}
	defer chatClient.Close()

	router := httpserver.NewRouter(httpserver.Dependencies{
		DB:         db,
		Redis:      redisClient,
		RabbitMQ:   rmq,
		ChatClient: chatClient,
	})

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("http listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
