// internal/storage/csv_storage.go

package storage

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	pb "momentum-trading-platform/api/proto/portfolio_service"
)

type CSVStorage struct {
	portfolioFile string
	tradeFile     string
	mu            sync.RWMutex
}

func NewCSVStorage(portfolioFile, tradeFile string) *CSVStorage {
	return &CSVStorage{
		portfolioFile: portfolioFile,
		tradeFile:     tradeFile,
	}
}

func (s *CSVStorage) SavePortfolioState(ctx context.Context, state *pb.PortfolioStatus, isSnapshot bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.OpenFile(s.portfolioFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open portfolio file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	timestamp := time.Now().Format(time.RFC3339)
	totalValue := state.CashBalance
	for _, position := range state.Positions {
		totalValue += position.MarketValue
		record := []string{
			timestamp,
			strconv.FormatBool(isSnapshot),
			position.Symbol,
			strconv.FormatInt(int64(position.Quantity), 10),
			strconv.FormatFloat(position.AveragePrice, 'f', 2, 64),
			strconv.FormatFloat(position.CurrentPrice, 'f', 2, 64),
			strconv.FormatFloat(position.MarketValue, 'f', 2, 64),
			strconv.FormatFloat(state.CashBalance, 'f', 2, 64),
			strconv.FormatFloat(totalValue, 'f', 2, 64),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write portfolio state: %v", err)
		}
	}

	return nil
}

func (s *CSVStorage) LoadPortfolioState(ctx context.Context) (*pb.PortfolioStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Open(s.portfolioFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open portfolio file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read portfolio file: %v", err)
	}

	if len(records) == 0 {
		return &pb.PortfolioStatus{}, nil
	}

	// Get the latest snapshot
	var latestSnapshot []string
	for i := len(records) - 1; i >= 0; i-- {
		if records[i][1] == "true" { // IsSnapshot
			latestSnapshot = records[i]
			break
		}
	}

	if latestSnapshot == nil {
		return &pb.PortfolioStatus{}, nil
	}

	state := &pb.PortfolioStatus{}
	state.CashBalance, _ = strconv.ParseFloat(latestSnapshot[7], 64)
	state.LastUpdateDate = latestSnapshot[0]

	// Reconstruct positions from the snapshot
	for _, record := range records {
		if record[0] == latestSnapshot[0] && record[1] == "true" {
			quantity, _ := strconv.ParseInt(record[3], 10, 32)
			averagePrice, _ := strconv.ParseFloat(record[4], 64)
			currentPrice, _ := strconv.ParseFloat(record[5], 64)
			marketValue, _ := strconv.ParseFloat(record[6], 64)
			position := &pb.Position{
				Symbol:       record[2],
				Quantity:     int32(quantity),
				AveragePrice: averagePrice,
				CurrentPrice: currentPrice,
				MarketValue:  marketValue,
			}
			state.Positions = append(state.Positions, position)
		}
	}

	return state, nil
}

func (s *CSVStorage) SaveTrade(ctx context.Context, trade *pb.Trade, cashBalanceAfter float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.OpenFile(s.tradeFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open trade file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		time.Now().Format(time.RFC3339),
		trade.Symbol,
		trade.Type.String(),
		strconv.FormatInt(int64(trade.Quantity), 10),
		strconv.FormatFloat(trade.Price, 'f', 2, 64),
		strconv.FormatFloat(cashBalanceAfter, 'f', 2, 64),
	}

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write trade: %v", err)
	}

	return nil
}

func (s *CSVStorage) GetTradeHistory(ctx context.Context, startDate, endDate string) ([]*pb.Trade, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Open(s.tradeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open trade file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read trade file: %v", err)
	}

	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)

	var trades []*pb.Trade
	for _, record := range records {
		tradeTime, _ := time.Parse(time.RFC3339, record[0])
		if tradeTime.After(start) && tradeTime.Before(end) {
			quantity, _ := strconv.ParseInt(record[3], 10, 32)
			price, _ := strconv.ParseFloat(record[4], 64)
			trade := &pb.Trade{
				Symbol:   record[1],
				Type:     pb.TradeType(pb.TradeType_value[record[2]]),
				Quantity: int32(quantity),
				Price:    price,
			}
			trades = append(trades, trade)
		}
	}

	return trades, nil
}
