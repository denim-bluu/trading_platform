package strategy

import (
	"context"
	"fmt"
	datapb "momentum-trading-platform/api/proto/data_service"
	pb "momentum-trading-platform/api/proto/strategy_service"

	log "github.com/sirupsen/logrus"
)

func (s *Server) GenerateSignals(ctx context.Context, req *pb.SignalRequest) (*pb.SignalResponse, error) {
	s.Logger.WithFields(log.Fields{
		"symbols":     req.Symbols,
		"start":       req.StartDate,
		"end":         req.EndDate,
		"interval":    req.Interval,
		"marketIndex": req.MarketIndex,
	}).Info("Generating signals")

	// Fetch market index data (e.g., S&P 500)
	indexResp, err := s.fetchIndexData(ctx, req.MarketIndex, req.StartDate, req.EndDate, req.Interval)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch index data: %v", err)
	}

	batchResp, err := s.fetchBatchStockData(ctx, req)
	if err != nil {
		return nil, err
	}

	// Use the momentum strategy (or allow selection of strategy in the future)
	strategy := s.Strategies["momentum"]
	signals, err := strategy.GenerateSignals(batchResp.StockData, indexResp)
	if err != nil {
		return nil, err
	}

	return &pb.SignalResponse{
		Signals: signals,
	}, nil
}

func (s *Server) fetchBatchStockData(ctx context.Context, req *pb.SignalRequest) (*datapb.BatchStockResponse, error) {
	batchReq := &datapb.BatchStockRequest{
		Symbols:   req.Symbols,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Interval:  req.Interval,
	}
	s.Logger.Infof("📡 Fetching following stock data: %v for %v to %v", req.Symbols, req.StartDate, req.EndDate)

	batchResp, err := s.Clients.DataClient.GetBatchStockData(ctx, batchReq)
	if err != nil {
		s.Logger.WithError(err).Error("❌ Failed to fetch batch stock data")
		return nil, fmt.Errorf("❌ failed to fetch stock data: %v", err)
	}
	if batchResp.Errors != nil {
		for symbol, errMsg := range batchResp.Errors {
			s.Logger.WithFields(log.Fields{
				"symbol": symbol,
				"error":  errMsg,
			}).Warn("❗ Failed to fetch stock data")
		}
	}
	return batchResp, nil
}

func (s *Server) fetchIndexData(ctx context.Context, indexSymbol, startDate, endDate, interval string) (*datapb.StockResponse, error) {
	s.Logger.Infof("📡 Fetching index data for %s from %v to %v", indexSymbol, startDate, endDate)
	indexReq := &datapb.StockRequest{
		Symbol:    indexSymbol,
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  interval,
	}
	return s.Clients.DataClient.GetStockData(ctx, indexReq)
}
