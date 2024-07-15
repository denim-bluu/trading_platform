package main

import (
	"log"
	"net"

	"momentum-trading-platform/internal/data"

	pb "momentum-trading-platform/api/proto/data_service"

	"google.golang.org/grpc"
)

func main() {
	clients, err := data.NewClients()
	if err != nil {
		log.Fatalf("Failed to create gRPC clients: %v", err)
	}
	defer clients.Close()

	s := data.NewServer()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		s.Logger.WithError(err).Fatal("Failed to listen")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataServiceServer(grpcServer, s)

	s.Logger.WithField("address", lis.Addr().String()).Info("Portfolio service starting")
	if err := grpcServer.Serve(lis); err != nil {
		s.Logger.WithError(err).Fatal("Failed to serve")
	}
}
