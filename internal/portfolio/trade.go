package portfolio

import (
	"context"
	pb "momentum-trading-platform/api/proto/portfolio_service"
	strategypb "momentum-trading-platform/api/proto/strategy_service"
	"momentum-trading-platform/internal/utils"
)

func (s *Server) sellPosition(symbol string) *pb.Trade {
	s.Logger.WithField("symbol", symbol).Info("Selling position")
	s.mu.Lock()
	defer s.mu.Unlock()

	position := s.Portfolio[symbol]
	trade := &pb.Trade{
		Symbol:   symbol,
		Type:     pb.TradeType_SELL,
		Quantity: position.Quantity,
		Price:    position.CurrentPrice,
	}
	s.CashBalance += float64(position.Quantity) * position.CurrentPrice
	delete(s.Portfolio, symbol)

	if err := s.Storage.SaveTrade(context.Background(), trade, s.CashBalance); err != nil {
		s.Logger.WithError(err).Error("Failed to save trade")
	}
	return trade
}

func (s *Server) buyPosition(signal *strategypb.StockSignal) *pb.Trade {
	s.mu.Lock()
	defer s.mu.Unlock()

	quantity := utils.CalculatePositionSize(signal.RiskUnit, s.CashBalance)
	cost := float64(quantity) * signal.CurrentPrice
	s.Logger.Infof("Buying %d shares of %s at $%.2f", quantity, signal.Symbol, signal.CurrentPrice)
	if cost > s.CashBalance {
		quantity = int32(s.CashBalance / signal.CurrentPrice)
		cost = float64(quantity) * signal.CurrentPrice
	}
	if quantity == 0 {
		return nil
	}
	trade := &pb.Trade{
		Symbol:   signal.Symbol,
		Type:     pb.TradeType_BUY,
		Quantity: quantity,
		Price:    signal.CurrentPrice,
	}
	s.Portfolio[signal.Symbol] = &pb.Position{
		Symbol:       signal.Symbol,
		Quantity:     quantity,
		AveragePrice: signal.CurrentPrice,
		CurrentPrice: signal.CurrentPrice,
		MarketValue:  cost,
	}
	s.CashBalance -= cost

	s.Logger.WithField("Trade", trade).Info("Trade executed")

	if err := s.Storage.SaveTrade(context.Background(), trade, s.CashBalance); err != nil {
		s.Logger.WithError(err).Error("Failed to save trade")
	}
	return trade
}

func (s *Server) adjustPosition(signal *strategypb.StockSignal, position *pb.Position) *pb.Trade {
	s.Logger.WithField("symbol", signal.Symbol).Info("Adjusting position")
	s.mu.Lock()
	defer s.mu.Unlock()

	targetSize := utils.CalculatePositionSize(signal.RiskUnit, s.CashBalance)
	diff := targetSize - position.Quantity
	if diff == 0 {
		return nil
	}
	tradeType := pb.TradeType_BUY
	if diff < 0 {
		tradeType = pb.TradeType_SELL
		diff = -diff
	}
	trade := &pb.Trade{
		Symbol:   signal.Symbol,
		Type:     tradeType,
		Quantity: diff,
		Price:    signal.CurrentPrice,
	}
	if tradeType == pb.TradeType_BUY {
		position.Quantity += diff
		s.CashBalance -= float64(diff) * signal.CurrentPrice
	} else {
		position.Quantity -= diff
		s.CashBalance += float64(diff) * signal.CurrentPrice
	}
	position.CurrentPrice = signal.CurrentPrice
	position.MarketValue = float64(position.Quantity) * signal.CurrentPrice

	if err := s.Storage.SaveTrade(context.Background(), trade, s.CashBalance); err != nil {
		s.Logger.WithError(err).Error("Failed to save trade")
	}
	return trade
}
