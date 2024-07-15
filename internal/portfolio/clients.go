package portfolio

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	datapb "momentum-trading-platform/api/proto/data_service"
	strategypb "momentum-trading-platform/api/proto/strategy_service"
)

type Clients struct {
	DataClient     datapb.DataServiceClient
	StrategyClient strategypb.StrategyServiceClient
	connections    []*grpc.ClientConn
}

func NewClients() (*Clients, error) {
	dataConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to data service: %v", err)
	}

	strategyConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to strategy service: %v", err)
	}

	return &Clients{
		DataClient:     datapb.NewDataServiceClient(dataConn),
		StrategyClient: strategypb.NewStrategyServiceClient(strategyConn),
		connections:    []*grpc.ClientConn{dataConn, strategyConn},
	}, nil
}

func (c *Clients) Close() {
	for _, conn := range c.connections {
		conn.Close()
	}
}
