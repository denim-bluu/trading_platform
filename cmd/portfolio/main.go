// cmd/portfolio/main.go
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	datapb "momentum-trading-platform/api/proto/data_service"
	pb "momentum-trading-platform/api/proto/portfolio_service"
	statepb "momentum-trading-platform/api/proto/portfolio_state_service"
	strategypb "momentum-trading-platform/api/proto/strategy_service"
	"momentum-trading-platform/utils"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedPortfolioServiceServer
	logger         *log.Logger
	dataClient     datapb.DataServiceClient
	strategyClient strategypb.StrategyServiceClient
	stateClient    statepb.PortfolioStateServiceClient
	portfolio      map[string]*pb.Position
	cashBalance    float64
	lastUpdateDate string
	mu             sync.Mutex
}

func newServer(dataClient datapb.DataServiceClient, strategyClient strategypb.StrategyServiceClient, stateClient statepb.PortfolioStateServiceClient) *server {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)

	s := &server{
		logger:         logger,
		dataClient:     dataClient,
		strategyClient: strategyClient,
		stateClient:    stateClient,
		portfolio:      make(map[string]*pb.Position),
		lastUpdateDate: time.Now().Format("2006-01-02:15:04:05"),
	}

	// Load the last saved state
	if err := s.loadLastState(); err != nil {
		s.logger.WithError(err).Error("Failed to load last portfolio state")
	}

	return s
}

// Weekly Rebalance (WeeklyRebalance method):
//
//   - Occurs every Wednesday
//
//   - Steps:
//     a. Check if the S&P 500 index is in a positive trend (above 200-day MA)
//     b. Get signals from the Strategy Service
//     c. Sell positions: Sell any position not present in the current signals
//     d. If index is in a positive trend, for each signal:
//     d.1. If it's an existing position: Adjust the position size
//     d.2. If it's a new position: Buy it (if cash is available)
//     d.3. Update the portfolio status
func (s *server) WeeklyRebalance(ctx context.Context, req *pb.RebalanceRequest) (*pb.PortfolioUpdate, error) {
	s.logger.WithField("date", req.Date).Info("Performing weekly rebalance")
	if err := s.loadLastState(); err != nil {
		return nil, fmt.Errorf("failed to load latest state before rebalance: %v", err)
	}

	if !isWednesday(req.Date) {
		return nil, fmt.Errorf("weekly rebalance can only be performed on Wednesdays")
	}

	isIndexPositive, err := s.isIndexInPositiveTrend(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check index trend: %v", err)
	}

	signals, err := s.getStrategySignals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get strategy signals: %v", err)
	}

	var trades []*pb.Trade

	// Sell positions that are no longer in the signals
	for symbol := range s.portfolio {
		if !s.isInSignals(symbol, signals.Signals) {
			trade := s.sellPosition(symbol)
			trades = append(trades, trade)
		}
	}

	// Buy new positions or adjust existing ones if index is in positive trend
	if isIndexPositive {
		for _, signal := range signals.Signals {
			if position, exists := s.portfolio[signal.Symbol]; exists {
				trade := s.adjustPosition(signal, position)
				if trade != nil {
					trades = append(trades, trade)
				}
			} else if s.getPortfolioStatus().GetCashBalance() > 0 {
				trade := s.buyPosition(signal)
				if trade != nil {
					trades = append(trades, trade)
				}
			}
		}
	}

	s.lastUpdateDate = req.Date
	if err := s.saveState(); err != nil {
		return nil, fmt.Errorf("failed to save state after rebalance: %v", err)
	}

	return &pb.PortfolioUpdate{
		Trades:        trades,
		UpdatedStatus: s.getPortfolioStatus(),
	}, nil
}

// BiWeekly Rebalance:
//
//   - Occurs on the second Wednesday of each month
//
//   - Steps:
//     a. Perform the regular weekly rebalance
//     b. Get fresh signals from the strategy service
//     c. For all current positions, adjust position sizes based on the new signals
//     d. Update the portfolio status
func (s *server) BiWeeklyRebalance(ctx context.Context, req *pb.RebalanceRequest) (*pb.PortfolioUpdate, error) {
	s.logger.WithField("date", req.Date).Info("Performing bi-weekly rebalance")
	if err := s.loadLastState(); err != nil {
		return nil, fmt.Errorf("failed to load latest state before rebalance: %v", err)
	}

	if !isSecondWednesdayOfMonth(req.Date) {
		return nil, fmt.Errorf("bi-weekly rebalance can only be performed on the second Wednesday of the month")
	}

	// Perform the regular weekly rebalance
	weeklyUpdate, err := s.WeeklyRebalance(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform weekly rebalance: %v", err)
	}

	// Then, adjust all position sizes
	signals, err := s.getStrategySignals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get strategy signals: %v", err)
	}

	var additionalTrades []*pb.Trade

	for _, signal := range signals.Signals {
		if position, exists := s.portfolio[signal.Symbol]; exists {
			trade := s.adjustPosition(signal, position)
			if trade != nil {
				additionalTrades = append(additionalTrades, trade)
			}
		}
	}

	weeklyUpdate.Trades = append(weeklyUpdate.Trades, additionalTrades...)
	weeklyUpdate.UpdatedStatus = s.getPortfolioStatus()

	if err := s.saveState(); err != nil {
		return nil, fmt.Errorf("failed to save state after rebalance: %v", err)
	}

	return weeklyUpdate, nil
}

func (s *server) sellPosition(symbol string) *pb.Trade {
	position := s.portfolio[symbol]
	trade := &pb.Trade{
		Symbol:   symbol,
		Type:     pb.TradeType_SELL,
		Quantity: position.Quantity,
		Price:    position.CurrentPrice,
	}
	s.cashBalance += float64(position.Quantity) * position.CurrentPrice
	delete(s.portfolio, symbol)
	return trade
}

func (s *server) buyPosition(signal *strategypb.StockSignal) *pb.Trade {
	quantity := utils.CalculatePositionSize(signal.RiskUnit, signal.CurrentPrice)
	cost := float64(quantity) * signal.CurrentPrice
	if cost > s.cashBalance {
		quantity = int32(s.cashBalance / signal.CurrentPrice)
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
	s.portfolio[signal.Symbol] = &pb.Position{
		Symbol:       signal.Symbol,
		Quantity:     quantity,
		AveragePrice: signal.CurrentPrice,
		CurrentPrice: signal.CurrentPrice,
		MarketValue:  cost,
	}
	s.cashBalance -= cost
	return trade
}

func (s *server) adjustPosition(signal *strategypb.StockSignal, position *pb.Position) *pb.Trade {
	targetSize := utils.CalculatePositionSize(signal.RiskUnit, signal.CurrentPrice)
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
		s.cashBalance -= float64(diff) * signal.CurrentPrice
	} else {
		position.Quantity -= diff
		s.cashBalance += float64(diff) * signal.CurrentPrice
	}
	position.CurrentPrice = signal.CurrentPrice
	position.MarketValue = float64(position.Quantity) * signal.CurrentPrice
	return trade
}

func (s *server) isInSignals(symbol string, signals []*strategypb.StockSignal) bool {
	for _, signal := range signals {
		if signal.Symbol == symbol {
			return true
		}
	}
	return false
}

func (s *server) getTotalPortfolioValue() float64 {
	total := s.cashBalance
	for _, position := range s.portfolio {
		total += position.MarketValue
	}
	return total
}

func (s *server) getPortfolioStatus() *pb.PortfolioStatus {
	positions := make([]*pb.Position, 0, len(s.portfolio))
	for _, position := range s.portfolio {
		positions = append(positions, position)
	}
	return &pb.PortfolioStatus{
		Positions:      positions,
		CashBalance:    s.cashBalance,
		TotalValue:     s.getTotalPortfolioValue(),
		LastUpdateDate: s.lastUpdateDate,
	}
}

func (s *server) isIndexInPositiveTrend(ctx context.Context) (bool, error) {
	resp, err := s.dataClient.GetStockData(ctx, &datapb.StockRequest{
		Symbol:    "^GSPC", // S&P 500 index
		StartDate: fmt.Sprintf("%d", time.Now().AddDate(0, 0, -200).Unix()),
		EndDate:   fmt.Sprintf("%d", time.Now().Unix()),
		Interval:  "1d",
	})
	if err != nil {
		return false, err
	}

	ma200 := utils.CalculateMovingAverage(resp.DataPoints, 200)
	currentPrice := resp.DataPoints[len(resp.DataPoints)-1].Close
	return currentPrice > ma200, nil
}

func (s *server) getStrategySignals(ctx context.Context) (*strategypb.SignalResponse, error) {
	// Implementation to get signals from the strategy service
	return s.strategyClient.GenerateSignals(ctx, &strategypb.SignalRequest{
		Symbols:   s.getSymbolsToAnalyze(),
		StartDate: time.Now().AddDate(0, 0, -90).Format("2006-01-02"),
		EndDate:   time.Now().Format("2006-01-02"),
		Interval:  "1d",
	})
}

func (s *server) getSymbolsToAnalyze() []string {
	// This should return all symbols in the S&P 500 index
	// Placeholder implementation
	return []string{"AAPL", "GOOGL", "MSFT", "AMZN", "FB"}
}

func isWednesday(date string) bool {
	t, _ := time.Parse("2006-01-02", date)
	return t.Weekday() == time.Wednesday
}

func isSecondWednesdayOfMonth(date string) bool {
	t, _ := time.Parse("2006-01-02", date)
	return t.Weekday() == time.Wednesday && (t.Day()-1)/7 == 1
}

func (s *server) convertStatePositions(positions []*statepb.Position) map[string]*pb.Position {
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

func (s *server) loadLastState() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	lastState, err := s.stateClient.LoadPortfolioState(ctx, &statepb.LoadRequest{
		Date: time.Now().Format("2006-01-02"), // Or the last known date
	})
	if err != nil {
		return fmt.Errorf("failed to load last portfolio state: %v", err)
	}
	s.logger.Infof("Loaded last portfolio state: %s", lastState.Date)

	s.portfolio = s.convertStatePositions(lastState.Positions)
	s.cashBalance = lastState.CashBalance
	s.lastUpdateDate = lastState.Date

	return nil
}

func (s *server) saveState() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	positions := make([]*statepb.Position, 0, len(s.portfolio))
	for _, pos := range s.portfolio {
		positions = append(positions, &statepb.Position{
			Symbol:       pos.Symbol,
			Quantity:     pos.Quantity,
			AveragePrice: pos.AveragePrice,
			CurrentPrice: pos.CurrentPrice,
			MarketValue:  pos.MarketValue,
		})
	}

	_, err := s.stateClient.SavePortfolioState(ctx, &statepb.PortfolioState{
		Date:        time.Now().Format("2006-01-02"),
		Positions:   positions,
		CashBalance: s.cashBalance,
		TotalValue:  s.getTotalPortfolioValue(),
	})

	return err
}

func main() {
	// Set up connections to the data and strategy services
	dataConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to data service: %v", err)
	}
	defer dataConn.Close()
	dataClient := datapb.NewDataServiceClient(dataConn)

	strategyConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to strategy service: %v", err)
	}
	defer strategyConn.Close()
	strategyClient := strategypb.NewStrategyServiceClient(strategyConn)

	stateConn, err := grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to strategy service: %v", err)
	}
	defer strategyConn.Close()
	stateClient := statepb.NewPortfolioStateServiceClient(stateConn)

	s := newServer(dataClient, strategyClient, stateClient)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		s.logger.WithError(err).Fatal("Failed to listen")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPortfolioServiceServer(grpcServer, s)

	s.logger.WithField("address", lis.Addr().String()).Info("Portfolio service starting")
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.WithError(err).Fatal("Failed to serve")
	}
}
