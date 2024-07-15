package main

import (
	"log"
	"net"

	"momentum-trading-platform/internal/portfolio"
	"momentum-trading-platform/internal/storage"

	pb "momentum-trading-platform/api/proto/portfolio_service"

	"google.golang.org/grpc"
)

func main() {
	clients, err := portfolio.NewClients()
	if err != nil {
		log.Fatalf("Failed to create gRPC clients: %v", err)
	}
	defer clients.Close()

	csvStorage := storage.NewCSVStorage("portfolio.csv", "trades.csv")
	s := portfolio.NewServer(clients, csvStorage)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		s.Logger.WithError(err).Fatal("Failed to listen")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPortfolioServiceServer(grpcServer, s)

	s.Logger.WithField("address", lis.Addr().String()).Info("Portfolio service starting")
	if err := grpcServer.Serve(lis); err != nil {
		s.Logger.WithError(err).Fatal("Failed to serve")
	}
}
