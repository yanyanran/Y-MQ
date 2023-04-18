package yerMQ

import "sync"

type yerMQ struct {
	sync.Mutex
	ctx context.Context
}