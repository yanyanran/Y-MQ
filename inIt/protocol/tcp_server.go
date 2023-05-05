package protocol

import (
	"log"
	"net"
	"sync"
)

// TCPHandler 任何实现了【Handle方法】的类都可以作为TCPServer的第二个参
type TCPHandler interface {
	Handle(conn net.Conn)
}

func TCPServer(listener net.Listener, handler TCPHandler) error { // 启动监听
	log.Printf("[TCPServer LOG] TCP Listening in %s\n", listener.Addr())
	var wg sync.WaitGroup

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			// TODO:临时性错误处理-> break/continue
			return err
		}

		wg.Add(1)
		go func() {
			handler.Handle(clientConn)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Printf("TCP: closing %s", listener.Addr())

	return nil
}
