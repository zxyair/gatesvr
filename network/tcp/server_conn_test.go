package tcp

import (
	"gatesvr/network"
	"gatesvr/utils/xtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCheckAndSendPendingMessages(t *testing.T) {
	// 初始化serverConnMgr
	mgr := &serverConnMgr{
		pendingMessages: sync.Map{},
	}

	// 初始化serverConn
	conn := &serverConn{
		connMgr: mgr,
		uid:     123,
		chWrite: make(chan chWrite, 10),
		rw:      sync.RWMutex{},
	}
	atomic.StoreInt32(&conn.state, int32(network.ConnOpened))

	// 添加测试消息
	testMsg := []byte("test message")
	mgr.pendingMessages.Store(int64(123), &uidTimestamp{
		timestamp: xtime.Now().UnixNano(),
		messages:  []pendingMsg{{msg: testMsg}},
	})

	// 验证消息是否发送
	select {
	case msg := <-conn.chWrite:
		if msg.typ != dataPacket {
			t.Errorf("Expected dataPacket type, got %v", msg.typ)
		}
		if string(msg.msg) != string(testMsg) {
			t.Errorf("Expected message '%s', got '%s'", string(testMsg), string(msg.msg))
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected message to be sent, but none received")
	}

	// 验证消息是否被清理
	_, ok := mgr.pendingMessages.Load(int64(123))
	if ok {
		t.Error("Expected pending messages to be deleted, but they still exist")
	}
}

func TestCheckAndSendPendingMessages_NoMessages(t *testing.T) {
	mgr := &serverConnMgr{
		pendingMessages: sync.Map{},
	}

	conn := &serverConn{
		connMgr: mgr,
		uid:     123,
		chWrite: make(chan chWrite, 10),
		rw:      sync.RWMutex{},
	}
	atomic.StoreInt32(&conn.state, int32(network.ConnOpened))

	err := conn.checkAndSendPendingMessages()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	select {
	case <-conn.chWrite:
		t.Error("Expected no message to be sent, but got one")
	case <-time.After(100 * time.Millisecond):
		// 正常情况
	}
}

func TestCheckAndSendPendingMessages_ClosedConnection(t *testing.T) {
	mgr := &serverConnMgr{
		pendingMessages: sync.Map{},
	}

	conn := &serverConn{
		connMgr: mgr,
		uid:     123,
		chWrite: make(chan chWrite, 10),
		rw:      sync.RWMutex{},
	}
	atomic.StoreInt32(&conn.state, int32(network.ConnClosed))

	mgr.pendingMessages.Store(int64(123), &uidTimestamp{
		timestamp: xtime.Now().UnixNano(),
		messages:  []pendingMsg{{msg: []byte("test message")}},
	})

	err := conn.checkAndSendPendingMessages()
	if err == nil {
		t.Error("Expected error for closed connection, but got nil")
	}
}
