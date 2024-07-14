// internal/storage/storage.go

package storage

import (
	"context"

	pb "momentum-trading-platform/api/proto/portfolio_state_service"
)

type Storage interface {
	SavePortfolioState(ctx context.Context, state *pb.PortfolioState) error
	LoadPortfolioState(ctx context.Context, date string) (*pb.PortfolioState, error)
	GetPortfolioHistory(ctx context.Context, startDate, endDate string) ([]*pb.PortfolioState, error)
}
