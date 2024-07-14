// internal/storage/csv_storage.go

package storage

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	pb "momentum-trading-platform/api/proto/portfolio_state_service"
)

type CSVStorage struct {
	filePath string
	mutex    sync.Mutex
}

func NewCSVStorage(filePath string) *CSVStorage {
	return &CSVStorage{filePath: filePath}
}

func (s *CSVStorage) SavePortfolioState(ctx context.Context, state *pb.PortfolioState) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	positionsJSON, err := json.Marshal(state.Positions)
	if err != nil {
		return fmt.Errorf("failed to marshal positions: %v", err)
	}

	record := []string{
		state.Date,
		string(positionsJSON),
		strconv.FormatFloat(state.CashBalance, 'f', 2, 64),
		strconv.FormatFloat(state.TotalValue, 'f', 2, 64),
	}

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %v", err)
	}

	return nil
}

func (s *CSVStorage) LoadPortfolioState(ctx context.Context, date string) (*pb.PortfolioState, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	for i := len(records) - 1; i >= 0; i-- {
		record := records[i]
		if record[0] == date {
			return s.parseRecord(record)
		}
	}

	return nil, fmt.Errorf("no portfolio state found for date: %s", date)
}

func (s *CSVStorage) GetPortfolioHistory(ctx context.Context, startDate, endDate string) ([]*pb.PortfolioState, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var history []*pb.PortfolioState
	for _, record := range records {
		if record[0] >= startDate && record[0] <= endDate {
			state, err := s.parseRecord(record)
			if err != nil {
				return nil, err
			}
			history = append(history, state)
		}
	}

	return history, nil
}

func (s *CSVStorage) parseRecord(record []string) (*pb.PortfolioState, error) {
	var positions []*pb.Position
	err := json.Unmarshal([]byte(record[1]), &positions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal positions: %v", err)
	}

	cashBalance, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cash balance: %v", err)
	}

	totalValue, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse total value: %v", err)
	}

	return &pb.PortfolioState{
		Date:        record[0],
		Positions:   positions,
		CashBalance: cashBalance,
		TotalValue:  totalValue,
	}, nil
}
