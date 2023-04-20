package yerMQ

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/golang/glog"
	"sync"
	"time"
)

var UuidChan = make(chan []byte, 1000)

type uuid int64

type uuidFactory struct {
}

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
		node.NewUuid()
	}
}

const (
	startTime           = int64(1577808000000)              // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69年
	timestampBits       = uint(41)                          // 时间戳占用位数
	workerBits          = uint(5)                           // 机器所占位数（区分集群下的不同机器）
	numberBits          = uint(12)                          // 序列号所占位数
	timeMax             = int64(-1 ^ (-1 << timestampBits)) // 时间戳最大值
	workerMax     int64 = -1 ^ (-1 << workerBits)           // 支持的最大机器id数
	numberMax     int64 = -1 ^ (-1 << numberBits)           // 支持的最大序列id数
	timeShift           = workerBits + numberBits           // 时间戳左移位数
	workerShift         = numberBits                        // 机器id左移位数
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

func (w *Worker) NewUuid() uuid {
	w.mu.Lock()
	now := time.Now().UnixNano() / 1000000 // 转毫秒
	if w.timestamp == now {
		w.number = (w.number + 1) & numberMax // 或(&)判断同1ms内序列号是否溢出
		if w.number == 0 {
			// 当前序列超过12bit-> 等待下1ms(number=0)
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		w.number = 0
	}
	t := now - startTime
	if t > timeMax {
		w.mu.Unlock()
		glog.Errorf("epoch must be between 0 and %d", numberMax-1)
		return 0
	}
	w.timestamp = now
	ID := uuid((now-startTime)<<timeShift | (w.workerId << workerShift) | (w.number))
	w.mu.Unlock()
	return ID
}

func (u uuid) Hex() MessageID {
	var h MessageID
	var b [8]byte

	b[0] = byte(u >> 56)
	b[1] = byte(u >> 48)
	b[2] = byte(u >> 40)
	b[3] = byte(u >> 32)
	b[4] = byte(u >> 24)
	b[5] = byte(u >> 16)
	b[6] = byte(u >> 8)
	b[7] = byte(u)

	hex.Encode(h[:], b[:])
	return h
}

// GetWorkerID 获取机器ID
func GetWorkerID(sid int64) (workerId int64) {
	workerId = (sid >> workerShift) & workerMax
	return
}

// GetTimestamp 获取时间戳
func GetTimestamp(sid int64) (timestamp int64) {
	timestamp = (sid >> timeShift) & timeMax
	return
}

// GetGenTimestamp 获取创建ID时的时间戳
func GetGenTimestamp(sid int64) (timestamp int64) {
	timestamp = GetTimestamp(sid) + startTime
	return
}

// GetGenTime 获取创建ID时的时间字符串(s)
func GetGenTime(sid int64) (t string) {
	t = time.Unix(GetGenTimestamp(sid)/1000, 0).Format("2006-01-02 15:04:05") // /1000换成秒
	return
}

// GetTimestampStatus 获取时间戳已使用的占比：范围（0.0 - 1.0）
func GetTimestampStatus() (state float64) {
	state = float64(time.Now().UnixNano()/1000000-startTime) / float64(timeMax)
	return
}
