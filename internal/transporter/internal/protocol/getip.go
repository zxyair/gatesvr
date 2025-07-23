package protocol

import (
	"encoding/binary"
	"gatesvr/internal/transporter/internal/codes"

	"gatesvr/core/buffer"
	"gatesvr/errors"
	"gatesvr/internal/transporter/internal/route"
	"gatesvr/session"
	"gatesvr/utils/xnet"

	"io"
)

const (
	getIPReqBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + b8 + b64
	getIPResBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + defaultCodeBytes + b32
)

// EncodeGetIPReq 编码获取IP请求
// 协议：size + header + route + seq + session kind + target
func EncodeGetIPReq(seq uint64, kind session.Kind, target int64) buffer.Buffer {
	buf := buffer.NewNocopyBuffer()
	writer := buf.Malloc(getIPReqBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(getIPReqBytes-defaultSizeBytes))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.GetIP)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteUint8s(uint8(kind))
	writer.WriteInt64s(binary.BigEndian, target)

	return buf
}

// DecodeGetIPReq 解码获取IP请求
// 协议：size + header + route + seq + session kind + target
func DecodeGetIPReq(data []byte) (seq uint64, kind session.Kind, target int64, err error) {
	if len(data) != getIPReqBytes {
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

// EncodeGetIPRes 编码获取IP响应
// 协议：size + header + route + seq + code + [ip]
func EncodeGetIPRes(seq uint64, code uint16, ip ...string) buffer.Buffer {
	size := getIPResBytes - defaultSizeBytes
	if code != codes.OK || len(ip) == 0 || ip[0] == "" {
		size -= 4
	}

	buf := buffer.NewNocopyBuffer()
	writer := buf.Malloc(getIPResBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(size))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.GetIP)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteUint16s(binary.BigEndian, code)

	if code == codes.OK && len(ip) > 0 && ip[0] != "" {
		writer.WriteUint32s(binary.BigEndian, xnet.IP2Long(ip[0]))
	}

	return buf
}

func DecodeGetIPRes(data []byte) (code uint16, ip string, err error) {
	if len(data) != getIPResBytes && len(data) != getIPResBytes-4 {
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

	if code == codes.OK && len(data) == getIPResBytes {
		var v uint32
		if v, err = reader.ReadUint32(binary.BigEndian); err != nil {
			return
		} else {
			ip = xnet.Long2IP(v)
		}
	}

	return
}
