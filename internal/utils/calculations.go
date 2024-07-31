package utils

import (
	"math"
	"time"

	datapb "momentum-trading-platform/api/proto/data_service"
	stratpb "momentum-trading-platform/api/proto/strategy_service"

	"gonum.org/v1/gonum/stat"
)

func IsWednesday(date string) bool {
	t, _ := time.Parse("2006-01-02", date)
	return t.Weekday() == time.Wednesday
}

func IsSecondWednesdayOfMonth(date string) bool {
	t, _ := time.Parse("2006-01-02", date)
	return t.Weekday() == time.Wednesday && (t.Day()-1)/7 == 1
}

func CalculateMomentumScore(dataPoints []*datapb.StockDataPoint, period int) float64 {
	if len(dataPoints) < period {
		return 0
	}

	x := make([]float64, period)
	y := make([]float64, period)

	data := dataPoints[len(dataPoints)-period:]

	for i, dp := range data {
		x[i] = float64(i) + 1.0
		y[i] = math.Log(dp.AdjustedClose)
	}

	alpha, beta := stat.LinearRegression(x, y, nil, false)
	r2 := stat.RSquared(x, y, nil, alpha, beta)

	annualizedSlope := math.Exp(beta*252) - 1 // Assuming 252 trading days in a year
	return annualizedSlope * r2
}

// GenerateSignal generates a trading signal based on momentum score, price, and moving average
func GenerateSignal(momentumScore float64, price float64, movingAverage float64) stratpb.SignalType {
	if momentumScore > 0 && price > movingAverage {
		return stratpb.SignalType_BUY
	}
	return stratpb.SignalType_HOLD
}

// CalculateRiskUnit calculates the position size based on ATR and risk factor. This is later used to calculate the position size.
func CalculateRiskUnit(atr, riskFactor float64) float64 {
	return riskFactor / atr
}

// CalculateATR calculates the Average True Range (ATR)
func CalculateATR(dataPoints []*datapb.StockDataPoint, period int) float64 {
	if len(dataPoints) < period {
		return 0
	}

	trueRanges := make([]float64, len(dataPoints)-1)
	for i := 1; i < len(dataPoints); i++ {
		high := dataPoints[i].High
		low := dataPoints[i].Low
		prevClose := dataPoints[i-1].Close
		trueRanges[i-1] = math.Max(high-low, math.Max(math.Abs(high-prevClose), math.Abs(low-prevClose)))
	}

	return stat.Mean(trueRanges[len(trueRanges)-period:], nil)
}

// CalculateMovingAverage calculates the moving average for a given period
func CalculateMovingAverage(dataPoints []*datapb.StockDataPoint, period int) float64 {
	if len(dataPoints) < period {
		return 0
	}

	sum := 0.0
	for i := len(dataPoints) - period; i < len(dataPoints); i++ {
		sum += dataPoints[i].Close
	}

	return sum / float64(period)
}

func HasRecentLargeGap(dataPoints []*datapb.StockDataPoint, period int, max_gap float64) bool {
	for i := 1; i < len(dataPoints) && i <= period; i++ {
		prevClose := dataPoints[i-1].Close
		currOpen := dataPoints[i].Open
		gap := math.Abs(currOpen-prevClose) / prevClose
		if gap > max_gap {
			return true
		}
	}
	return false
}
