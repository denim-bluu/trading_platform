package strategy

import (
	datapb "momentum-trading-platform/api/proto/data_service"
	pb "momentum-trading-platform/api/proto/strategy_service"
)

type MarketRegime int

const (
	Bull MarketRegime = iota
	Bear
	Neutral
)

type Strategy interface {
	GenerateSignals(stockData map[string]*datapb.StockResponse, indexData *datapb.StockResponse) ([]*pb.StockSignal, error)
	CalculateRisk(stockData *datapb.StockResponse) float64
	GetParameters() map[string]interface{}
	SetParameters(params map[string]interface{}) error
	DetectMarketRegime(marketIndexData *datapb.StockResponse) MarketRegime
}
