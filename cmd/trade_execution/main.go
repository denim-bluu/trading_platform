package main

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "momentum-trading-platform/api/proto/trade_execution_service"
	tradeexecution "momentum-trading-platform/internal/trade_execution"
)

func main() {
	clients, err := tradeexecution.NewClients()
	if err != nil {
		log.Fatalf("Failed to create clients: %v", err)
	}
	defer clients.Close()

	server := tradeexecution.NewServer(clients)

	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTradeExecutionServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Printf("Trade Execution Service listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
