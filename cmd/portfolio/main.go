// cmd/portfolio/main.go
package main

import (
	"context"
	"math/rand"
	"net"

	"github.com/charmbracelet/log"

	pb "momentum-trading-platform/api/proto/portfolio_service"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPortfolioServiceServer
	portfolio map[string]*pb.Position
	cash      float64
}

func newServer() *server {
	return &server{
		portfolio: make(map[string]*pb.Position),
		cash:      1000000,
	}
}

func (s *server) GetPortfolioStatus(ctx context.Context, req *pb.PortfolioRequest) (*pb.PortfolioStatus, error) {
	positions := make([]*pb.Position, 0, len(s.portfolio))
	totalValue := s.cash

	for _, pos := range s.portfolio {
		positions = append(positions, pos)
		totalValue += pos.MarketValue
	}

	return &pb.PortfolioStatus{
		Positions:   positions,
		CashBalance: s.cash,
		TotalValue:  totalValue,
	}, nil
}

func (s *server) ProcessTradingSignals(ctx context.Context, signals *pb.TradingSignals) (*pb.PortfolioUpdate, error) {
	trades := make([]*pb.Trade, 0)

	for _, signal := range signals.Signals {
		switch signal.Type {
		case pb.TradeType_BUY:
			trade := s.executeBuy(signal.Symbol, signal.PositionSize)
			if trade != nil {
				trades = append(trades, trade)
			}
			log.Infof("Buy signal for %s with position size %.2f", signal.Symbol, signal.PositionSize)
		case pb.TradeType_SELL:
			trade := s.executeSell(signal.Symbol)
			if trade != nil {
				trades = append(trades, trade)
			}
			log.Infof("Sell signal for %s", signal.Symbol)
		}
	}

	status, _ := s.GetPortfolioStatus(ctx, &pb.PortfolioRequest{})
	return &pb.PortfolioUpdate{
		Trades:        trades,
		UpdatedStatus: status,
	}, nil
}

func (s *server) executeBuy(symbol string, positionSize float64) *pb.Trade {
	// Mock price between 100 and 200, this should be replaced with real data from Data Service
	price := 100 + rand.Float64()*100

	// Calculate quantity based on position size ($ terms) and price
	quantity := int32(positionSize / price)

	if quantity > 0 {
		cost := float64(quantity) * price
		if cost > s.cash {
			quantity = int32(s.cash / price)
			cost = float64(quantity) * price
		}
		s.cash -= cost

		if pos, exists := s.portfolio[symbol]; exists {
			pos.Quantity += quantity
			pos.AveragePrice = (pos.AveragePrice*float64(pos.Quantity-quantity) + cost) / float64(pos.Quantity)
			pos.CurrentPrice = price
			pos.MarketValue = float64(pos.Quantity) * price
		} else {
			s.portfolio[symbol] = &pb.Position{
				Symbol:       symbol,
				Quantity:     quantity,
				AveragePrice: price,
				CurrentPrice: price,
				MarketValue:  cost,
			}
		}
		log.Infof("Bought %d shares of %s at %.2f", quantity, symbol, price)

		return &pb.Trade{
			Symbol:   symbol,
			Type:     pb.TradeType_BUY,
			Quantity: quantity,
			Price:    price,
		}
	}

	return nil
}

func (s *server) executeSell(symbol string) *pb.Trade {
	if pos, exists := s.portfolio[symbol]; exists {
		price := pos.CurrentPrice * (1 + (rand.Float64()*0.1 - 0.05)) // Price +/- 5%
		s.cash += float64(pos.Quantity) * price

		trade := &pb.Trade{
			Symbol:   symbol,
			Type:     pb.TradeType_SELL,
			Quantity: pos.Quantity,
			Price:    price,
		}

		log.Infof("Sold %d shares of %s at %.2f", pos.Quantity, symbol, price)

		delete(s.portfolio, symbol)
		return trade
	}

	return nil
}
func (s *server) RebalancePortfolio(ctx context.Context, req *pb.RebalanceRequest) (*pb.PortfolioUpdate, error) {
	// Simplified rebalancing: sell 10% of each position
	trades := make([]*pb.Trade, 0)

	for symbol, pos := range s.portfolio {
		sellQuantity := int32(float64(pos.Quantity) * 0.1)
		if sellQuantity > 0 {
			trade := s.executeSell(symbol)
			if trade != nil {
				trades = append(trades, trade)
			}
		}
		log.Infof("Rebalanced %s, sold %d shares", symbol, sellQuantity)
	}

	status, _ := s.GetPortfolioStatus(ctx, &pb.PortfolioRequest{})
	return &pb.PortfolioUpdate{
		Trades:        trades,
		UpdatedStatus: status,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPortfolioServiceServer(s, newServer())
	log.Printf("Portfolio service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
