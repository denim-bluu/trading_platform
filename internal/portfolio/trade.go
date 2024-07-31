// internal/portfolio/trade.go
package portfolio

import (
	"context"
	pb "momentum-trading-platform/api/proto/portfolio_service"
	portfoliostatepb "momentum-trading-platform/api/proto/portfolio_state_service"
	strategypb "momentum-trading-platform/api/proto/strategy_service"
	tradepb "momentum-trading-platform/api/proto/trade_execution_service"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GenerateAndSubmitOrders(ctx context.Context, req *pb.GenerateOrdersRequest) (*pb.GenerateOrdersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Logger.Info("Generating and submitting orders based on signals")

	// Get current portfolio state
	currentState, err := s.getCurrentPortfolioState(ctx)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to get current portfolio state")
		return nil, status.Errorf(codes.Internal, "failed to get current portfolio state: %v", err)
	}

	// Calculate desired portfolio
	desiredPortfolio, err := s.calculateDesiredPortfolio(req.Signals)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to calculate desired portfolio")
		return nil, status.Errorf(codes.Internal, "failed to calculate desired portfolio: %v", err)
	}

	// Generate orders
	orders := s.generateOrders(currentState, desiredPortfolio)

	// Update desired portfolio state (this is internal to the Portfolio Service)
	s.DesiredPortfolio = desiredPortfolio

	// Submit orders to Trade Execution Service
	executionResp, err := s.Clients.TradeExecutionClient.ExecuteTrades(ctx, &tradepb.ExecuteTradesRequest{
		Orders: orders,
	})
	if err != nil {
		s.Logger.WithError(err).Error("Failed to submit orders to Trade Execution Service")
		return nil, status.Errorf(codes.Internal, "failed to submit orders: %v", err)
	}

	s.Logger.WithFields(log.Fields{
		"orderCount":  len(orders),
		"executionId": executionResp.ExecutionId,
	}).Info("Orders generated and submitted successfully")

	return &pb.GenerateOrdersResponse{
		Orders: orders,
	}, nil
}

// Keep the generateOrders, calculateDesiredPortfolio, and getCurrentPortfolioState methods as they are

func (s *Server) getCurrentPortfolioState(ctx context.Context) (map[string]*pb.Position, error) {
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

func (s *Server) calculateDesiredPortfolio(signals []*strategypb.StockSignal) (map[string]*pb.Position, error) {
	desiredPortfolio := make(map[string]*pb.Position)
	totalRiskUnits := 0.0
	for _, signal := range signals {
		if signal.Signal == strategypb.SignalType_BUY {
			totalRiskUnits += signal.RiskUnit
		}
	}

	totalValue := s.calculateTotalValue() // This should include cash balance
	for _, signal := range signals {
		if signal.Signal == strategypb.SignalType_BUY {
			allocation := (signal.RiskUnit / totalRiskUnits) * totalValue
			quantity := int32(allocation / signal.CurrentPrice)
			desiredPortfolio[signal.Symbol] = &pb.Position{
				Symbol:       signal.Symbol,
				Quantity:     quantity,
				CurrentPrice: signal.CurrentPrice,
				MarketValue:  float64(quantity) * signal.CurrentPrice,
			}
		}
	}

	return desiredPortfolio, nil
}

func (s *Server) generateOrders(currentPortfolio, desiredPortfolio map[string]*pb.Position) []*pb.Order {
	var orders []*pb.Order

	for symbol, desiredPos := range desiredPortfolio {
		currentPos, exists := currentPortfolio[symbol]
		if !exists {
			// New position to open
			orders = append(orders, &pb.Order{
				Symbol:   symbol,
				Type:     pb.OrderType_MARKET,
				Quantity: desiredPos.Quantity,
			})
		} else if desiredPos.Quantity != currentPos.Quantity {
			// Existing position to adjust
			orderQuantity := desiredPos.Quantity - currentPos.Quantity
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
	for symbol, currentPos := range currentPortfolio {
		if _, exists := desiredPortfolio[symbol]; !exists {
			orders = append(orders, &pb.Order{
				Symbol:   symbol,
				Type:     pb.OrderType_MARKET,
				Quantity: -currentPos.Quantity,
			})
		}
	}

	return orders
}
