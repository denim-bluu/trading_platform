// cmd/trade_execution/main.go
package main

import (
	"context"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/charmbracelet/log"

	pb "momentum-trading-platform/api/proto/trade_execution_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedTradeExecutionServiceServer
	orders map[string]*pb.TradeResult
	mu     sync.Mutex
}

func newServer() *server {
	return &server{
		orders: make(map[string]*pb.TradeResult),
	}
}

func (s *server) ExecuteTrade(ctx context.Context, order *pb.TradeOrder) (*pb.TradeResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Simulate trade execution
	executedPrice := s.simulateExecutionPrice(order)
	executedQuantity := order.Quantity
	status := pb.ExecutionStatus_FILLED

	// Simulate partial fills for larger orders
	if order.Quantity > 1000 {
		executedQuantity = int32(float64(order.Quantity) * rand.Float64())
		if executedQuantity < order.Quantity {
			status = pb.ExecutionStatus_PARTIALLY_FILLED
		}
	}

	result := &pb.TradeResult{
		OrderId:          order.OrderId,
		Status:           status,
		ExecutedPrice:    executedPrice,
		ExecutedQuantity: executedQuantity,
		ExecutionTime:    time.Now().Format(time.RFC3339),
	}

	log.Infof("Executed trade: %+v", result)

	s.orders[order.OrderId] = result

	return result, nil
}

func (s *server) GetExecutionStatus(ctx context.Context, req *pb.ExecutionStatusRequest) (*pb.TradeResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if result, exists := s.orders[req.OrderId]; exists {
		log.Infof("Order status for %s: %+v", req.OrderId, result)
		return result, nil
	}

	return nil, status.Errorf(codes.NotFound, "order not found")
}

func (s *server) simulateExecutionPrice(order *pb.TradeOrder) float64 {
	if order.Type == pb.OrderType_LIMIT {
		// For limit orders, execute at limit price or better
		return order.LimitPrice * (1 - rand.Float64()*0.01) // Up to 1% better than limit
	}

	// For market orders, simulate some slippage
	basePrice := 100.0                     // Assume a base price of $100
	slippage := rand.Float64()*0.02 - 0.01 // -1% to +1% slippage
	return basePrice * (1 + slippage)
}

func main() {
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTradeExecutionServiceServer(s, newServer())
	log.Printf("Trade Execution service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
