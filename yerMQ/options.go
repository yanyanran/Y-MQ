package yerMQ

import (
	"crypto/md5"
	"hash/crc32"
	"io"
	"log"
	"os"
	"time"
)

type Options struct {
	ID int64

	TCPAddress  string
	HTTPAddress string

	BroadcastAddress  string
	BroadcastTCPPort  int
	BroadcastHTTPPort int

	HTTPClientConnectTimeout time.Duration
	HTTPClientRequestTimeout time.Duration

	// diskQueue
	DataPath        string
	MemQueueSize    int64
	MaxBytesPerFile int64
	SyncEvery       int64
	SyncTimeout     time.Duration

	// msg
	MsgTimeout    time.Duration
	MaxMsgTimeout time.Duration
	MaxMsgSize    int64
	MaxBodySize   int64
	MaxReqTimeout time.Duration
	ClientTimeout time.Duration

	// client
	FlushBufferTimeout time.Duration
}

func NewOptions() *Options {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	h := md5.New()
	io.WriteString(h, hostName)
	defaultID := int64(crc32.ChecksumIEEE(h.Sum(nil)) % 1024)

	return &Options{
		ID:                       defaultID,
		TCPAddress:               "122.0.0.0:5150",
		HTTPAddress:              "122.0.0.0:5151",
		BroadcastAddress:         hostName,
		BroadcastTCPPort:         0,
		BroadcastHTTPPort:        0,
		HTTPClientConnectTimeout: 2 * time.Second,
		HTTPClientRequestTimeout: 5 * time.Second,

		MemQueueSize:    10000,
		MaxBytesPerFile: 100 * 1024 * 1024,
		SyncEvery:       2500,
		SyncTimeout:     2 * time.Second,

		MsgTimeout:    60 * time.Second,
		MaxMsgTimeout: 15 * time.Minute,
		MaxMsgSize:    1024 * 1024,
		MaxBodySize:   5 * 1024 * 1024,
		MaxReqTimeout: 1 * time.Hour,
		ClientTimeout: 60 * time.Second,
	}
}