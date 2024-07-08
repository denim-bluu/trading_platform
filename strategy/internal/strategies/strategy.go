package strategies

import pb "trading_platform/strategy/proto"

type TradingStrategy interface {
	Evaluate(data []*pb.DataPoint, startDate, endDate string, indexValue float64) []*pb.TradeAction
}
