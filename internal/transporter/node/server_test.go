package node_test

import (
	"context"
	"fmt"
	"gatesvr/cluster"
	"gatesvr/internal/transporter/node"
	"gatesvr/log"

	"testing"
	"time"
)

func TestServer(t *testing.T) {
	server, err := node.NewServer(":49898", &provider{})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("server listen on: %s", server.ListenAddr())

	go server.Start()

	<-time.After(20 * time.Second)
}

type provider struct {
}

// Trigger 触发事件
func (p *provider) Trigger(ctx context.Context, gid string, cid, uid int64, event cluster.Event) error {
	return nil
}

// Deliver 投递消息
func (p *provider) Deliver(ctx context.Context, gid, nid string, cid, uid int64, message []byte) error {
	log.Infof("gid: %s, nid: %s, cid: %d, uid: %d message: %s", gid, nid, cid, uid, string(message))
	return nil
}

func TestTimeout(t *testing.T) {
	ctx := context.Background()

	ctx1, _ := context.WithTimeout(ctx, 10*time.Second)

	ctx2, _ := context.WithTimeout(ctx1, 5*time.Second)

	fmt.Println(time.Now().Unix())

	select {
	case <-ctx1.Done():
		fmt.Println(1, time.Now().Unix())
	case <-ctx2.Done():
		fmt.Println(2, time.Now().Unix())
	}
}

// GetState 获取状态
func (p *provider) GetState() (cluster.State, error) {
	return cluster.Work, nil
}

// SetState 设置状态
func (p *provider) SetState(state cluster.State) error {
	return nil
}
