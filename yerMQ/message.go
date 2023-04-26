package yerMQ

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

const (
	IDLen     = 16
	minMsgLen = IDLen + 8 + 2
)

type MessageID [IDLen]byte

// Message byte 构造
//	[x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x][x]...
//	|       (int64)        ||    ||      (hex string encoded in ASCII)           || (binary)
//	|       8-byte         ||    ||                 16-byte                      || N-byte
//	------------------------------------------------------------------------------------------...
//	       timestamp          ^^                     msgID                          msgBody
//	                       (uint16)
//	                        2-byte
//	                        resend

type Message struct {
	ID        MessageID // 唯一ID
	msgBody   []byte
	Timestamp int64  // 时间戳 8byte
	resend    uint16 // 重发次数 2byte

	// 进入inFightQueue后的信息
	inTime     time.Time // 进入时长
	clientID   int64
	level      int64 // 优先级
	levelIndex int
}

func NewMessage(id MessageID, body []byte) *Message {
	return &Message{
		ID:        id,
		msgBody:   body,
		Timestamp: time.Now().UnixNano(),
	}
}

// 写msg的ID、msgBody、Timestamp、resend
func (m *Message) Write(w io.Writer) (int64, error) {
	// TODO
	return 0, nil
}

// 解码msg
func decodeMessage(b []byte) (*Message, error) {
	var msg Message
	if len(b) < minMsgLen {
		return nil, fmt.Errorf("invalid message buffer size (%d)", len(b))
	}

	msg.Timestamp = int64(binary.BigEndian.Uint64(b[:8]))
	msg.resend = uint16(binary.BigEndian.Uint64(b[8:10]))
	//msg.ID = binary.BigEndian.Uint64(b[10:10+IDLen])
	copy(msg.ID[:], b[10:10+IDLen])
	msg.msgBody = b[10+IDLen:]
	return &msg, nil
}

// TODO:msg写入后台磁盘队列中

