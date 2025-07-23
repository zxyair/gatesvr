package gate_test

import (
	"context"
	"gatesvr/cluster"
	"gatesvr/session"
	"gatesvr/utils/xuuid"

	"gatesvr/internal/transporter/gate"

	"testing"
)

func TestBuilder(t *testing.T) {
	builder := gate.NewBuilder(&gate.Options{
		InsID:   xuuid.UUID(),
		InsKind: cluster.Node,
	})

	client, err := builder.Build("127.0.0.1:49899")
	if err != nil {
		t.Fatal(err)
	}

	ip, miss, err := client.GetIP(context.Background(), session.User, 1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("miss: %v ip: %v", miss, ip)

	ip, miss, err = client.GetIP(context.Background(), session.User, 1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("miss: %v ip: %v", miss, ip)
}
