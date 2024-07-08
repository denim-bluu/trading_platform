package indicators

import (
	"math"
)

func CalculateAdjustedRegressionSlope(data []float64, days int) float64 {
	n := float64(days)
	sumX := n * (n - 1) / 2
	sumX2 := n * (n - 1) * (2*n - 1) / 6

	sumY := 0.0
	sumXY := 0.0
	for i := 0; i < days; i++ {
		sumY += data[i]
		sumXY += float64(i) * data[i]
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	return slope
}

func CalculateATR(high, low, close []float64, period int) []float64 {
	atr := make([]float64, len(close))
	for i := period; i < len(close); i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			tr := math.Max(high[i-j], close[i-j-1]) - math.Min(low[i-j], close[i-j-1])
			sum += tr
		}
		atr[i] = sum / float64(period)
	}
	return atr
}

func CalculateMaxGap(data []float64, days int) float64 {
	maxGap := 0.0
	for i := 1; i < days; i++ {
		gap := math.Abs(data[i]-data[i-1]) / data[i-1]
		if gap > maxGap {
			maxGap = gap
		}
	}
	return maxGap
}

func CalculateMovingAverage(data []float64, period int) []float64 {
	if period <= 0 || len(data) < period {
		return nil
	}
	result := make([]float64, len(data)-period+1)
	for i := range result {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += data[j]
		}
		result[i] = sum / float64(period)
	}
	return result
}
