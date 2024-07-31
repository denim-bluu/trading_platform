package tradeexecution

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	portfoliostatepb "momentum-trading-platform/api/proto/portfolio_state_service"
	pb "momentum-trading-platform/api/proto/trade_execution_service"
)

type Server struct {
	pb.UnimplementedTradeExecutionServiceServer
	Logger     *log.Logger
	Clients    *Clients
	executions map[string]*pb.ExecutionStatus
	mu         sync.Mutex
}

func NewServer(clients *Clients) *Server {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	return &Server{
		Logger:     logger,
		Clients:    clients,
		executions: make(map[string]*pb.ExecutionStatus),
	}
}

func (s *Server) ExecuteTrades(ctx context.Context, req *pb.ExecuteTradesRequest) (*pb.ExecuteTradesResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	executionID := uuid.New().String()
	s.Logger.WithField("executionID", executionID).Info("Executing trades")

	results := make([]*pb.OrderExecutionResult, len(req.Orders))
	for i, order := range req.Orders {
		// Simulate trade execution
		result := s.simulateExecution(order)
		results[i] = result
	}

	executionStatus := &pb.ExecutionStatus{
		ExecutionId: executionID,
		Status:      pb.ExecutionStatusType_COMPLETED,
		Results:     results,
	}
	s.executions[executionID] = executionStatus

	// Update portfolio state
	err := s.updatePortfolioState(ctx, results)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to update portfolio state")
		return nil, status.Errorf(codes.Internal, "failed to update portfolio state: %v", err)
	}

	return &pb.ExecuteTradesResponse{
		ExecutionId: executionID,
		Results:     results,
	}, nil
}

func (s *Server) GetExecutionStatus(ctx context.Context, req *pb.GetExecutionStatusRequest) (*pb.ExecutionStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	executionStatus, exists := s.executions[req.ExecutionId]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "execution ID not found")
	}

	return executionStatus, nil
}

func (s *Server) updatePortfolioState(ctx context.Context, results []*pb.OrderExecutionResult) error {
	// First, get the current portfolio state
	currentState, err := s.Clients.PortfolioStateClient.GetPortfolioState(ctx, &portfoliostatepb.GetPortfolioStateRequest{})
	if err != nil {
		return fmt.Errorf("failed to get current portfolio state: %w", err)
	}

	// Update the positions based on execution results
	positions := currentState.Positions
	cashBalance := currentState.CashBalance

	for _, result := range results {
		found := false
		for i, position := range positions {
			if position.Symbol == result.Symbol {
				newQuantity := position.Quantity + result.FilledQuantity
				if newQuantity == 0 {
					// Remove the position if quantity becomes 0
					positions = append(positions[:i], positions[i+1:]...)
				} else {
					position.Quantity = newQuantity
					position.CurrentPrice = result.AveragePrice
					position.MarketValue = float64(newQuantity) * result.AveragePrice
				}
				found = true
				break
			}
		}
		if !found && result.FilledQuantity > 0 {
			// Add new position
			positions = append(positions, &portfoliostatepb.Position{
				Symbol:       result.Symbol,
				Quantity:     result.FilledQuantity,
				CurrentPrice: result.AveragePrice,
				MarketValue:  float64(result.FilledQuantity) * result.AveragePrice,
			})
		}
		// Update cash balance
		cashBalance -= float64(result.FilledQuantity) * result.AveragePrice
	}

	// Update the portfolio state
	_, err = s.Clients.PortfolioStateClient.UpdatePortfolioState(ctx, &portfoliostatepb.UpdatePortfolioStateRequest{
		Positions:   positions,
		CashBalance: cashBalance,
	})
	if err != nil {
		return fmt.Errorf("failed to update portfolio state: %w", err)
	}

	return nil
}

func (s *Server) simulateExecution(order *pb.Order) *pb.OrderExecutionResult {
	// TODO: Replace this with actual broker API interaction
	s.Logger.WithFields(log.Fields{
		"symbol":   order.Symbol,
		"quantity": order.Quantity,
		"type":     order.Type,
	}).Info("Executing order")

	// Simulate some latency
	time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond)

	filledQuantity := order.Quantity
	averagePrice := order.Price

	if order.Type == pb.OrderType_MARKET {
		// Simulate some slippage for market orders
		slippage := (1 + (0.01 * (2*rand.Float64() - 1))) // +/- 1% slippage
		averagePrice *= slippage
	}

	// Simulate partial fills
	if rand.Float32() < 0.1 { // 10% chance of partial fill
		filledQuantity = int32(float32(order.Quantity) * rand.Float32())
	}

	return &pb.OrderExecutionResult{
		Symbol:         order.Symbol,
		Status:         pb.ExecutionStatusType_COMPLETED,
		FilledQuantity: filledQuantity,
		AveragePrice:   averagePrice,
	}
}
