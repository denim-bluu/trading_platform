package data

import (
	"context"
	"fmt"
	"time"

	pb "momentum-trading-platform/api/proto/data_service"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS stock_data (
    symbol VARCHAR(10),
    timestamp BIGINT,
    open DECIMAL(10,2),
    high DECIMAL(10,2),
    low DECIMAL(10,2),
    close DECIMAL(10,2),
    adjusted_close DECIMAL(10,2),
    volume BIGINT,
    PRIMARY KEY (symbol, timestamp)
);`

// Add this function
func (s *Server) initDatabase() error {
	_, err := s.DB.Exec(createTableSQL)
	return err
}

func (s *Server) getStockDataFromDB(symbol, startDate, endDate string) (*pb.StockResponse, error) {
	query := `SELECT timestamp, open, high, low, close, adjusted_close, volume 
              FROM stock_data 
              WHERE symbol = $1 AND timestamp BETWEEN $2 AND $3 
              ORDER BY timestamp`

	startTimestamp, _ := time.Parse("2006-01-02", startDate)
	endTimestamp, _ := time.Parse("2006-01-02", endDate)

	rows, err := s.DB.Query(query, symbol, startTimestamp.Unix(), endTimestamp.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataPoints []*pb.StockDataPoint
	for rows.Next() {
		var dp pb.StockDataPoint
		err := rows.Scan(&dp.Timestamp, &dp.Open, &dp.High, &dp.Low, &dp.Close, &dp.AdjustedClose, &dp.Volume)
		if err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, &dp)
	}

	if len(dataPoints) == 0 {
		return nil, fmt.Errorf("no data found for symbol %s", symbol)
	}

	return &pb.StockResponse{
		Symbol:     symbol,
		DataPoints: dataPoints,
	}, nil
}

func (s *Server) storeStockDataInDB(data *pb.StockResponse) error {
	query := `INSERT INTO stock_data (symbol, timestamp, open, high, low, close, adjusted_close, volume) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
              ON CONFLICT (symbol, timestamp) DO UPDATE 
              SET open = $3, high = $4, low = $5, close = $6, adjusted_close = $7, volume = $8`

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	for _, dp := range data.DataPoints {
		_, err := tx.Exec(query, data.Symbol, dp.Timestamp, dp.Open, dp.High, dp.Low, dp.Close, dp.AdjustedClose, dp.Volume)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Server) updateDatabaseWithLatestData(symbols []string) error {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -252)

	for _, symbol := range symbols {
		req := &pb.StockRequest{
			Symbol:    symbol,
			StartDate: startDate.Format("2006-01-02"),
			EndDate:   endDate.Format("2006-01-02"),
			Interval:  "1d",
		}

		data, err := s.fetchStockData(context.Background(), req.Symbol, req.StartDate, req.EndDate, req.Interval)
		if err != nil {
			s.Logger.WithError(err).Errorf("Failed to fetch latest data for %s", symbol)
			return fmt.Errorf("failed to fetch latest data for %s: %v", symbol, err)
		}

		err = s.storeStockDataInDB(data)
		if err != nil {
			s.Logger.WithError(err).Errorf("Failed to store latest data for %s", symbol)
			return fmt.Errorf("failed to store latest data for %s: %v", symbol, err)
		}
	}

	return nil
}
