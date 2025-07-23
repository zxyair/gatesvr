package protocol_test

import (
	"gatesvr/internal/transporter/internal/codes"
	"gatesvr/internal/transporter/internal/protocol"
	"gatesvr/session"
	"testing"
)

func TestDecodeIsOnlineReq(t *testing.T) {
	buffer := protocol.EncodeIsOnlineReq(1, session.User, 1)

	seq, kind, target, err := protocol.DecodeIsOnlineReq(buffer.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("seq: %v", seq)
	t.Logf("kind: %v", kind)
	t.Logf("target: %v", target)
}

func TestDecodeIsOnlineRes(t *testing.T) {
	buffer := protocol.EncodeIsOnlineRes(1, codes.NotFoundSession, false)

	code, isOnline, err := protocol.DecodeIsOnlineRes(buffer.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("code: %v", code)
	t.Logf("isOnline: %v", isOnline)
}
