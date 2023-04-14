package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var UuidChan = make(chan []byte, 1000)

func UuidFactory(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			UuidChan <- newUuid()
		}
	}
}

func newUuid() []byte {
	// 生成节点实例
	node, err := NewWorker(1)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println(node.GetId())
	}
}

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift         = workerBits + numberBits
	workerShift       = numberBits
	startTime   int64 = 1525705533000 // 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker ID excess of quantity")
	}
	return &Worker{ // 生成一个新节点
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w *Worker) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	ID := (now-startTime)<<timeShift | (w.workerId << workerShift) | (w.number)
	return ID
}

func main() {
	newUuid()
}
