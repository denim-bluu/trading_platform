package strategies

import (
	"log"
	"sort"
	"time"
	"trading_platform/strategy/internal/indicators"
	pb "trading_platform/strategy/proto"
)

type MomentumStrategy struct{}

func NewMomentumStrategy() *MomentumStrategy {
	return &MomentumStrategy{}
}

func (s *MomentumStrategy) Evaluate(data []*pb.DataPoint, startDate, endDate string, indexData []*pb.DataPoint, accountValue float64) []*pb.TradeAction {
	log.Println("Evaluating momentum strategy...")

	// Step 1: Filter data to start from the requested start date
	filteredData := s.filterData(data, startDate, endDate)

	// Step 2: Calculate indicators
	warmupPeriod := s.calculateIndicators(filteredData)

	// Step 3: Rank and validate data
	validData := s.rankAndValidateData(filteredData[warmupPeriod:])

	// Step 4: Evaluate trades
	tradeActions := s.evaluateTrades(validData, indexData, accountValue, startDate, endDate)

	log.Println("Momentum strategy evaluation completed.")
	return tradeActions
}

func (s *MomentumStrategy) filterData(data []*pb.DataPoint, startDate, endDate string) []*pb.DataPoint {
	log.Println("Filtering data based on date range...")
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

	log.Printf("Filtered data points: %d", len(filteredData))
	return filteredData
}

const (
	MA20Period  = 20
	MA100Period = 100
	MA200Period = 200
	ATR20Period = 20
	SlopePeriod = 90
	GapPeriod   = 90
)

func (s *MomentumStrategy) calculateIndicators(data []*pb.DataPoint) int {
	log.Println("Calculating indicators...")

	closePrices := make([]float64, len(data))
	highPrices := make([]float64, len(data))
	lowPrices := make([]float64, len(data))
	for i, dp := range data {
		closePrices[i] = dp.Adjclose
		highPrices[i] = dp.High
		lowPrices[i] = dp.Low
	}

	ma20 := indicators.CalculateMovingAverage(closePrices, MA20Period)
	ma100 := indicators.CalculateMovingAverage(closePrices, MA100Period)
	ma200 := indicators.CalculateMovingAverage(closePrices, MA200Period)
	atr20 := indicators.CalculateATR(highPrices, lowPrices, closePrices, ATR20Period)
	adjRegressionSlope := indicators.CalculateAdjustedRegressionSlope(closePrices, SlopePeriod)
	maxGap := indicators.CalculateMaxGap(closePrices, GapPeriod)

	for i := range data {
		if data[i].Indicators == nil {
			data[i].Indicators = make(map[string]float64)
		}

		if i >= MA20Period-1 {
			data[i].Indicators["MA20"] = ma20[i-(MA20Period-1)]
		}
		if i >= MA100Period-1 {
			data[i].Indicators["MA100"] = ma100[i-(MA100Period-1)]
		}
		if i >= MA200Period-1 {
			data[i].Indicators["MA200"] = ma200[i-(MA200Period-1)]
		}
		if i >= ATR20Period-1 {
			data[i].Indicators["ATR20"] = atr20[i-(ATR20Period-1)]
		}
		if i >= SlopePeriod-1 {
			data[i].Indicators["AdjustedRegressionSlope"] = adjRegressionSlope
		}
		if i >= GapPeriod-1 {
			data[i].Indicators["MaxGap"] = maxGap
		}
	}

	warmupPeriod := max(MA20Period, MA100Period, MA200Period, ATR20Period, SlopePeriod, GapPeriod)
	log.Printf("Warm-up period determined: %d days", warmupPeriod)

	return warmupPeriod
}

func max(values ...int) int {
	maxValue := values[0]
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func (s *MomentumStrategy) rankAndValidateData(data []*pb.DataPoint) []*pb.DataPoint {
	log.Println("Ranking and validating data...")
	sort.Slice(data, func(i, j int) bool {
		return data[i].Indicators["AdjustedRegressionSlope"] > data[j].Indicators["AdjustedRegressionSlope"]
	})

	var validData []*pb.DataPoint
	for _, dp := range data {
		if dp.Close >= dp.Indicators["MA100"] && dp.Indicators["MaxGap"] <= 0.15 {
			validData = append(validData, dp)
		}
	}

	log.Printf("Valid data points after ranking and validation: %d", len(validData))
	return validData
}

func (s *MomentumStrategy) evaluateTrades(data []*pb.DataPoint, indexData []*pb.DataPoint, accountValue float64, startDate, endDate string) []*pb.TradeAction {
	log.Println("Evaluating trades based on strategy rules...")
	var actions []*pb.TradeAction
	var inPosition bool
	var entryPrice float64
	var positionSize float64

	layout := "2006-01-02"
	start, _ := time.Parse(layout, startDate)
	end, _ := time.Parse(layout, endDate)

	var indexInPositiveTrend bool
	for _, idx := range indexData {
		idxTime, _ := time.Parse(layout, idx.Timestamp)
		if idxTime.Equal(start) || idxTime.After(start) && idxTime.Before(end) {
			if idx.Close > idx.Indicators["MA200"] {
				indexInPositiveTrend = true
			}
		}
	}

	for i := 0; i < len(data); i++ {
		ma200 := data[i].Indicators["MA200"]
		ma100 := data[i].Indicators["MA100"]
		atr20 := data[i].Indicators["ATR20"]
		price := data[i].Close

		// Calculate position size using ATR
		positionSize = accountValue * 0.001 / atr20

		// Check entry condition (Wednesday check skipped for simplicity)
		if !inPosition && price > ma200 && data[i-1].Close <= data[i-1].Indicators["MA200"] && indexInPositiveTrend {
			// Enter position
			inPosition = true
			entryPrice = price
			actions = append(actions, &pb.TradeAction{
				Symbol:    data[i].Symbol,
				Action:    "BUY",
				Price:     price,
				Size:      positionSize,
				Timestamp: data[i].Timestamp,
			})
			log.Printf("BUY action for %s at %f on %s", data[i].Symbol, price, data[i].Timestamp)
		} else if inPosition && (price < ma200 || price < entryPrice-2*atr20) {
			// Exit position
			inPosition = false
			actions = append(actions, &pb.TradeAction{
				Symbol:    data[i].Symbol,
				Action:    "SELL",
				Price:     price,
				Size:      positionSize,
				Timestamp: data[i].Timestamp,
			})
			log.Printf("SELL action for %s at %f on %s", data[i].Symbol, price, data[i].Timestamp)
		} else if inPosition && price > ma100 {
			// Add to position
			actions = append(actions, &pb.TradeAction{
				Symbol:    data[i].Symbol,
				Action:    "ADD",
				Price:     price,
				Size:      positionSize,
				Timestamp: data[i].Timestamp,
			})
			log.Printf("ADD action for %s at %f on %s", data[i].Symbol, price, data[i].Timestamp)
		}
	}

	log.Printf("Trade actions generated: %d", len(actions))
	return actions
}
