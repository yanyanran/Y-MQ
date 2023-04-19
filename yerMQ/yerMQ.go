package yerMQ

import (
	"context"
	"sync"
)

type yerMQ struct {
	sync.Mutex
	ctx context.Context
}
