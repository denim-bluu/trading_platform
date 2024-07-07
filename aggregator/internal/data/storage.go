package data

import (
	"os"

	pb "trading_platform/aggregator/proto"

	"google.golang.org/protobuf/proto"
)

type DataStorage interface {
	SaveData(filename string, data *pb.HistoricalData) error
	LoadData(filename string) (*pb.HistoricalData, error)
}

type dataStorage struct{}

func NewDataStorage() DataStorage {
	return &dataStorage{}
}

func (ds *dataStorage) SaveData(filename string, data *pb.HistoricalData) error {
	out, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, out, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStorage) LoadData(filename string) (*pb.HistoricalData, error) {
	in, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	data := &pb.HistoricalData{}
	err = proto.Unmarshal(in, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
