package heartbeat

import (
	"encoding/binary"
	"net"

	"google.golang.org/protobuf/proto"
)

type HeartBeatHandler struct {
	conn net.Conn
}

func NewHeartBeatHandler(conn net.Conn) *HeartBeatHandler {
	heartbeatHandler := &HeartBeatHandler{
		conn: conn,
	}

	return heartbeatHandler
}

func (heartbeatHandler *HeartBeatHandler) readN(buff []byte) error {
	byteRead := uint64(0)

	for byteRead < uint64(len(buff)) {
		n, error := heartbeatHandler.conn.Read(buff)

		if error != nil {
			return error
		}

		byteRead += uint64(n)
	}

	return nil
}

func (heartbeatHandler *HeartBeatHandler) writeN(buff []byte) error {
	byteWrite := uint64(0)

	for byteWrite < uint64(len(buff)) {
		n, error := heartbeatHandler.conn.Write(buff)

		if error != nil {
			return error
		}

		byteWrite += uint64(n)
	}

	return nil
}

func (heartbeatHandler *HeartBeatHandler) Send(heartbeat *HeartbeatMessage) error {
	serialized, err := proto.Marshal(heartbeat)

	if err != nil {
		return err
	}

	prefix := make([]byte, 8)
	binary.LittleEndian.PutUint64(prefix, uint64(len(serialized)))
	heartbeatHandler.writeN(prefix)
	heartbeatHandler.writeN(serialized)

	return nil
}

func (heartbeatHandler *HeartBeatHandler) Receive() (*HeartbeatMessage, error) {
	// serialized, err := proto.Marshal(wrapper)
	prefix := make([]byte, 8)
	heartbeatHandler.readN(prefix)

	payloadSize := binary.LittleEndian.Uint64(prefix)

	payload := make([]byte, payloadSize)

	heartbeatHandler.readN(payload)

	heartbeat := &HeartbeatMessage{}
	err := proto.Unmarshal(payload, heartbeat)

	return heartbeat, err
}

func (heartbeatHandler *HeartBeatHandler) Close() {
	heartbeatHandler.conn.Close()
}
