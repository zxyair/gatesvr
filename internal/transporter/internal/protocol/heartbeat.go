package protocol

import (
	"encoding/binary"
	"gatesvr/core/buffer"
)

var heartbeat []byte

func init() {
	writer := buffer.NewWriter(defaultSizeBytes + defaultHeaderBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(defaultHeaderBytes))
	writer.WriteUint8s(heartbeatBit)
	heartbeat = writer.Bytes()
}

func Heartbeat() []byte {
	return heartbeat
}
