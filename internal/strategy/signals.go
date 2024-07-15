package strategy

import (
	"context"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"

	datapb "momentum-trading-platform/api/proto/data_service"
	pb "momentum-trading-platform/api/proto/strategy_service"
	"momentum-trading-platform/internal/utils"
)

func (s *Server) GenerateSignals(ctx context.Context, req *pb.SignalRequest) (*pb.SignalResponse, error) {
	s.Logger.WithFields(log.Fields{
		"symbols":  req.Symbols,
		"start":    req.StartDate,
		"end":      req.EndDate,
		"interval": req.Interval,
	}).Info("📧 Sending request to data service")

	batchResp, err := s.fetchBatchStockData(ctx, req)
	if err != nil {
		return nil, err
	}

	signals := s.processStockData(batchResp)

	signals = s.sortAndFilterSignalsbyMomentum(signals)

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

	batchResp, err := s.Clients.DataClient.GetBatchStockData(ctx, batchReq)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch batch stock data")
		return nil, fmt.Errorf("failed to fetch stock data: %v", err)
	}
	s.Logger.Info("✅ Received batch stock data")
	return batchResp, nil
}

func (s *Server) processStockData(batchResp *datapb.BatchStockResponse) []*pb.StockSignal {
	var signals []*pb.StockSignal
	for symbol, stockResp := range batchResp.StockData {
		s.Logger.WithField("symbol", symbol).Info("📊 Processing stock data")

		if s.isStockDisqualified(stockResp) {
			continue
		}

		signal := s.determineStockSignal(symbol, stockResp)
		signals = append(signals, signal)
	}
	return signals
}

func (s *Server) isStockDisqualified(stockResp *datapb.StockResponse) bool {
	if utils.HasRecentLargeGap(stockResp.DataPoints, 90, 0.15) {
		s.Logger.WithField("symbol", stockResp.Symbol).Info("❌ Stock disqualified due to recent large gap")
		return true
	}

	lastPrice := stockResp.DataPoints[len(stockResp.DataPoints)-1].Close
	movingAverage := utils.CalculateMovingAverage(stockResp.DataPoints, 100)
	if lastPrice < movingAverage {
		s.Logger.WithField("symbol", stockResp.Symbol).Info("❌ Stock disqualified due to being below 100MA")
		return true
	}

	momentumScore := utils.CalculateMomentumScore(stockResp.DataPoints, 90)
	if momentumScore < 0 {
		s.Logger.WithField("symbol", stockResp.Symbol).Info("❌ Stock disqualified due to negative momentum score")
		return true
	}

	return false
}

func (s *Server) determineStockSignal(symbol string, stockResp *datapb.StockResponse) *pb.StockSignal {
	lastPrice := stockResp.DataPoints[len(stockResp.DataPoints)-1].Close
	momentumScore := utils.CalculateMomentumScore(stockResp.DataPoints, 90)
	atr := utils.CalculateATR(stockResp.DataPoints, 20)
	riskUnit := utils.CalculateRiskUnit(atr, 0.001)

	s.Logger.WithFields(log.Fields{
		"symbol":         symbol,
		"momentum_score": momentumScore,
		"atr":            atr,
		"risk_unit":      riskUnit,
	}).Trace("🔢 Stock metrics")

	return &pb.StockSignal{
		Symbol:        symbol,
		Signal:        pb.SignalType_BUY,
		RiskUnit:      riskUnit,
		MomentumScore: momentumScore,
		CurrentPrice:  lastPrice,
	}
}

func (s *Server) sortAndFilterSignalsbyMomentum(signals []*pb.StockSignal) []*pb.StockSignal {
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].MomentumScore > signals[j].MomentumScore
	})
	s.Logger.WithField("signals", signals).Info("Rank ordered Signals")

	topCount := int(float64(len(signals)) * 0.2)
	s.Logger.WithField("top_count", topCount).Info("Filtering top 20% of signals")
	return signals[:topCount]
}
