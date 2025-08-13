package tcp_test

import (
	"gatesvr/log"
	"gatesvr/network/tcp"
	"gatesvr/packet"
	"sync/atomic"

	"gatesvr/network"
	"net/http"
	_ "net/http/pprof"
	"testing"
)

var (
	seqCounter int32 // 新增包级计数器
)

func TestServer_Simple(t *testing.T) {
	server := tcp.NewServer()

	server.OnStart(func() {
		log.Info("server is started")
	})

	server.OnStop(func() {
		log.Info("server is stopped")
	})

	server.OnConnect(func(conn network.Conn) {
		log.Infof("connection is opened, connection id: %d", conn.ID())
	})

	server.OnDisconnect(func(conn network.Conn) {
		log.Infof("connection is closed, connection id: %d", conn.ID())
	})

	server.OnReceive(func(conn network.Conn, msg []byte) {
		message, err := packet.UnpackMessage(msg)
		if err != nil {
			log.Errorf("unpack message failed: %v", err)
			return
		}

		log.Infof("receive message from client, cid: %d, seq: %d, route: %d, msg: %s", conn.ID(), message.Seq, message.Route, string(message.Buffer))
		msg, err = packet.PackMessage(&packet.Message{
			Seq:    atomic.AddInt32(&seqCounter, 1), // 使用包级计数器
			Route:  1,
			Buffer: []byte("I'm fine~~"),
		})
		if err != nil {
			log.Errorf("pack message failed: %v", err)
			return
		}

		if err = conn.Push(msg); err != nil {
			log.Errorf("push message failed: %v", err)
		}
	})

	if err := server.Start(); err != nil {
		log.Fatalf("start server failed: %v", err)
	}

	select {}
}

func TestServer_Benchmark(t *testing.T) {
	server := tcp.NewServer(
		tcp.WithServerHeartbeatInterval(0),
	)

	server.OnStart(func() {
		log.Info("server is started")
	})

	server.OnReceive(func(conn network.Conn, msg []byte) {
		message, err := packet.UnpackMessage(msg)
		if err != nil {
			log.Errorf("unpack message failed: %v", err)
			return
		}

		data, err := packet.PackMessage(&packet.Message{
			Seq:    message.Seq,
			Route:  message.Route,
			Buffer: message.Buffer,
		})
		if err != nil {
			log.Errorf("pack message failed: %v", err)
			return
		}

		if err = conn.Push(data); err != nil {
			log.Errorf("push message failed: %v", err)
			return
		}
	})

	if err := server.Start(); err != nil {
		log.Fatalf("start server failed: %v", err)
	}

	go func() {
		err := http.ListenAndServe(":8089", nil)
		if err != nil {
			log.Errorf("pprof server start failed: %v", err)
		}
	}()

	select {}
}
func TestServer_MaxConnectionCount(t *testing.T) {
	server := tcp.NewServer()

	server.OnConnect(func(conn network.Conn) {
		log.Infof("connection is opened, connection id: %d", conn.ID())
	})

	server.OnDisconnect(func(conn network.Conn) {
		log.Infof("connection is closed, connection id: %d", conn.ID())
	})

	if err := server.Start(); err != nil {
		log.Fatalf("start server failed: %v", err)
	}

	select {}
}
