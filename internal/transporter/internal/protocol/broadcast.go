package protocol

import (
	"encoding/binary"
	"gatesvr/internal/transporter/internal/codes"

	"gatesvr/core/buffer"
	"gatesvr/errors"
	"gatesvr/internal/transporter/internal/route"
	"gatesvr/session"

	"io"
)

const (
	broadcastReqBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + b8
	broadcastResBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + defaultCodeBytes + b64
)

// EncodeBroadcastReq 编码广播请求
// 协议：size + header + route + seq + session kind + <message packet>
func EncodeBroadcastReq(seq uint64, kind session.Kind, message buffer.Buffer) buffer.Buffer {
	buf := buffer.NewNocopyBuffer()
	writer := buf.Malloc(broadcastReqBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(broadcastReqBytes-defaultSizeBytes+message.Len()))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.Broadcast)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteUint8s(uint8(kind))
	buf.Mount(message)

	return buf
}

// DecodeBroadcastReq 解码广播请求
// 协议：size + header + route + seq + session kind + <message packet>
func DecodeBroadcastReq(data []byte) (seq uint64, kind session.Kind, message []byte, err error) {
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

	message = data[broadcastReqBytes:]

	return
}

// EncodeBroadcastRes 编码广播响应
// 协议：size + header + route + seq + code + [total]
func EncodeBroadcastRes(seq uint64, code uint16, total ...uint64) buffer.Buffer {
	size := broadcastResBytes - defaultSizeBytes
	if code != codes.OK || len(total) == 0 || total[0] == 0 {
		size -= b64
	}

	buf := buffer.NewNocopyBuffer()
	writer := buf.Malloc(broadcastResBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(size))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.Broadcast)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteUint16s(binary.BigEndian, code)

	if code == codes.OK && len(total) > 0 && total[0] != 0 {
		writer.WriteUint64s(binary.BigEndian, total[0])
	}

	return buf
}

// DecodeBroadcastRes 解码广播响应
// 协议：size + header + route + seq + code + [total]
func DecodeBroadcastRes(data []byte) (code uint16, total uint64, err error) {
	if len(data) != broadcastResBytes && len(data) != broadcastResBytes-b64 {
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

	if code == codes.OK && len(data) == broadcastResBytes {
		total, err = reader.ReadUint64(binary.BigEndian)
	}

	return
}
