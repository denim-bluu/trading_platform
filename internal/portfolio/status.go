package portfolio

import (
	"context"

	pb "momentum-trading-platform/api/proto/portfolio_service"
)

func (s *Server) GetPortfolioStatus(ctx context.Context, req *pb.PortfolioStatusRequest) (*pb.PortfolioStatus, error) {
	return s.getPortfolioStatus(), nil
}

func (s *Server) getPortfolioStatus() *pb.PortfolioStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	positions := make([]*pb.Position, 0, len(s.Portfolio))
	for _, position := range s.Portfolio {
		positions = append(positions, position)
	}
	return &pb.PortfolioStatus{
		Positions:      positions,
		CashBalance:    s.CashBalance,
		TotalValue:     s.getTotalPortfolioValue(),
		LastUpdateDate: s.LastUpdateDate,
	}
}

func (s *Server) getTotalPortfolioValue() float64 {
	total := s.CashBalance
	for _, position := range s.Portfolio {
		total += position.MarketValue
	}
	return total
}
