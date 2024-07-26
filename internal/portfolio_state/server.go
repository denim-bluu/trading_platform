// internal/portfolio_state/server.go
package portfoliostate

import (
	"context"
	"database/sql"
	"encoding/json"

	pb "momentum-trading-platform/api/proto/portfolio_state_service"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const createTableSQL = `
	CREATE TABLE IF NOT EXISTS portfolio_state (
		id SERIAL PRIMARY KEY,
		positions JSONB,
		cash_balance FLOAT,
		total_value FLOAT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

type Server struct {
	pb.UnimplementedPortfolioStateServiceServer
	Logger *log.Logger
	DB     *sql.DB
}

func NewServer(dbConnectionString string) (*Server, error) {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Logger: logger,
		DB:     db,
	}

	if err := s.initDB(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) initDB() error {
	_, err := s.DB.Exec(createTableSQL)
	return err
}

func (s *Server) GetPortfolioState(ctx context.Context, req *pb.GetPortfolioStateRequest) (*pb.PortfolioState, error) {
	s.Logger.Info("Getting portfolio state")

	var positionsJSON []byte
	var cashBalance, totalValue float64
	err := s.DB.QueryRow("SELECT positions, cash_balance, total_value FROM portfolio_state ORDER BY timestamp DESC LIMIT 1").Scan(&positionsJSON, &cashBalance, &totalValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.PortfolioState{
				Positions:   make([]*pb.Position, 0),
				CashBalance: 1000000, // Default starting cash
				TotalValue:  1000000,
			}, nil
		}
		return nil, err
	}

	var positions []*pb.Position
	if err := json.Unmarshal(positionsJSON, &positions); err != nil {
		return nil, err
	}

	return &pb.PortfolioState{
		Positions:   positions,
		CashBalance: cashBalance,
		TotalValue:  totalValue,
	}, nil
}

func (s *Server) UpdatePortfolioState(ctx context.Context, req *pb.UpdatePortfolioStateRequest) (*pb.UpdatePortfolioStateResponse, error) {
	s.Logger.Info("Updating portfolio state")

	positionsJSON, err := json.Marshal(req.Positions)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to marshal positions")
		return nil, status.Errorf(codes.Internal, "failed to process positions: %v", err)
	}

	totalValue := req.CashBalance
	for _, pos := range req.Positions {
		totalValue += pos.MarketValue
	}

	s.Logger.WithFields(log.Fields{
		"positionCount": len(req.Positions),
		"cashBalance":   req.CashBalance,
		"totalValue":    totalValue,
	}).Info("Updating portfolio state")

	_, err = s.DB.Exec("INSERT INTO portfolio_state (positions, cash_balance, total_value) VALUES ($1, $2, $3)",
		positionsJSON, req.CashBalance, totalValue)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				s.Logger.WithError(err).Error("Duplicate portfolio state entry")
				return nil, status.Errorf(codes.AlreadyExists, "portfolio state already exists for this timestamp")
			case "23502": // not_null_violation
				s.Logger.WithError(err).Error("Null value in portfolio state update")
				return nil, status.Errorf(codes.InvalidArgument, "all required fields must be provided")
			default:
				s.Logger.WithError(err).Error("Unexpected database error")
				return nil, status.Errorf(codes.Internal, "unexpected database error: %v", err)
			}
		}
		s.Logger.WithError(err).Error("Failed to insert portfolio state")
		return nil, status.Errorf(codes.Internal, "failed to update portfolio state: %v", err)
	}

	s.Logger.Info("Portfolio state updated successfully")
	return &pb.UpdatePortfolioStateResponse{
		Success: true,
		Message: "Portfolio state updated successfully",
	}, nil
}
