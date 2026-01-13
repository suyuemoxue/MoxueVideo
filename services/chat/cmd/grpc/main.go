package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"moxuevideo/chat/internal/config"
	"moxuevideo/chat/internal/infra/mq"
	"moxuevideo/chat/internal/infra/mysql"
	"moxuevideo/chat/internal/infra/persistence/mysqlrepo"
	"moxuevideo/chat/internal/transport/grpcapi"
	"moxuevideo/chat/internal/transport/wsapi"
	chatuc "moxuevideo/chat/internal/usecase/chat"
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

	rmq, err := mq.Open(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("rabbitmq open: %v", err)
	}
	defer rmq.Close()

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}

	srv := grpcapi.NewServer()
	repo := mysqlrepo.New(db)
	pub := mq.NewPublisher(rmq)
	chatService := chatuc.New(repo, pub)
	wsSrv := wsapi.NewServer(nil, chatService)

	go func() {
		log.Printf("grpc listening on %s", cfg.GRPCAddr)
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	go func() {
		log.Printf("ws listening on %s", cfg.WSAddr)
		if err := http.ListenAndServe(cfg.WSAddr, wsSrv.Handler()); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ws serve: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcapi.GracefulStop(ctx, srv)
}
