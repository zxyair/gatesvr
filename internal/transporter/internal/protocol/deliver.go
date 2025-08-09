package protocol

import (
	"encoding/binary"
	"gatesvr/core/buffer"
	"gatesvr/errors"
	"gatesvr/internal/transporter/internal/route"
	"io"
)

const (
	deliverReqBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + b64 + b64
	deliverResBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + defaultCodeBytes
)

// EncodeDeliverReq 编码投递消息请求
// 协议：size4 + header1 + route1 + seq8 + cid8 + uid8 + <message packet>
func EncodeDeliverReq(seq uint64, cid int64, uid int64, message []byte) buffer.Buffer {
	buf := buffer.NewNocopyBuffer()
	writer := buf.Malloc(deliverReqBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(deliverReqBytes-defaultSizeBytes+len(message)))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.Deliver)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteInt64s(binary.BigEndian, cid, uid)
	buf.Mount(message)
	//log.Debugf("client 对请求protocol编码后的消息内容: %v,长度为%d", buf.Bytes(), len(buf.Bytes())) //输出buf中的内容，用log输出，用于调试
	return buf
}

// DecodeDeliverReq 解码投递消息请求
func DecodeDeliverReq(data []byte) (seq uint64, cid int64, uid int64, message []byte, err error) {
	reader := buffer.NewReader(data)

	if _, err = reader.Seek(defaultSizeBytes+defaultHeaderBytes+defaultRouteBytes, io.SeekStart); err != nil {
		return
	}

	if seq, err = reader.ReadUint64(binary.BigEndian); err != nil {
		return
	}

	if cid, err = reader.ReadInt64(binary.BigEndian); err != nil {
		return
	}

	if uid, err = reader.ReadInt64(binary.BigEndian); err != nil {
		return
	}

	message = data[deliverReqBytes:]

	//log.Debugf("node对请求protocol解码后的消息内容，seq: %v, cid: %v, uid: %v, message: %v", seq, cid, uid, message)
	return
}

// EncodeDeliverRes 编码投递消息响应
// 协议：size + header + route + seq + code
func EncodeDeliverRes(seq uint64, code uint16) buffer.Buffer {
	buf := buffer.NewNocopyBuffer()
	writer := buf.Malloc(deliverResBytes)
	writer.WriteUint32s(binary.BigEndian, uint32(deliverResBytes-defaultSizeBytes))
	writer.WriteUint8s(dataBit)
	writer.WriteUint8s(route.Deliver)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteUint16s(binary.BigEndian, code)

	return buf
}

// DecodeDeliverRes 解码投递消息响应
// 协议：size + header + route + seq + code
func DecodeDeliverRes(data []byte) (code uint16, err error) {
	if len(data) != deliverResBytes {
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
