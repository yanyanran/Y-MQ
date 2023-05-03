package yerMQ

import (
	"sync"
	"sync/atomic"
)

// TODO: 加入后台磁盘队列

type Topic struct {
	mu sync.RWMutex

	name       string
	channelMap map[string]*Channel
	msgChan    chan *Message // 带缓存的中间管道
	startChan  chan int
	exitChan   chan int

	stopped  int32 // 用户可选择让此topic暂停服务
	stopChan chan int
}

func NewTopic(name string, inMemSize int) *Topic {
	// TODO
}

func (t *Topic) GetChannel(ch string) *Channel {
	// TODO
}

func (t *Topic) Start() {
	select {
	case t.startChan <- 1:
	default:
	}
}

func (t *Topic) Stop() error {
	return t.doStop(true)
}

func (t *Topic) UnStop() error {
	return t.doStop(false)
}

func (t *Topic) doStop(pause bool) error {
	if pause {
		atomic.StoreInt32(&t.stopped, 1)
	} else {
		atomic.StoreInt32(&t.stopped, 0)
	}

	select {
	case t.stopChan <- 1:
	case <-t.stopChan:
	}

	return nil
}

func (t *Topic) IsStopped() bool {
	return atomic.LoadInt32(&t.stopped) == 1
}
