// cmd/strategy/main.go
package main

import (
	"context"
	"math/rand"
	"net"

	"github.com/charmbracelet/log"

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
	isMarketRegimePositive := rand.Float32() < 0.7

	signals := make([]*pb.StockSignal, 0)
	for _, symbol := range mockStocks {
		signalType := pb.SignalType_HOLD
		if isMarketRegimePositive && rand.Float32() < 0.3 {
			signalType = pb.SignalType_BUY
		} else if !isMarketRegimePositive && rand.Float32() < 0.3 {
			signalType = pb.SignalType_SELL
		}

		// Calculate position size based on risk parity
		atr := 1.0 + rand.Float64()*4.0       // Mock ATR between 1 and 5
		positionSize := 1000000 * 0.001 / atr // Assuming $1M account and 10 bps risk per stock

		log.Info("Signal for %s: %s, Position Size: %.2f", symbol, signalType, positionSize)
		signals = append(signals, &pb.StockSignal{
			Symbol:       symbol,
			Signal:       signalType,
			PositionSize: positionSize,
		})
	}

	return &pb.SignalResponse{
		Signals:                signals,
		IsMarketRegimePositive: isMarketRegimePositive,
	}, nil
}

func main() {
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
