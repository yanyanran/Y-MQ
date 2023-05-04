package yerMQ

import (
	"sync"
	"sync/atomic"
)

type Consumer interface { // TODO???
	UnStop()
	Stop()
	Close() error
	TimedOutMessage()
	Stats(string) ClientStats
	Empty()
}

// TODO:加入持久化后台队列

type Channel struct {
	mu         sync.RWMutex
	requeueNum uint64
	msgNum     uint64
	timeoutNum uint64

	name      string
	topicName string
	yerMq     *YERMQ

	clients map[int64]Consumer
	stopped int32
}

func NewChannel() *Channel {
	// TODO
}

func (c *Channel) Stop() error {
	return c.doStop(true)
}

func (c *Channel) UnStop() error {
	return c.doStop(false)
}

func (c *Channel) doStop(pause bool) error {
	if pause {
		atomic.StoreInt32(&c.stopped, 1)
	} else {
		atomic.StoreInt32(&c.stopped, 0)
	}

	c.mu.RLock()
	for _, client := range c.clients {
		if pause {
			client.Stop() // TODO???
		} else {
			client.UnStop()
		}
	}
	c.mu.RUnlock()
	return nil
}

func (c *Channel) IsStopped() bool {
	return atomic.LoadInt32(&c.stopped) == 1
}
