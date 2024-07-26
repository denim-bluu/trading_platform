package main

import (
	"net"

	"github.com/charmbracelet/log"

	"momentum-trading-platform/internal/portfolio"

	pb "momentum-trading-platform/api/proto/portfolio_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	clients, err := portfolio.NewClients()
	if err != nil {
		log.Fatalf("Failed to create gRPC clients: %v", err)
	}

	s, err := portfolio.NewServer(clients)
	if err != nil {
		s.Logger.WithError(err).Fatal("Failed to create server")
	}

	lis, err := net.Listen("tcp", "0.0.0.0:50054")
	if err != nil {
		s.Logger.WithError(err).Fatal("Failed to listen")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPortfolioServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	s.Logger.WithField("address", lis.Addr().String()).Info("Portfolio service starting")
	if err := grpcServer.Serve(lis); err != nil {
		s.Logger.WithError(err).Fatal("Failed to serve")
	}
}
