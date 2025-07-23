package protocol_test

import (
	"gatesvr/core/buffer"
	"gatesvr/internal/transporter/internal/codes"
	"gatesvr/internal/transporter/internal/protocol"
	"gatesvr/packet"
	"gatesvr/session"
	"testing"
)

func TestEncodeMulticastReq(t *testing.T) {
	message, err := packet.PackMessage(&packet.Message{
		Route:  1,
		Seq:    2,
		Buffer: []byte("hello w	orld"),
	})
	if err != nil {
		t.Fatal(err)
	}

	buf := protocol.EncodeMulticastReq(1, session.User, []int64{1, 2, 3}, buffer.NewNocopyBuffer(message))

	t.Log(buf.Bytes())
}

func TestDecodeMulticastReq(t *testing.T) {
	message, err := packet.PackMessage(&packet.Message{
		Route:  1,
		Seq:    2,
		Buffer: []byte("hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}

	buf := protocol.EncodeMulticastReq(1, session.User, []int64{1, 2, 3}, buffer.NewNocopyBuffer(message))

	seq, kind, targets, message, err := protocol.DecodeMulticastReq(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("seq: %v", seq)
	t.Logf("kind: %v", kind)
	t.Logf("targets: %v", targets)
	t.Logf("message: %v", string(message))
}

func TestEncodeMulticastRes(t *testing.T) {
	buf := protocol.EncodeMulticastRes(1, codes.OK, 20)

	t.Log(buf.Bytes())
}

func TestDecodeMulticastRes(t *testing.T) {
	buf := protocol.EncodeMulticastRes(1, codes.OK, 20)

	code, total, err := protocol.DecodeMulticastRes(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("code: %v", code)
	t.Logf("total: %v", total)
}
