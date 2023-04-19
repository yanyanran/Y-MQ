package yerMQ

import "sync"

// TODO:加入持久化后台队列

type Channel struct {
	mu         sync.RWMutex
	requeueNum uint64
	msgNum     uint64
	timeoutNum uint64

	name      string
	topicName string
}
