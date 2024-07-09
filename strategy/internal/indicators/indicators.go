package indicators

import (
	"math"
)

func CalculateAdjustedRegressionSlope(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	// Log-transform the prices
	logPrices := make([]float64, period)
	for i := 0; i < period; i++ {
		logPrices[i] = math.Log(prices[i])
	}

	// Calculate linear regression on log-transformed prices
	var sumX, sumY, sumXY, sumX2 float64
	for i := 0; i < period; i++ {
		x := float64(i)
		y := logPrices[i]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	n := float64(period)
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	// Calculate R^2 (coefficient of determination)
	var ssTot, ssRes float64
	meanY := sumY / n
	for i := 0; i < period; i++ {
		y := logPrices[i]
		x := float64(i)
		estimatedY := slope*x + (sumY-slope*sumX)/n
		ssTot += (y - meanY) * (y - meanY)
		ssRes += (y - estimatedY) * (y - estimatedY)
	}
	rSquared := 1 - (ssRes / ssTot)

	return slope * rSquared
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
