// cmd/portfolio_state/main.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net"
	"os"

	log "github.com/sirupsen/logrus"

	pb "momentum-trading-platform/api/proto/portfolio_state_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedPortfolioStateServiceServer
	db     *sql.DB
	logger *log.Logger
}

func newServer() *server {
	db, err := sql.Open("postgres", "postgresql://username:password@localhost/portfolio_db?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)

	return &server{logger: logger, db: db}
}

func (s *server) SavePortfolioState(ctx context.Context, state *pb.PortfolioState) (*pb.SaveResponse, error) {
	positionsJSON, err := json.Marshal(state.Positions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to marshal positions: %v", err)
	}
	s.logger.Infof("Saving portfolio state for date: %s", state.Date)
	_, err = s.db.ExecContext(ctx,
		"INSERT INTO portfolio_states (date, positions, cash_balance, total_value) VALUES ($1, $2, $3, $4)",
		state.Date, positionsJSON, state.CashBalance, state.TotalValue)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save portfolio state: %v", err)
	}

	return &pb.SaveResponse{Success: true, Message: "Portfolio state saved successfully"}, nil
}

func (s *server) LoadPortfolioState(ctx context.Context, req *pb.LoadRequest) (*pb.PortfolioState, error) {
	var state pb.PortfolioState
	var positionsJSON []byte
	s.logger.Infof("Loading portfolio state for date: %s", req.Date)

	err := s.db.QueryRowContext(ctx,
		"SELECT date, positions, cash_balance, total_value FROM portfolio_states WHERE date = $1",
		req.Date).Scan(&state.Date, &positionsJSON, &state.CashBalance, &state.TotalValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "No portfolio state found for date: %s", req.Date)
		}
		return nil, status.Errorf(codes.Internal, "Failed to load portfolio state: %v", err)
	}

	err = json.Unmarshal(positionsJSON, &state.Positions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to unmarshal positions: %v", err)
	}

	return &state, nil
}

func (s *server) GetPortfolioHistory(ctx context.Context, req *pb.HistoryRequest) (*pb.PortfolioHistory, error) {
	s.logger.Infof("Getting portfolio history between %s and %s", req.StartDate, req.EndDate)
	rows, err := s.db.QueryContext(ctx,
		"SELECT date, positions, cash_balance, total_value FROM portfolio_states WHERE date BETWEEN $1 AND $2 ORDER BY date",
		req.StartDate, req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query portfolio history: %v", err)
	}
	defer rows.Close()

	var history pb.PortfolioHistory
	for rows.Next() {
		var state pb.PortfolioState
		var positionsJSON []byte
		err := rows.Scan(&state.Date, &positionsJSON, &state.CashBalance, &state.TotalValue)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan row: %v", err)
		}

		err = json.Unmarshal(positionsJSON, &state.Positions)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to unmarshal positions: %v", err)
		}

		history.States = append(history.States, &state)
	}

	return &history, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPortfolioStateServiceServer(s, newServer())
	log.Printf("Portfolio State service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
