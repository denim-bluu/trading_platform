package portfolio

import (
	"context"
	"fmt"
	"time"

	pb "momentum-trading-platform/api/proto/portfolio_service"
	statepb "momentum-trading-platform/api/proto/portfolio_state_service"
)

func (s *Server) loadLastState() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	lastState, err := s.Clients.StateClient.LoadPortfolioState(ctx, &statepb.LoadRequest{
		Date: time.Now().Format("2006-01-02"), // Or the last known date
	})
	if err != nil {
		return fmt.Errorf("failed to load last portfolio state: %v", err)
	}
	s.Logger.Infof("Loaded last portfolio state: %s", lastState.Date)

	s.Portfolio = s.convertStatePositions(lastState.Positions)
	s.CashBalance = lastState.CashBalance
	s.LastUpdateDate = lastState.Date

	return nil
}

func (s *Server) saveState() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	positions := make([]*statepb.Position, 0, len(s.Portfolio))
	for _, pos := range s.Portfolio {
		positions = append(positions, &statepb.Position{
			Symbol:       pos.Symbol,
			Quantity:     pos.Quantity,
			AveragePrice: pos.AveragePrice,
			CurrentPrice: pos.CurrentPrice,
			MarketValue:  pos.MarketValue,
		})
	}

	_, err := s.Clients.StateClient.SavePortfolioState(ctx, &statepb.PortfolioState{
		Date:        time.Now().Format("2006-01-02"),
		Positions:   positions,
		CashBalance: s.CashBalance,
		TotalValue:  s.getTotalPortfolioValue(),
	})

	if err != nil {
		return fmt.Errorf("failed to save portfolio state: %v", err)
	}

	s.Logger.Info("Portfolio state saved successfully")
	return nil
}

func (s *Server) convertStatePositions(positions []*statepb.Position) map[string]*pb.Position {
	portfolio := make(map[string]*pb.Position)
	for _, pos := range positions {
		portfolio[pos.Symbol] = &pb.Position{
			Symbol:       pos.Symbol,
			Quantity:     pos.Quantity,
			AveragePrice: pos.AveragePrice,
			CurrentPrice: pos.CurrentPrice,
			MarketValue:  pos.MarketValue,
		}
	}
	return portfolio
}
