package storage

import (
	"context"

	pb "momentum-trading-platform/api/proto/portfolio_service"
)

type Storage interface {
	SavePortfolioState(ctx context.Context, state *pb.PortfolioStatus) error
	LoadPortfolioState(ctx context.Context) (*pb.PortfolioStatus, error)
	SaveTrade(ctx context.Context, trade *pb.Trade) error
	GetTradeHistory(ctx context.Context, startDate, endDate string) ([]*pb.Trade, error)
}
