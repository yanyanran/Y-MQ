package protocol

import (
	"encoding/binary"
	"io"
	"net"
)

type Protocol interface {
	NewClient(net.Conn) Client
	IOLoop(Client) error
}

type Client interface {
	Close() error
}

func SendResponse(w io.Writer, data []byte) (int, error) {
	// add a dataLen header
	err := binary.Write(w, binary.BigEndian, int32(len(data)))
	if err != nil {
		return 0, err
	}

	n, err := w.Write(data)
	if err != nil {
		return 0, nil
	}
	return n + 4, nil
}

func SendTypeResponse(w io.Writer, ty int32, data []byte) (int, error) {
	bBuf := make([]byte, 4)
	size := uint32(len(data)) + 4

	binary.BigEndian.PutUint32(bBuf, size)
	n, err := w.Write(bBuf)
	if err != nil {
		return n, err
	}

	binary.BigEndian.PutUint32(bBuf, uint32(ty))
	n, err = w.Write(bBuf)
	if err != nil {
		return n + 4, nil
	}

	n, err = w.Write(data)
	return n + 8, err
}
