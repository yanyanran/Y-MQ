package yerMQ

import "sync"

// TODO: 加入后台磁盘队列

type Topic struct {
	mu sync.RWMutex

	name       string
	channelMap map[string]*Channel
	msgChan    chan *Message // 带缓存的中间管道

}
