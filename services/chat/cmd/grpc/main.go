package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"moxuevideo/chat/internal/config"
	"moxuevideo/chat/internal/transport/grpcapi"
)

func main() {
	cfg := config.Load()

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}

	srv := grpcapi.NewServer()

	go func() {
		log.Printf("grpc listening on %s", cfg.GRPCAddr)
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcapi.GracefulStop(ctx, srv)
}
