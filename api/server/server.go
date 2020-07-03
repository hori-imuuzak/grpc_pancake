package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"pancake/maker/api/gen/api"
	"pancake/maker/api/handler"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = 50051
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	grpc_zap.ReplaceGrpcLogger(zapLogger)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_zap.UnaryServerInterceptor(zapLogger)),
	)
	api.RegisterPancakeBakerServiceServer(server, handler.NewBakerHandler())
	reflection.Register(server)

	go func() {
		log.Printf("start gRPC server port :%v", port)
		server.Serve(lis)
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	server.GracefulStop()
}
