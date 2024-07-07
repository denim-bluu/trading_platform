package aggregator

import (
	"log"
	"time"

	"trading_platform/aggregator/internal/data"
	"trading_platform/aggregator/internal/indicators"
	pb "trading_platform/aggregator/proto"
)

type Aggregator struct {
	fetcher    data.DataFetcher
	calculator indicators.IndicatorCalculator
	storage    data.DataStorage
}

func NewAggregator(fetcher data.DataFetcher, calculator indicators.IndicatorCalculator, storage data.DataStorage) *Aggregator {
	return &Aggregator{fetcher: fetcher, calculator: calculator, storage: storage}
}

func (a *Aggregator) AggregateHistoricalData(symbol string, startDate string, endDate string, filename string) error {
	log.Printf("Fetching historical data for symbol: %s", symbol)
	data, err := a.fetcher.FetchHistoricalData(symbol, startDate, endDate)
	if err != nil {
		log.Printf("Failed to fetch historical data: %v", err)
		return err
	}
	log.Printf("Fetched %d data points for symbol: %s", len(data), symbol)

	a.calculateAndStoreIndicators(data, filename)
	return nil
}

func (a *Aggregator) UpdateLiveData(symbol string, filename string) error {
	log.Printf("Fetching live data for symbol: %s", symbol)
	liveData, err := a.fetcher.FetchLiveData(symbol)
	if err != nil {
		log.Printf("Failed to fetch live data: %v", err)
		return err
	}

	log.Printf("Loading existing data from file: %s", filename)
	storedData, err := a.storage.LoadData(filename)
	if err != nil {
		log.Printf("Failed to load data: %v", err)
		return err
	}

	log.Println("Updating stored data with live data")
	storedData.DataPoints = append(storedData.DataPoints, liveData)
	a.calculateAndStoreIndicators(storedData.DataPoints, filename)
	return nil
}

func (a *Aggregator) GetHistoricalData(symbol string, filename string, startDate string, endDate string) ([]*pb.DataPoint, error) {
	log.Printf("Loading data from file: %s", filename)
	storedData, err := a.storage.LoadData(filename)
	if err != nil {
		log.Printf("Failed to load data: %v", err)
		return nil, err
	}

	var filteredData []*pb.DataPoint
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		log.Printf("Failed to parse start date: %v", err)
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		log.Printf("Failed to parse end date: %v", err)
		return nil, err
	}

	log.Printf("Filtering data points for symbol: %s from %s to %s", symbol, startDate, endDate)
	for _, dp := range storedData.DataPoints {
		if dp.Symbol != symbol {
			continue
		}
		var dpTime time.Time
		dpTime, err = time.Parse("2006-01-02 15:04:05", dp.Timestamp)
		if err != nil {
			dpTime, err = time.Parse("2006-01-02", dp.Timestamp)
			if err != nil {
				log.Printf("Failed to parse data point timestamp: %v", err)
				continue // Skip this data point if parsing fails
			}
		}
		if (dpTime.Equal(start) || dpTime.After(start)) && (dpTime.Equal(end) || dpTime.Before(end)) {
			filteredData = append(filteredData, dp)
		}
	}

	log.Printf("Returning %d filtered data points", len(filteredData))
	return filteredData, nil
}

func (a *Aggregator) calculateAndStoreIndicators(data []*pb.DataPoint, filename string) {
	log.Println("Calculating indicators")
	AdjClosePrices := make([]float64, len(data))
	highPrices := make([]float64, len(data))
	lowPrices := make([]float64, len(data))
	for i, dp := range data {
		AdjClosePrices[i] = dp.Adjclose
		highPrices[i] = dp.High
		lowPrices[i] = dp.Low
	}

	ma20 := a.calculator.CalculateMovingAverage(AdjClosePrices, 20)
	ma100 := a.calculator.CalculateMovingAverage(AdjClosePrices, 100)
	ma200 := a.calculator.CalculateMovingAverage(AdjClosePrices, 200)
	atr20 := a.calculator.CalculateATR(highPrices, lowPrices, AdjClosePrices, 20)

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
	}

	log.Printf("Saving aggregated data to file: %s", filename)
	protoData := &pb.HistoricalData{DataPoints: data}
	err := a.storage.SaveData(filename, protoData)
	if err != nil {
		log.Printf("Failed to save data: %v", err)
		return
	}

	log.Println("Successfully saved aggregated data")
}
