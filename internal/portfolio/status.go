// internal/portfolio/status.go
package portfolio

import (
	"context"
	pb "momentum-trading-platform/api/proto/portfolio_service"
)

func (s *Server) GetDesiredPortfolioState(ctx context.Context, req *pb.GetDesiredPortfolioStateRequest) (*pb.PortfolioState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	positions := make([]*pb.Position, 0, len(s.DesiredPortfolio))
	for _, pos := range s.DesiredPortfolio {
		positions = append(positions, pos)
	}

	return &pb.PortfolioState{
		Positions:   positions,
		CashBalance: s.CashBalance,
		TotalValue:  s.calculateTotalValue(),
	}, nil
}
