package yerMQ

import (
	"Y-MQ/inIt/protocol"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"
)

// Type
const (
	TypeResponse int32 = 0
	TypeError    int32 = 1
	TypeMessage  int32 = 2
)

var separatorBytes = []byte(" ")
var heartbeatBytes = []byte("_heartbeat_")
var okBytes = []byte("OK")

type yProtocol struct {
	yerMq *YERMQ
}

func (p *yProtocol) NewClient(conn net.Conn) protocol.Client {
	clientID := atomic.AddInt64(&p.yerMq.clientIDIndex, 1)
	return newClient(clientID, conn, p.yerMq)
}

// IOLoop dispose tcp conn
func (p *yProtocol) IOLoop(c protocol.Client) error {
	var (
		err      error
		line     []byte
		zeroTime time.Time
	)

	client := c.(*client)

	msgPumpStartedChan := make(chan bool)
	// one conn have one msgPump goroutine, communicate by chan
	go p.msgPump(client, msgPumpStartedChan)
	// init over can be start
	<-msgPumpStartedChan

	for {
		if client.WhenHeartbeat > 0 {
			client.SetReadDeadline(time.Now().Add(client.WhenHeartbeat * 2))
		} else {
			client.SetReadDeadline(zeroTime)
		}

		// read util encounter '\n'
		line, err = client.Reader.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("failed to read command - %s", err)
			}
			break
		}

		// trim
		line = line[:len(line)-1]
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		// parse the command
		params := bytes.Split(line, separatorBytes)

		log.Printf("PROTOCOL: [%s] %s", client, params)

		var response []byte
		// assign instruction correspondence methods
		response, err = p.Execute(client, params)
		if err != nil {
			// send error and determine error's type
			ctx := ""
			if parentErr := err.(protocol.ChildErr).Parent(); parentErr != nil {
				ctx = " - " + parentErr.Error()
			}
			log.Printf("[PROTOCOL ERROR][%s] - %s%s", client, err, ctx)

			sendErr := p.Send(client, TypeError, []byte(err.Error()))
			// send failed
			if sendErr != nil {
				log.Printf("[PROTOCOL ERROR][%s] - %s%s", client, sendErr, ctx)
				break
			}

			// error is fatal
			if _, ok := err.(*protocol.FatalClientErr); ok {
				break
			}
			continue
		}

		if response != nil {
			err = p.Send(client, TypeResponse, response)
			if err != nil {
				err = fmt.Errorf("failed to send response - %s", err)
				break
			}
		}
	}

	log.Printf("PROTOCOL: [%s] exiting ioloop", client)
	close(client.ExitChan)
	if client.Channel != nil {
		client.Channel.RemoveClient(client.ID)
	}

	return err
}

func (p *yProtocol) msgPump(c *client, startedChan chan bool) {
	// todo 内部监听多个事件 【消息接收，消息投递，订阅channel，发送心跳包】
}

// Execute tcp command
func (p *yProtocol) Execute(c *client, params [][]byte) ([]byte, error) {
	switch {
	case bytes.Equal(params[0], []byte("FIN")):
		return p.FIN(c, params)
		// TODO RDY
		// TODO SUB
		// TODO REQ
		// TODO PUB
		// TODO NPUB
		// TODO DPUB
		// TODO NOP
	}
}

func (p *yProtocol) Send(c *client, ty int32, data []byte) error {
	c.writeLock.Lock()
	var zeroTime time.Time
	if c.WhenHeartbeat > 0 {
		c.SetReadDeadline(time.Now().Add(c.WhenHeartbeat))
	} else {
		c.SetReadDeadline(zeroTime)
	}

	_, err := protocol.SendTypeResponse(c.Writer, ty, data)
	if err != nil {
		c.writeLock.Unlock()
		return err
	}

	if ty != TypeMessage {
		err = c.Flush()
	}

	c.writeLock.Unlock()
	return err
}

func (p *yProtocol) FIN(c *client, param [][]byte) ([]byte, error) {
	//TODO
}

func (p *yProtocol) RDY(c *client, param [][]byte) ([]byte, error) {
	//TODO
}

func (p *yProtocol) SUB(c *client, param [][]byte) ([]byte, error) {
	//TODO
}

func (p *yProtocol) REQ(c *client, param [][]byte) ([]byte, error) {
	//TODO
}

func (p *yProtocol) PUB(c *client, param [][]byte) ([]byte, error) {
	//TODO
}

func (p *yProtocol) NPUB(c *client, param [][]byte) ([]byte, error) {
	//TODO
}

func (p *yProtocol) DPUB(c *client, param [][]byte) ([]byte, error) {
	//TODO
}

func (p *yProtocol) NOP(c *client, param [][]byte) ([]byte, error) {
	//TODO
}
