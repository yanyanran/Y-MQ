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

}
