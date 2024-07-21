package main

import (
	"log"
	"net"

	"momentum-trading-platform/internal/strategy"

	pb "momentum-trading-platform/api/proto/strategy_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	clients, err := strategy.NewClients()
	if err != nil {
		log.Fatalf("Failed to create gRPC clients: %v", err)
	}
	defer clients.Close()

	s := strategy.NewServer(clients)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		s.Logger.WithError(err).Fatal("Failed to listen")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStrategyServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	s.Logger.WithField("address", lis.Addr().String()).Info("Strategy service starting")
	if err := grpcServer.Serve(lis); err != nil {
		s.Logger.WithError(err).Fatal("Failed to serve")
	}
}
