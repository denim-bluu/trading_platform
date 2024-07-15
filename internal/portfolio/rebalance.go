package portfolio

import (
	"context"
	"fmt"
	"time"

	datapb "momentum-trading-platform/api/proto/data_service"
	pb "momentum-trading-platform/api/proto/portfolio_service"
	strategypb "momentum-trading-platform/api/proto/strategy_service"
	"momentum-trading-platform/internal/utils"
)

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
func (s *Server) WeeklyRebalance(ctx context.Context, req *pb.RebalanceRequest) (*pb.PortfolioUpdate, error) {
	s.Logger.WithField("date", req.Date).Info("Performing weekly rebalance")
	if err := s.loadLastState(); err != nil {
		return nil, fmt.Errorf("failed to load latest state before rebalance: %v", err)
	}

	if !utils.IsWednesday(req.Date) {
		// return nil, fmt.Errorf("weekly rebalance can only be performed on Wednesdays")
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
	for symbol := range s.Portfolio {
		if !s.isSymbolInSignals(symbol, signals.Signals) {
			trade := s.sellPosition(symbol)
			trades = append(trades, trade)
		}
	}

	// Buy new positions or adjust existing ones if index is in positive trend
	if isIndexPositive {
		for _, signal := range signals.Signals {
			if position, exists := s.Portfolio[signal.Symbol]; exists {
				trade := s.adjustPosition(signal, position)
				if trade != nil {
					trades = append(trades, trade)
				}
			} else if s.CashBalance > 0 {
				trade := s.buyPosition(signal)
				if trade != nil {
					trades = append(trades, trade)
				}
			}
		}
	}

	s.LastUpdateDate = req.Date
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
func (s *Server) BiWeeklyRebalance(ctx context.Context, req *pb.RebalanceRequest) (*pb.PortfolioUpdate, error) {
	s.Logger.WithField("date", req.Date).Info("Performing bi-weekly rebalance")
	if err := s.loadLastState(); err != nil {
		return nil, fmt.Errorf("failed to load latest state before rebalance: %v", err)
	}

	if !utils.IsSecondWednesdayOfMonth(req.Date) {
		// return nil, fmt.Errorf("bi-weekly rebalance can only be performed on the second Wednesday of the month")
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
		if position, exists := s.Portfolio[signal.Symbol]; exists {
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

func (s *Server) isIndexInPositiveTrend(ctx context.Context) (bool, error) {
	resp, err := s.Clients.DataClient.GetStockData(ctx, &datapb.StockRequest{
		Symbol:    "^GSPC", // S&P 500 index
		StartDate: time.Now().AddDate(0, 0, -200).Format("2006-01-02"),
		EndDate:   time.Now().Format("2006-01-02"),
		Interval:  "1d",
	})
	if err != nil {
		return false, err
	}

	ma200 := utils.CalculateMovingAverage(resp.DataPoints, 200)
	currentPrice := resp.DataPoints[len(resp.DataPoints)-1].Close
	return currentPrice > ma200, nil
}

func (s *Server) getStrategySignals(ctx context.Context) (*strategypb.SignalResponse, error) {
	return s.Clients.StrategyClient.GenerateSignals(ctx, &strategypb.SignalRequest{
		Symbols:   s.getSymbolsToAnalyze(),
		StartDate: time.Now().AddDate(0, 0, -90).Format("2006-01-02"),
		EndDate:   time.Now().Format("2006-01-02"),
		Interval:  "1d",
	})
}

func (s *Server) getSymbolsToAnalyze() []string {
	// This should return all symbols in the S&P 500 index
	// Placeholder implementation
	return []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "NVDA", "NFLX", "PYPL", "ADBE",
		"INTC", "CSCO", "CMCSA", "PEP", "AVGO", "TXN", "COST", "QCOM", "TMUS", "AMGN", "SBUX",
		"INTU", "AMD", "ISRG", "GILD", "MDLZ", "BKNG", "MU", "ADP", "REGN", "ATVI"}
}

func (s *Server) isSymbolInSignals(symbol string, signals []*strategypb.StockSignal) bool {
	for _, signal := range signals {
		if signal.Symbol == symbol {
			return true
		}
	}
	return false
}
