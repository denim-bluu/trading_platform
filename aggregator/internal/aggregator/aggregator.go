package aggregator

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"trading_platform/aggregator/internal/data"
	pb "trading_platform/aggregator/proto"
)

type Aggregator struct {
	fetcher data.DataFetcher
	storage data.DataStorage
	dataDir string
	cache   *data.Cache
}

func parseTimestamp(timestamp string) (time.Time, error) {
	dpTime, err := time.Parse("2006-01-02", timestamp)
	if err != nil {
		dpTime, err = time.Parse("2006-01-02 15:04:05", timestamp)
	}
	return dpTime, err
}

func NewAggregator(fetcher data.DataFetcher, storage data.DataStorage, dataDir string, cacheTTL time.Duration) *Aggregator {
	os.MkdirAll(dataDir, os.ModePerm) // Ensure the data directory exists
	return &Aggregator{fetcher: fetcher, storage: storage, dataDir: dataDir, cache: data.NewCache(cacheTTL)}
}

func (a *Aggregator) getFilename(symbol string) string {
	return filepath.Join(a.dataDir, symbol+".pb")
}

func (a *Aggregator) AggregateHistoricalData(symbol string, startDate string, endDate string) error {
	log.Printf("Fetching historical data for symbol: %s", symbol)
	data, err := a.fetcher.FetchHistoricalData(symbol, startDate, endDate)
	if err != nil {
		log.Printf("Failed to fetch historical data: %v", err)
		return err
	}
	log.Printf("Fetched %d data points for symbol: %s", len(data), symbol)

	filename := a.getFilename(symbol)
	log.Printf("Saving aggregated data to file: %s", filename)
	protoData := &pb.HistoricalData{DataPoints: data}
	err = a.storage.SaveData(filename, protoData)
	if err != nil {
		log.Printf("Failed to save data: %v", err)
		return err
	}

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	a.cache.Set(symbol, start, end, data)

	log.Println("Successfully saved aggregated data")
	return nil
}

func (a *Aggregator) UpdateLiveData(symbol string) error {
	log.Printf("Fetching live data for symbol: %s", symbol)
	liveData, err := a.fetcher.FetchLiveData(symbol)
	if err != nil {
		log.Printf("Failed to fetch live data: %v", err)
		return err
	}

	filename := a.getFilename(symbol)
	log.Printf("Loading existing data from file: %s", filename)
	storedData, err := a.storage.LoadData(filename)
	if err != nil {
		log.Printf("Failed to load data: %v", err)
		return err
	}

	log.Println("Updating stored data with live data")
	storedData.DataPoints = append(storedData.DataPoints, liveData)
	err = a.storage.SaveData(filename, storedData)
	if err != nil {
		log.Printf("Failed to save updated data: %v", err)
		return err
	}
	start, err := parseTimestamp(storedData.DataPoints[0].Timestamp)
	if err != nil {
		log.Printf("Failed to parse start timestamp: %v", err)
		return err
	}
	end, err := parseTimestamp(storedData.DataPoints[len(storedData.DataPoints)-1].Timestamp)
	if err != nil {
		log.Printf("Failed to parse end timestamp: %v", err)
		return err
	}

	a.cache.Set(symbol, start, end, storedData.DataPoints)
	log.Println("Successfully updated live data")
	return nil
}

func (a *Aggregator) GetHistoricalData(symbol string, startDate string, endDate string) ([]*pb.DataPoint, error) {
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

	if cachedData, found := a.cache.Get(symbol, start, end); found {
		log.Printf("Cache hit for symbol: %s", symbol)
		return cachedData, nil
	}
	filename := a.getFilename(symbol)
	log.Printf("Loading data from file: %s", filename)
	storedData, err := a.storage.LoadData(filename)
	if err != nil {
		log.Printf("Failed to load data: %v", err)
		return nil, err
	}

	var filteredData []*pb.DataPoint
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

	a.cache.Set(symbol, start, end, filteredData)

	log.Printf("Returning %d filtered data points", len(filteredData))
	return filteredData, nil
}
