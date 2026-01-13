package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"moxuevideo/core/internal/config"
	"moxuevideo/core/internal/domain"
	"moxuevideo/core/internal/infra/grpcchat"
	"moxuevideo/core/internal/infra/health"
	"moxuevideo/core/internal/infra/mq"
	"moxuevideo/core/internal/infra/mysql"
	"moxuevideo/core/internal/infra/osssts"
	"moxuevideo/core/internal/infra/persistence/mysqlrepo"
	"moxuevideo/core/internal/infra/redisc"
	"moxuevideo/core/internal/transport/httpapi"
	ossuc "moxuevideo/core/internal/usecase/oss"
	useruc "moxuevideo/core/internal/usecase/user"
	videouc "moxuevideo/core/internal/usecase/video"
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

	var ossSTS *osssts.Service
	if cfg.OSS.AccessKeyID != "" && cfg.OSS.AccessKeySecret != "" && cfg.OSS.RoleARN != "" && cfg.OSS.Bucket != "" {
		ossSTS, err = osssts.New(
			cfg.OSS.AccessKeyID,
			cfg.OSS.AccessKeySecret,
			cfg.OSS.Region,
			cfg.OSS.Endpoint,
			cfg.OSS.Bucket,
			cfg.OSS.RoleARN,
			cfg.OSS.RoleSessionName,
			cfg.OSS.DurationSeconds,
		)
		if err != nil {
			log.Fatalf("oss sts init: %v", err)
		}
	}
	userRepo := mysqlrepo.NewUserRepository(db)
	videoRepo := mysqlrepo.NewVideoRepository(db)

	userService := useruc.New(userRepo)
	videoService := videouc.New(videoRepo)
	ossService := ossuc.New(ossSTS)

	h := httpapi.New(httpapi.Deps{
		User:  userService,
		Video: videoService,
		OSS:   ossService,
		Health: httpapi.HealthDeps{
			MySQL:    health.NewMySQLChecker(db),
			Redis:    health.NewRedisChecker(redisClient),
			RabbitMQ: health.NewRabbitMQChecker(rmq),
			GRPC:     health.NewGRPCChecker(chatClient),
		},
	})
	router := httpapi.NewRouter(h)

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

	go func() {
		if rmq == nil || rmq.Channel == nil {
			return
		}
		const exchange = "moxuevideo.events"
		const routingKey = "chat.message.created"
		const queue = "core.notify"

		if err := rmq.Channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
			log.Printf("mq exchange declare: %v", err)
			return
		}
		q, err := rmq.Channel.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			log.Printf("mq queue declare: %v", err)
			return
		}
		if err := rmq.Channel.QueueBind(q.Name, routingKey, exchange, false, nil); err != nil {
			log.Printf("mq queue bind: %v", err)
			return
		}
		msgs, err := rmq.Channel.Consume(q.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Printf("mq consume: %v", err)
			return
		}
		for d := range msgs {
			if len(d.Body) == 0 {
				continue
			}
			var evt domain.ChatMessageCreated
			if err := json.Unmarshal(d.Body, &evt); err != nil {
				continue
			}
			if evt.ReceiverID == 0 {
				continue
			}
			h.PushChatMessageCreated(evt.ReceiverID, map[string]any{
				"kind":       "chat_message",
				"message_id": evt.MessageID,
				"thread_id":  evt.ThreadID,
				"sender_id":  evt.SenderID,
				"content":    evt.Content,
				"created_at": evt.CreatedAt,
			})
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
