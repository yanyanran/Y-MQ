package yerMQ

import (
	"Y-MQ/inIt/protocol"
	"log"
	"net"
	"sync"
)

const (
	typeConsumer = iota
	typeProducer
)

type tcpServer struct {
	yerMQ *YERMQ
	conns sync.Map
}

// Handle implement interface function Handle
func (t *tcpServer) Handle(conn net.Conn) {
	log.Printf("TCP: new client %s", conn.RemoteAddr())
	var prot protocol.Protocol

	client := prot.NewClient(conn)           // new connect object [IOLoop]
	t.conns.Store(conn.RemoteAddr(), client) // put client into map where key is remoteAddr

	err := prot.IOLoop(client)
	if err != nil {
		log.Printf("[ERROR] client(%s) - %s", conn.RemoteAddr(), err)
	}

	t.conns.Delete(conn.RemoteAddr())
	client.Close()
}

func (t *tcpServer) Close() {
	log.Printf("[TCPServer] Client Closing...")
	t.conns.Range(func(k, v any) bool {
		v.(protocol.Client).Close()
		return true
	})
}
