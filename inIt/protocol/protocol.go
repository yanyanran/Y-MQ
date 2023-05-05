package protocol

import "net"

type Protocol interface {
	NewClient(net.Conn) Client
	IOLoop(Client) error
}

type Client interface {
	Close() error
}
