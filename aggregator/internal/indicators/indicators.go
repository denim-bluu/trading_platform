package indicators

import (
	"math"
)

type IndicatorCalculator interface {
	CalculateMovingAverage(data []float64, period int) []float64
	CalculateATR(high, low, close []float64, period int) []float64
}

type indicatorCalculator struct{}

func NewIndicatorCalculator() IndicatorCalculator {
	return &indicatorCalculator{}
}

func (ic *indicatorCalculator) CalculateMovingAverage(data []float64, period int) []float64 {
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

func (ic *indicatorCalculator) CalculateATR(high, low, close []float64, period int) []float64 {
	if period <= 0 || len(high) < period || len(low) < period || len(close) < period {
		return nil
	}
	tr := make([]float64, len(high))
	for i := range high {
		if i == 0 {
			tr[i] = high[i] - low[i]
		} else {
			tr[i] = math.Max(high[i]-low[i], math.Max(math.Abs(high[i]-close[i-1]), math.Abs(low[i]-close[i-1])))
		}
	}
	atr := make([]float64, len(high)-period+1)
	atr[0] = average(tr[:period])
	for i := period; i < len(high); i++ {
		atr[i-period+1] = (atr[i-period]*float64(period-1) + tr[i]) / float64(period)
	}
	return atr
}

func average(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}
