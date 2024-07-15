package data

import (
	"google.golang.org/grpc"
)

type Clients struct {
	connections []*grpc.ClientConn
}

func NewClients() (*Clients, error) {
	return &Clients{
		connections: []*grpc.ClientConn{},
	}, nil
}

func (c *Clients) Close() {
	for _, conn := range c.connections {
		conn.Close()
	}
}
