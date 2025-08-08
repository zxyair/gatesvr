package protocol

import (
	"encoding/binary"
	"gatesvr/core/buffer"
	"gatesvr/errors"
	"gatesvr/internal/transporter/internal/route"
	"gatesvr/session"
	"io"
)

const (
	pushReqBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + b8 + b64
	pushResBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + defaultCodeBytes
)

// EncodePushReq 编码推送请求
// 协议：size4 + header1 + route1 + seq8 + session kind1 + target8 + <message packet>
func EncodePushReq(seq uint64, kind session.Kind, target int64, message buffer.Buffer) buffer.Buffer {
	buf := buffer.NewNocopyBuffer()
	writer := buf.Malloc(pushReqBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(pushReqBytes-defaultSizeBytes+message.Len()))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.Push)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteUint8s(uint8(kind))
	writer.WriteInt64s(binary.BigEndian, target)
	buf.Mount(message)

	//log.Debugf("node返回响应后protocol编码后为: %v", buf.Bytes())

	return buf
}

// DecodePushReq 解码推送消息
// 协议：size + header + route + seq + session kind + target + <message packet>
func DecodePushReq(data []byte) (seq uint64, kind session.Kind, target int64, message []byte, err error) {
	reader := buffer.NewReader(data)

	if _, err = reader.Seek(defaultSizeBytes+defaultHeaderBytes+defaultRouteBytes, io.SeekStart); err != nil {
		return
	}

	if seq, err = reader.ReadUint64(binary.BigEndian); err != nil {
		return
	}

	var k uint8
	if k, err = reader.ReadUint8(); err != nil {
		return
	} else {
		kind = session.Kind(k)
	}

	if target, err = reader.ReadInt64(binary.BigEndian); err != nil {
		return
	}

	message = data[pushReqBytes:]
	//log.Debugf("gate收到返回响应后protocol解码后为: %v,长度为%d", message, len(message))
	return
}

// EncodePushRes 编码推送响应
// 协议：size + header + route + seq + code
func EncodePushRes(seq uint64, code uint16) buffer.Buffer {
	buf := buffer.NewNocopyBuffer()
	writer := buf.Malloc(pushResBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(pushResBytes-defaultSizeBytes))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.Push)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteUint16s(binary.BigEndian, code)

	return buf
}

// DecodePushRes 解码推送响应
// 协议：size + header + route + seq + code
func DecodePushRes(data []byte) (code uint16, err error) {
	if len(data) != pushResBytes {
		err = errors.ErrInvalidMessage
		return
	}

	reader := buffer.NewReader(data)

	if _, err = reader.Seek(-defaultCodeBytes, io.SeekEnd); err != nil {
		return
	}

	if code, err = reader.ReadUint16(binary.BigEndian); err != nil {
		return
	}

	return
}
