package main

import (
	"context"
	"log"
	"net"

	pb "momentum-trading-platform/api/proto/data_service"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDataServiceServer
}

func (s *server) GetStockData(ctx context.Context, req *pb.StockRequest) (*pb.StockResponse, error) {
	// Mock implementation
	return &pb.StockResponse{
		Symbol: req.Symbol,
		Date:   req.Date,
		Price:  100.0,   // Mock price
		Volume: 1000000, // Mock volume
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDataServiceServer(s, &server{})
	log.Printf("Data service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
