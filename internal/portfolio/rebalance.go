// internal/portfolio/rebalance.go
package portfolio

import (
	"context"
	"fmt"
	"time"

	pb "momentum-trading-platform/api/proto/portfolio_service"
	portfoliostatepb "momentum-trading-platform/api/proto/portfolio_state_service"
	strategypb "momentum-trading-platform/api/proto/strategy_service"
)

func (s *Server) PerformRebalance(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Logger.Info("Starting portfolio rebalance")

	// Get latest signals
	signals, err := s.getLatestSignals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest signals: %w", err)
	}
	s.Logger.WithField("signalCount", len(signals)).Info("Received signals for rebalance")

	// Generate orders based on signals
	ordersResp, err := s.GenerateAndSubmitOrders(ctx, &pb.GenerateOrdersRequest{Signals: signals})
	if err != nil {
		return fmt.Errorf("failed to generate and submit orders: %w", err)
	}
	s.Logger.WithField("orderCount", len(ordersResp.Orders)).Info("Generated and submitted rebalance orders")

	s.lastRebalanceTime = time.Now()
	s.Logger.WithField("lastRebalanceTime", s.lastRebalanceTime).Info("Rebalance completed successfully")
	return nil
}

func (s *Server) getLatestSignals(ctx context.Context) ([]*strategypb.StockSignal, error) {
	s.Logger.Info("Fetching latest signals from Strategy Service")

	// Assuming we want signals for all stocks in S&P 500
	symbols := []string{"AAPL", "GOOGL", "MSFT"} // This should be expanded to include all relevant symbols

	req := &strategypb.SignalRequest{
		Symbols:   symbols,
		StartDate: time.Now().AddDate(0, 0, -7).Format("2006-01-02"), // Last 7 days
		EndDate:   time.Now().Format("2006-01-02"),
		Interval:  "1d",
	}

	resp, err := s.Clients.StrategyClient.GenerateSignals(ctx, req)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch signals from Strategy Service")
		return nil, fmt.Errorf("failed to fetch signals: %w", err)
	}

	s.Logger.WithField("signalCount", len(resp.Signals)).Info("Received signals from Strategy Service")
	return resp.Signals, nil
}

func (s *Server) generateRebalanceOrders() []*pb.Order {
	s.Logger.Info("Generating rebalance orders")
	var orders []*pb.Order

	// Get current portfolio state from Portfolio State Service
	currentPortfolio, err := s.getCurrentPortfolio(context.Background())
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get current portfolio state")
		return orders
	}

	for symbol, desiredPosition := range s.DesiredPortfolio {
		currentPosition, exists := currentPortfolio[symbol]
		if !exists {
			// New position to open
			orders = append(orders, &pb.Order{
				Symbol:   symbol,
				Type:     pb.OrderType_MARKET,
				Quantity: desiredPosition.Quantity,
			})
		} else if desiredPosition.Quantity != currentPosition.Quantity {
			// Existing position to adjust
			orderQuantity := desiredPosition.Quantity - currentPosition.Quantity
			orderType := pb.OrderType_MARKET
			if orderQuantity > 0 {
				orders = append(orders, &pb.Order{
					Symbol:   symbol,
					Type:     orderType,
					Quantity: orderQuantity,
				})
			} else {
				orders = append(orders, &pb.Order{
					Symbol:   symbol,
					Type:     orderType,
					Quantity: -orderQuantity,
				})
			}
		}
	}

	// Check for positions to close
	for symbol, currentPosition := range currentPortfolio {
		if _, exists := s.DesiredPortfolio[symbol]; !exists {
			orders = append(orders, &pb.Order{
				Symbol:   symbol,
				Type:     pb.OrderType_MARKET,
				Quantity: -currentPosition.Quantity,
			})
		}
	}

	return orders
}

func (s *Server) sendOrdersToExecution(ctx context.Context, orders []*pb.Order) error {
	// Implementation to send orders to Trade Execution Service
	// This will involve making a gRPC call to the Trade Execution Service
	return nil
}

func (s *Server) getCurrentPortfolio(ctx context.Context) (map[string]*pb.Position, error) {
	state, err := s.Clients.PortfolioStateClient.GetPortfolioState(ctx, &portfoliostatepb.GetPortfolioStateRequest{})
	if err != nil {
		return nil, err
	}

	currentPortfolio := make(map[string]*pb.Position)
	for _, pos := range state.Positions {
		currentPortfolio[pos.Symbol] = &pb.Position{
			Symbol:       pos.Symbol,
			Quantity:     pos.Quantity,
			CurrentPrice: pos.CurrentPrice,
			MarketValue:  pos.MarketValue,
		}
	}

	return currentPortfolio, nil
}

func (s *Server) TriggerRebalance(ctx context.Context, req *pb.TriggerRebalanceRequest) (*pb.TriggerRebalanceResponse, error) {
	err := s.PerformRebalance(ctx)
	if err != nil {
		return &pb.TriggerRebalanceResponse{
			Success: false,
			Message: "Failed to perform rebalance: " + err.Error(),
		}, nil
	}
	return &pb.TriggerRebalanceResponse{
		Success: true,
		Message: "Rebalance performed successfully",
	}, nil
}

func (s *Server) UpdateRebalanceSchedule(ctx context.Context, req *pb.UpdateRebalanceScheduleRequest) (*pb.UpdateRebalanceScheduleResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.RebalanceSchedule = req.Schedule
	return &pb.UpdateRebalanceScheduleResponse{
		Success: true,
		Message: "Rebalance schedule updated successfully",
	}, nil
}
