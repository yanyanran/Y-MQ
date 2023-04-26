package yerMQ

import (
	"net"
	"sync"
)

type tcpServer struct {
	yerMQ *YERMQ
	conns sync.Map
}

func (t *tcpServer) Handle(conn net.Conn) {

}

func (t *tcpServer) Close() {

}
