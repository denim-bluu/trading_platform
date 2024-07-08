package strategies

import (
	"sort"
	"time"
	"trading_platform/strategy/internal/indicators"
	pb "trading_platform/strategy/proto"
)

type MomentumStrategy struct{}

func NewMomentumStrategy() *MomentumStrategy {
	return &MomentumStrategy{}
}

func (s *MomentumStrategy) Evaluate(data []*pb.DataPoint, startDate, endDate string, indexValue float64) []*pb.TradeAction {
	var actions []*pb.TradeAction
	var inPosition bool
	var entryPrice float64

	// Extract close prices for indicator calculations
	closePrices := make([]float64, len(data))
	highPrices := make([]float64, len(data))
	lowPrices := make([]float64, len(data))
	for i, dp := range data {
		closePrices[i] = dp.Adjclose
		highPrices[i] = dp.High
		lowPrices[i] = dp.Low
	}

	// Calculate indicators
	ma20 := indicators.CalculateMovingAverage(closePrices, 20)
	ma100 := indicators.CalculateMovingAverage(closePrices, 100)
	ma200 := indicators.CalculateMovingAverage(closePrices, 200)
	atr20 := indicators.CalculateATR(highPrices, lowPrices, closePrices, 20)
	adjRegressionSlope := indicators.CalculateAdjustedRegressionSlope(closePrices, 90)
	maxGap := indicators.CalculateMaxGap(closePrices, 90)

	// Add indicators to data points
	for i := range data {
		if data[i].Indicators == nil {
			data[i].Indicators = make(map[string]float64)
		}

		if i >= 19 {
			data[i].Indicators["MA20"] = ma20[i-19]
		}
		if i >= 99 {
			data[i].Indicators["MA100"] = ma100[i-99]
		}
		if i >= 199 {
			data[i].Indicators["MA200"] = ma200[i-199]
		}
		if i >= 19 {
			data[i].Indicators["ATR20"] = atr20[i-19]
		}
		if i >= 89 {
			data[i].Indicators["AdjustedRegressionSlope"] = adjRegressionSlope
		}
		if i >= 89 {
			data[i].Indicators["MaxGap"] = maxGap
		}
	}

	// Filter data to start from the requested start date
	layout := "2006-01-02"
	start, _ := time.Parse(layout, startDate)
	end, _ := time.Parse(layout, endDate)
	var filteredData []*pb.DataPoint
	for _, dp := range data {
		dpTime, _ := time.Parse("2006-01-02 15:04:05", dp.Timestamp)
		if dpTime.Before(start) || dpTime.After(end) {
			continue
		}
		filteredData = append(filteredData, dp)
	}

	// Step 1: Sort by adjusted regression slope (volatility adjusted momentum)
	sort.Slice(filteredData, func(i, j int) bool {
		return filteredData[i].Indicators["AdjustedRegressionSlope"] > filteredData[j].Indicators["AdjustedRegressionSlope"]
	})

	// Step 2: Filter out stocks below 100-day MA or with a recent gap > 15%
	validData := []*pb.DataPoint{}
	for _, dp := range filteredData {
		if dp.Close >= dp.Indicators["MA100"] && dp.Indicators["MaxGap"] <= 0.15 {
			validData = append(validData, dp)
		}
	}

	// Step 3: Check sell conditions and rebalance positions
	for i := 0; i < len(validData); i++ {
		if i < 200 {
			continue // Need at least 200 days of data for the strategy
		}

		ma200 := validData[i].Indicators["MA200"]
		ma100 := validData[i].Indicators["MA100"]
		atr20 := validData[i].Indicators["ATR20"]
		price := validData[i].Close

		if !inPosition && price > ma200 && validData[i-1].Close <= validData[i-1].Indicators["MA200"] && indexValue > ma200 {
			// Enter position
			inPosition = true
			entryPrice = price
			actions = append(actions, &pb.TradeAction{
				Symbol:    validData[i].Symbol,
				Action:    "BUY",
				Price:     price,
				Timestamp: validData[i].Timestamp,
			})
		} else if inPosition && (price < ma200 || price < entryPrice-2*atr20) {
			// Exit position
			inPosition = false
			actions = append(actions, &pb.TradeAction{
				Symbol:    validData[i].Symbol,
				Action:    "SELL",
				Price:     price,
				Timestamp: validData[i].Timestamp,
			})
		} else if inPosition && price > ma100 {
			// Add to position
			actions = append(actions, &pb.TradeAction{
				Symbol:    validData[i].Symbol,
				Action:    "ADD",
				Price:     price,
				Timestamp: validData[i].Timestamp,
			})
		}
	}

	return actions
}
