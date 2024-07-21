package strategy

import (
	datapb "momentum-trading-platform/api/proto/data_service"
	pb "momentum-trading-platform/api/proto/strategy_service"
	"momentum-trading-platform/internal/utils"
	"sort"

	log "github.com/sirupsen/logrus"
)

// MomentumStrategy defines the structure of the momentum trading strategy.
type MomentumStrategy struct {
	lookbackPeriod int
	topPercentage  float64
	riskFactor     float64
	ma200Period    int // Period for calculating 200-day moving average
}

func NewMomentumStrategy() *MomentumStrategy {
	return &MomentumStrategy{
		lookbackPeriod: 90,
		topPercentage:  0.2,
		riskFactor:     0.001,
		ma200Period:    200,
	}
}

// GenerateSignals generates trading signals based on the provided batch of stock data.
func (s *MomentumStrategy) GenerateSignals(batchStockData map[string]*datapb.StockResponse, indexData *datapb.StockResponse) ([]*pb.StockSignal, error) {
	regime := s.DetectMarketRegime(indexData)

	if regime == Bear {
		// In a bear market, we don't trade
		return []*pb.StockSignal{}, nil
	}

	var signals []*pb.StockSignal

	for symbol, data := range batchStockData {
		signal := s.generateSignal(symbol, data)
		signals = append(signals, signal)
	}

	return s.sortAndFilterSignals(signals), nil
}

// generateSignal generates a trading signal for a specific stock.
func (s *MomentumStrategy) generateSignal(symbol string, stockResp *datapb.StockResponse) *pb.StockSignal {
	lastPrice := stockResp.DataPoints[len(stockResp.DataPoints)-1].Close
	momentumScore := utils.CalculateMomentumScore(stockResp.DataPoints, s.lookbackPeriod)
	riskUnit := s.CalculateRisk(stockResp)

	if s.isStockDisqualified(stockResp) {
		return &pb.StockSignal{
			Symbol:        symbol,
			Signal:        pb.SignalType_SELL,
			RiskUnit:      -999, // Set to -999 to indicate that the stock is disqualified, and clear all positions
			MomentumScore: momentumScore,
			CurrentPrice:  lastPrice,
		}
	}

	return &pb.StockSignal{
		Symbol:        symbol,
		Signal:        pb.SignalType_BUY,
		RiskUnit:      riskUnit,
		MomentumScore: momentumScore,
		CurrentPrice:  lastPrice,
	}
}

// isStockDisqualified checks if a stock is disqualified based on certain criteria.
func (s *MomentumStrategy) isStockDisqualified(stockResp *datapb.StockResponse) bool {
	// Implementation similar to the current isStockDisqualified function
	if utils.HasRecentLargeGap(stockResp.DataPoints, 90, 0.15) {
		log.WithField("symbol", stockResp.Symbol).Info("❌ Stock disqualified due to recent large gap")
		return true
	}
	lastPrice := stockResp.DataPoints[len(stockResp.DataPoints)-1].Close
	movingAverage := utils.CalculateMovingAverage(stockResp.DataPoints, 100)
	if lastPrice < movingAverage {
		log.WithField("symbol", stockResp.Symbol).Info("❌ Stock disqualified due to being below 100MA")
		return true
	}

	momentumScore := utils.CalculateMomentumScore(stockResp.DataPoints, 90)
	if momentumScore < 0 {
		log.WithField("symbol", stockResp.Symbol).Info("❌ Stock disqualified due to negative momentum score")
		return true
	}

	return false
}

func (s *MomentumStrategy) sortAndFilterSignals(signals []*pb.StockSignal) []*pb.StockSignal {
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].MomentumScore > signals[j].MomentumScore
	})

	topCount := int(float64(len(signals)) * s.topPercentage)
	return signals[:topCount]
}

func (s *MomentumStrategy) CalculateRisk(stockData *datapb.StockResponse) float64 {
	atr := utils.CalculateATR(stockData.DataPoints, 20)
	return utils.CalculateRiskUnit(atr, s.riskFactor)
}

func (s *MomentumStrategy) DetectMarketRegime(indexData *datapb.StockResponse) MarketRegime {
	if len(indexData.DataPoints) < 200 {
		log.Warnf("❗ Not enough data points to detect market regime, expected 200 but got %d", len(indexData.DataPoints))
		return Neutral
	}

	currentPrice := indexData.DataPoints[len(indexData.DataPoints)-1].Close
	ma200 := utils.CalculateMovingAverage(indexData.DataPoints, 200)

	if currentPrice > ma200 {
		log.Infof("📈 Market regime: Bull, current price: %.2f, 200MA: %.2f", currentPrice, ma200)
		return Bull
	} else {
		log.Infof("📉 Market regime: Bear, current price: %.2f, 200MA: %.2f", currentPrice, ma200)
		return Bear
	}
}

func (s *MomentumStrategy) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"lookbackPeriod": s.lookbackPeriod,
		"topPercentage":  s.topPercentage,
		"riskFactor":     s.riskFactor,
		"ma200Period":    s.ma200Period,
	}
}

func (s *MomentumStrategy) SetParameters(params map[string]interface{}) error {
	if lookbackPeriod, ok := params["lookbackPeriod"].(int); ok {
		s.lookbackPeriod = lookbackPeriod
	}
	if topPercentage, ok := params["topPercentage"].(float64); ok {
		s.topPercentage = topPercentage
	}
	if riskFactor, ok := params["riskFactor"].(float64); ok {
		s.riskFactor = riskFactor
	}
	if ma200Period, ok := params["ma200Period"].(int); ok {
		s.ma200Period = ma200Period
	}
	return nil
}
