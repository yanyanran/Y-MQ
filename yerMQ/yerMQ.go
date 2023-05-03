package yerMQ

import (
	"Y-MQ/tools"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

type errStore struct {
	err error
}

type YERMQ struct {
	sync.Mutex
	ctx       context.Context
	ctxCancel context.CancelFunc

	clientIDIndex int64 // ID序列
	topicMap      map[string]*Topic
	opts          atomic.Value

	dl        *tools.DirLock
	isLoading int32
	isExiting int32
	errValue  atomic.Value
	startTime time.Time

	tcpServer    *tcpServer
	tcpListener  net.Listener
	httpListener net.Listener

	scanWorkerPoolSize int
	exitChan           chan int
	waitGroup          tools.WaitGroupWrapper
}

func New(opt *Options) (*YERMQ, error) {
	var err error
	dataPath := opt.DataPath
	if opt.DataPath == "" {
		pwd, _ := os.Getwd()
		dataPath = pwd
	}

	y := &YERMQ{
		startTime: time.Now(),
		topicMap:  make(map[string]*Topic),
		exitChan:  make(chan int),
		dl:        tools.NewDirLock(dataPath),
	}
	y.ctx, y.ctxCancel = context.WithCancel(context.Background())

	y.swapOpts(opt)
	y.errValue.Store(errStore{})

	err = y.dl.Lock() // dir lock
	if err != nil {
		return nil, fmt.Errorf("failed to lock data-path: %v", err)
	}

	// check config msg
	if opt.ID < 0 || opt.ID >= 1024 {
		return nil, errors.New("--node-id must be [0,1024)")
	}
	log.Printf(tools.String("yerMQ"))
	log.Printf("ID: %d", opt.ID)

	y.tcpServer = &tcpServer{yerMQ: y}
	y.httpListener, err = net.Listen(TypeOfAddr(opt.TCPAddress), opt.TCPAddress)
	if err != nil {
		return nil, fmt.Errorf("listen (%s) failed - %s", opt.TCPAddress, err)
	}

	if opt.HTTPAddress != "" {
		y.httpListener, err = net.Listen(TypeOfAddr(opt.HTTPAddress), opt.HTTPAddress)
		if err != nil {
			return nil, fmt.Errorf("listen (%s) failed - %s", opt.HTTPAddress, err)
		}
	}

	if opt.BroadcastHTTPPort == 0 { // ?
		tcpAddr, ok := y.RealHTTPAddr().(*net.TCPAddr)
		if ok {
			opt.BroadcastHTTPPort = tcpAddr.Port
		}
	}

	if opt.BroadcastTCPPort == 0 { // ?
		tcpAddr, ok := y.RealTCPAddr().(*net.TCPAddr)
		if ok {
			opt.BroadcastTCPPort = tcpAddr.Port
		}
	}

	return y, nil
}

// TODO: Metadata

type Metadata struct {
	Topics []TopicMetadata
}

type TopicMetadata struct {
	Name     string
	Stopped  bool
	Channels []ChannelMetadata
}

type ChannelMetadata struct {
	Name    string
	Stopped bool
}

func newMetadataFile(opts *Options) string {
	return path.Join(opts.DataPath, "yermq.dat")
}

func readMetadataFile(fn string) ([]byte, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read metadata from %s - %s", fn, err)
		}
	}
	return data, nil
}

func (y *YERMQ) LoadMetadata() error {
	atomic.StoreInt32(&y.isLoading, 1)
	defer atomic.StoreInt32(&y.isLoading, 0)

	fn := newMetadataFile(y.getOpts())
	data, err := readMetadataFile(fn)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}

	var m Metadata
	err = json.Unmarshal(data, &m) // 解析元数据
	if err != nil {
		return fmt.Errorf("failed to parse metadata in %s - %s", fn, err)
	}

	for _, t := range m.Topics {
		if !tools.IsValidTopicName(t.Name) {
			log.Fatalf("skipping creation of invalid topic %s\n", t.Name)
			continue
		}
		topic := y.GetTopic(t.Name)
		if t.Stopped {
			topic.Stop()
		}

		for _, c := range t.Channels {
			if !tools.IsValidChannelName(c.Name) {
				log.Fatalf("skipping creation of invalid channel %s", c.Name)
				continue
			}
			channel := topic.GetChannel(c.Name)
			if c.Stopped {
				channel.Stop()
			}
		}
		topic.Start()
	}
	return nil
}

func (y *YERMQ) Exit() {
	// TODO
}

func (y *YERMQ) GetTopic(topicName string) *Topic {
	// TODO
}

func (y *YERMQ) GetHaveTopic(topicName string) (*Topic, error) {
	// TODO
}

func (y *YERMQ) DeleteHaveTopic(topicName string) (*Topic, error) {
	// TODO
}

// TODO channel操作

func TypeOfAddr(addr string) string {
	if _, _, err := net.SplitHostPort(addr); err == nil {
		return "tcp"
	}
	return "unix"
}

func (y *YERMQ) getOpts() *Options {
	return y.opts.Load().(*Options)
}

func (y *YERMQ) swapOpts(opt *Options) {
	y.opts.Store(opt)
}

func (y *YERMQ) RealTCPAddr() net.Addr {
	if y.tcpListener == nil {
		return &net.TCPAddr{}
	}
	return y.tcpListener.Addr()
}

func (y *YERMQ) RealHTTPAddr() net.Addr {
	if y.httpListener == nil {
		return &net.TCPAddr{}
	}
	return y.httpListener.Addr()
}

func (y *YERMQ) Main() error {
	// TODO
}

func (y *YERMQ) queueScanWorker() {

}

func (y *YERMQ) queueScanLoop() {

}

func (y *YERMQ) Context() context.Context {
	return y.ctx
}
