package tcp

import (
	"gatesvr/utils/xtime"
	"sync"
	"testing"
)

func TestSavePendingMessages(t *testing.T) {
	mgr := &serverConnMgr{
		pendingMessages: sync.Map{},
	}

	conn := &serverConn{
		connMgr: mgr,
		uid:     123,
		chWrite: make(chan chWrite, 10),
	}

	// 发送测试消息
	msg := []byte("test message")
	conn.chWrite <- chWrite{typ: dataPacket, msg: msg}
	close(conn.chWrite)

	// 调用保存方法
	conn.savePendingMessages()

	// 验证消息是否保存
	value, ok := mgr.pendingMessages.Load(int64(123))
	t.Logf("Pending messages: %v", value)
	if !ok {
		t.Fatal("Expected pending messages for uid 123, but none found")
	}

	data := value.(*uidTimestamp)
	if len(data.messages) != 1 {
		t.Errorf("Expected 1 pending message, got %d", len(data.messages))
	}
	t.Logf("Pending message content: %v", data.messages[0].msg)
	if string(data.messages[0].msg) != string(msg) {
		t.Errorf("Expected message '%s', got '%s'", string(msg), string(data.messages[0].msg))
	}
}

func TestClearExpiredMessages(t *testing.T) {
	mgr := &serverConnMgr{
		pendingMessages: sync.Map{},
	}

	// 添加过期消息
	expiredTime := xtime.Now().UnixNano() - 4*60*1e9 // 4分钟前
	mgr.pendingMessages.Store(int64(123), &uidTimestamp{
		timestamp: expiredTime,
		messages:  []pendingMsg{{msg: []byte("expired message")}},
	})

	// 添加未过期消息
	validTime := xtime.Now().UnixNano() - 2*60*1e9 // 2分钟前
	mgr.pendingMessages.Store(int64(456), &uidTimestamp{
		timestamp: validTime,
		messages:  []pendingMsg{{msg: []byte("valid message")}},
	})

	// 调用清理方法
	mgr.clearExpiredMessages()

	// 验证过期消息是否被清理
	_, ok := mgr.pendingMessages.Load(int64(123))
	if ok {
		t.Error("Expected expired messages to be deleted, but they were found")
	}

	// 验证未过期消息是否保留
	value, ok := mgr.pendingMessages.Load(int64(456))
	if !ok {
		t.Fatal("Expected valid messages for uid 456, but none found")
	}

	data := value.(*uidTimestamp)
	if len(data.messages) != 1 {
		t.Errorf("Expected 1 pending message, got %d", len(data.messages))
	}
}
