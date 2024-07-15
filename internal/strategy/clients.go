package strategy

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	datapb "momentum-trading-platform/api/proto/data_service"
)

type Clients struct {
	DataClient  datapb.DataServiceClient
	connections []*grpc.ClientConn
}

func NewClients() (*Clients, error) {
	dataConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to data service: %v", err)
	}

	return &Clients{
		DataClient:  datapb.NewDataServiceClient(dataConn),
		connections: []*grpc.ClientConn{dataConn},
	}, nil
}

func (c *Clients) Close() {
	for _, conn := range c.connections {
		conn.Close()
	}
}
