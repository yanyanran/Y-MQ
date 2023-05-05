package yerMQ

import (
	"bufio"
	"net"
	"sync"
	"time"
)

const defaultBufferSize = 16 * 1024 // Default buffer

// client's state
const (
	stateInit = iota
	stateDisconnected
	stateConnected
	stateSubscribed // SUB
	stateClosing
)

type client struct {
	yerMq *YERMQ

	ID      int64
	State   int32
	Channel *Channel
	net.Conn

	Reader    *bufio.Reader
	Writer    *bufio.Writer
	writeLock sync.RWMutex
	metaLock  sync.RWMutex

	FlushBufferSize    int
	FlushBufferTimeout time.Duration // flush to remote
	MsgTimeout         time.Duration
	WhenHeartbeat      time.Duration

	ReadStateChan chan int // [client's state change] redetermine the client receive status
	ExitChan      chan int
	SubEventChan  chan *Channel

	// admin
	ReadyCount   int64
	InFightCount int64
	MsgCount     uint64
	FinishCount  uint64
	RequeueCount uint64
	pubCounts    map[string]uint64

	// can reuse buffer
	lenBuf   [4]byte
	lenSlice []byte
}

func newClient(id int64, conn net.Conn, yerMq *YERMQ) *client {
	c := &client{
		yerMq:              yerMq,
		ID:                 id,
		Conn:               conn,
		Reader:             bufio.NewReaderSize(conn, defaultBufferSize),
		Writer:             bufio.NewWriterSize(conn, defaultBufferSize),
		FlushBufferSize:    defaultBufferSize,
		FlushBufferTimeout: yerMq.getOpts().FlushBufferTimeout,
		MsgTimeout:         yerMq.getOpts().MsgTimeout,
		WhenHeartbeat:      yerMq.getOpts().ClientTimeout / 2,
		ReadStateChan:      make(chan int, 1),
		ExitChan:           make(chan int),
		SubEventChan:       make(chan *Channel, 1),
		pubCounts:          make(map[string]uint64),
	}
	c.lenSlice = c.lenBuf[:]
	return c
}

func (c *client) Type() int {
	c.metaLock.RLock()
	has := len(c.pubCounts) > 0
	c.metaLock.RUnlock()
	if has {
		return typeProducer
	}
	return typeConsumer
}

func (c *client) Flush() error {
	var zeroTime time.Time
	if c.WhenHeartbeat > 0 {
		c.SetWriteDeadline(time.Now().Add(c.WhenHeartbeat))
	} else {
		c.SetWriteDeadline(zeroTime)
	}

	err := c.Writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
