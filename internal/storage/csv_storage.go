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

	log "github.com/sirupsen/logrus"

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

func (s *CSVStorage) SavePortfolioState(ctx context.Context, state *pb.PortfolioStatus) error {
	log.Infof("Saving portfolio state: %v", state)
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create(s.portfolioFile)
	if err != nil {
		return fmt.Errorf("failed to create portfolio file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Symbol", "Quantity", "AveragePrice", "CurrentPrice", "MarketValue"}); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}

	// Write positions
	for _, position := range state.Positions {
		record := []string{
			position.Symbol,
			strconv.FormatInt(int64(position.Quantity), 10),
			strconv.FormatFloat(position.AveragePrice, 'f', 2, 64),
			strconv.FormatFloat(position.CurrentPrice, 'f', 2, 64),
			strconv.FormatFloat(position.MarketValue, 'f', 2, 64),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write position: %v", err)
		}
	}

	// Write cash balance and last update date
	if err := writer.Write([]string{"CashBalance", strconv.FormatFloat(state.CashBalance, 'f', 2, 64)}); err != nil {
		return fmt.Errorf("failed to write cash balance: %v", err)
	}
	if err := writer.Write([]string{"LastUpdateDate", state.LastUpdateDate}); err != nil {
		return fmt.Errorf("failed to write last update date: %v", err)
	}

	return nil
}

func (s *CSVStorage) LoadPortfolioState(ctx context.Context) (*pb.PortfolioStatus, error) {
	log.Info("Loading portfolio state")
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

	state := &pb.PortfolioStatus{}
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) == 2 {
			// Cash balance or last update date
			switch record[0] {
			case "CashBalance":
				state.CashBalance, _ = strconv.ParseFloat(record[1], 64)
			case "LastUpdateDate":
				state.LastUpdateDate = record[1]
			}
		} else if len(record) == 5 {
			// Position
			quantity, _ := strconv.ParseInt(record[1], 10, 32)
			averagePrice, _ := strconv.ParseFloat(record[2], 64)
			currentPrice, _ := strconv.ParseFloat(record[3], 64)
			marketValue, _ := strconv.ParseFloat(record[4], 64)
			position := &pb.Position{
				Symbol:       record[0],
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

func (s *CSVStorage) SaveTrade(ctx context.Context, trade *pb.Trade) error {
	log.Infof("Saving trade: %v", trade)
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
	}

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write trade: %v", err)
	}

	return nil
}

func (s *CSVStorage) GetTradeHistory(ctx context.Context, startDate, endDate string) ([]*pb.Trade, error) {
	log.Infof("Getting trade history between %s and %s", startDate, endDate)
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
