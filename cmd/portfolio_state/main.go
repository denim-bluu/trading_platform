// cmd/portfolio_state/main.go

package main

import (
	"context"
	"net"
	"os"

	log "github.com/sirupsen/logrus"

	pb "momentum-trading-platform/api/proto/portfolio_state_service"
	"momentum-trading-platform/internal/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedPortfolioStateServiceServer
	logger  *log.Logger
	storage storage.Storage
}

func newServer(storage storage.Storage) *server {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	return &server{logger: logger, storage: storage}
}

func (s *server) SavePortfolioState(ctx context.Context, state *pb.PortfolioState) (*pb.SaveResponse, error) {
	err := s.storage.SavePortfolioState(ctx, state)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save portfolio state: %v", err)
	}

	return &pb.SaveResponse{Success: true, Message: "Portfolio state saved successfully"}, nil
}

func (s *server) LoadPortfolioState(ctx context.Context, req *pb.LoadRequest) (*pb.PortfolioState, error) {
	state, err := s.storage.LoadPortfolioState(ctx, req.Date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to load portfolio state: %v", err)
	}

	return state, nil
}

func (s *server) GetPortfolioHistory(ctx context.Context, req *pb.HistoryRequest) (*pb.PortfolioHistory, error) {
	states, err := s.storage.GetPortfolioHistory(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get portfolio history: %v", err)
	}

	return &pb.PortfolioHistory{States: states}, nil
}

func main() {
	csvStorage := storage.NewCSVStorage("portfolio_states.csv")
	s := newServer(csvStorage)

	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPortfolioStateServiceServer(grpcServer, s)
	log.Printf("Portfolio State service listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
