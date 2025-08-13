package protocol_test

import (
	"gatesvr/core/buffer"
	"gatesvr/internal/transporter/internal/codes"
	"gatesvr/internal/transporter/internal/protocol"
	"gatesvr/packet"
	"gatesvr/session"
	"testing"
)

func TestEncodeBroadcastReq(t *testing.T) {
	message, err := packet.PackMessage(&packet.Message{
		Route:  1,
		Seq:    2,
		Buffer: []byte("hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}

	buf := protocol.EncodeBroadcastReq(1, session.User, buffer.NewNocopyBuffer(message))

	t.Log(buf.Bytes())
}

func TestDecodeBroadcastReq(t *testing.T) {
	message, err := packet.PackMessage(&packet.Message{
		Route:  1,
		Seq:    2,
		Buffer: []byte("hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}

	buf := protocol.EncodeBroadcastReq(1, session.User, buffer.NewNocopyBuffer(message))

	seq, kind, message, err := protocol.DecodeBroadcastReq(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("seq: %v", seq)
	t.Logf("kind: %v", kind)
	t.Logf("message: %v", string(message))
}

func TestEncodeBroadcastRes(t *testing.T) {
	buf := protocol.EncodeBroadcastRes(1, codes.OK, 20)

	t.Log(buf.Bytes())
}

func TestDecodeBroadcastRes(t *testing.T) {
	buf := protocol.EncodeBroadcastRes(1, codes.OK, 20)

	code, total, err := protocol.DecodeBroadcastRes(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("code: %v", code)
	t.Logf("total: %v", total)
}
