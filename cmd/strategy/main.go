// cmd/strategy/main.go
package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sort"
	"time"

	pb "momentum-trading-platform/api/proto/strategy_service"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStrategyServiceServer
}

type Stock struct {
	Symbol        string
	MomentumScore float64
}

var mockStocks = []string{"AAPL", "GOOGL", "MSFT", "AMZN", "FB", "TSLA", "NVDA", "JPM", "JNJ", "V"}

func (s *server) GetTradingSignals(ctx context.Context, req *pb.SignalRequest) (*pb.SignalResponse, error) {
	// Mock market regime (70% chance of positive regime)
	isMarketRegimePositive := rand.Float32() < 0.7

	// Generate mock momentum scores
	stocks := make([]Stock, len(mockStocks))
	for i, symbol := range mockStocks {
		stocks[i] = Stock{
			Symbol:        symbol,
			MomentumScore: rand.Float64() * 100, // Random score between 0 and 100
		}
	}

	// Sort stocks by momentum score
	sort.Slice(stocks, func(i, j int) bool {
		return stocks[i].MomentumScore > stocks[j].MomentumScore
	})

	// Generate signals
	signals := make([]*pb.StockSignal, 0)
	for i, stock := range stocks {
		var signal pb.SignalType
		if i < len(stocks)/5 && isMarketRegimePositive { // Top 20% and positive market regime
			signal = pb.SignalType_BUY
		} else if i >= len(stocks)/2 { // Bottom 50%
			signal = pb.SignalType_SELL
		} else {
			signal = pb.SignalType_HOLD
		}

		signals = append(signals, &pb.StockSignal{
			Symbol:        stock.Symbol,
			Signal:        signal,
			MomentumScore: stock.MomentumScore,
			PositionSize:  100000 / stock.MomentumScore, // Simplified position sizing
		})
	}

	return &pb.SignalResponse{
		Signals:                signals,
		IsMarketRegimePositive: isMarketRegimePositive,
	}, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStrategyServiceServer(s, &server{})
	log.Printf("Strategy service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
