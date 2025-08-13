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
	isOnlineReqBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + b8 + b64
	isOnlineResBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + defaultCodeBytes + b8
)

// EncodeIsOnlineReq 编码检测用户是否在线请求
// 协议：size + header + route + seq + session kind + target
func EncodeIsOnlineReq(seq uint64, kind session.Kind, target int64) buffer.Buffer {
	buf := buffer.NewNocopyBuffer()
	writer := buf.Malloc(isOnlineReqBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(isOnlineReqBytes-defaultSizeBytes))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.IsOnline)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteUint8s(uint8(kind))
	writer.WriteInt64s(binary.BigEndian, target)

	return buf
}

// DecodeIsOnlineReq 解码检测用户是否在线请求
// 协议：size + header + route + seq + session kind + target
func DecodeIsOnlineReq(data []byte) (seq uint64, kind session.Kind, target int64, err error) {
	if len(data) != isOnlineReqBytes {
		err = errors.ErrInvalidMessage
		return
	}

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

	return
}

// EncodeIsOnlineRes 编码检测用户是否在线响应
// 协议：size + header + route + seq + code + online state
func EncodeIsOnlineRes(seq uint64, code uint16, isOnline bool) buffer.Buffer {
	buf := buffer.NewNocopyBuffer()
	writer := buf.Malloc(isOnlineResBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(isOnlineResBytes-defaultSizeBytes))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.IsOnline)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteUint16s(binary.BigEndian, code)
	writer.WriteBools(isOnline)

	return buf
}

// DecodeIsOnlineRes 解码检测用户是否在线响应
// 协议：size + header + route + seq + code + online state
func DecodeIsOnlineRes(data []byte) (code uint16, isOnline bool, err error) {
	if len(data) != isOnlineResBytes {
		err = errors.ErrInvalidMessage
		return
	}
	reader := buffer.NewReader(data)

	if _, err = reader.Seek(defaultSizeBytes+defaultHeaderBytes+defaultRouteBytes+defaultSeqBytes, io.SeekStart); err != nil {
		return
	}

	if code, err = reader.ReadUint16(binary.BigEndian); err != nil {
		return
	}

	if isOnline, err = reader.ReadBool(); err != nil {
		return
	}

	return
}
