package logics

import (
	"fmt"
	"gatesvr/cluster"
	"gatesvr/cluster/client"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/clientstart/pojo"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/route"
	"gatesvr/log"
	"sync/atomic"
	"time"
)

// 推送消息
func PushMessage(conn *client.Conn, message string) {
	msg := &cluster.Message{
		Route:      route.Greet,
		IsCritical: true,
		Data: &pojo.GreetReq{
			Message: message,
		}}
	log.Debugf("client推送消原始消息为: Route: %d, Data: %+v", msg.Route, msg.Data)
	err := conn.Push(msg)
	if err != nil {
		log.Errorf("push message failed: %v", err)
	} else {
		log.Debugf("成功发送给网关")
	}
}

func RetryConnect(conn *client.Conn, proxy *client.Proxy) {
	conn.Close()
	if _, err := proxy.Dial(); err != nil {
		log.Errorf("connect server failed: %v", err)
		return
	}
	log.Debugf("重新发起连接")
}

// 推送一次消息，不停接受服务端推送
func ClientContinuePush(conn *client.Conn) {
	ticker := time.NewTicker(time.Millisecond) // 每秒 1000 次
	defer ticker.Stop()
	timeout := time.After(50 * time.Millisecond)
	var count int32
	for {
		select {
		case <-ticker.C:
			msg := &cluster.Message{
				Route:      route.StatefulGreetRoute,
				IsCritical: true,
				Seq:        atomic.AddInt32(&count, 1),
				Data: &pojo.GreetReq{
					Message: fmt.Sprintf("client  %d 推送消息+%d 验证有序性", conn.UID(), atomic.LoadInt32(&count)),
				}}
			err := conn.Push(msg)
			if err != nil {
				log.Errorf("push message failed: %v", err)
			}
		case <-timeout:
			return
		}
	}
}

// 高频推送消息，持续 3 秒
func TestLimiter(conn *client.Conn) {
	ticker := time.NewTicker(time.Millisecond) // 每秒 1000 次
	defer ticker.Stop()
	timeout := time.After(50 * time.Millisecond)
	var count int32
	for {
		select {
		case <-ticker.C:
			msg := &cluster.Message{
				Route: route.Greet,
				//IsCritical: true,
				Seq: atomic.AddInt32(&count, 1),
				Data: &pojo.GreetReq{
					Message: fmt.Sprintf("高频消息+%d", atomic.LoadInt32(&count)),
				}}
			err := conn.Push(msg)
			if err != nil {
				log.Errorf("push message failed: %v", err)
			}
		case <-timeout:
			return
		}
	}
}
func TestLimiterCriticalMessage(conn *client.Conn) {
	ticker := time.NewTicker(time.Millisecond) // 每秒 1000 次
	defer ticker.Stop()
	timeout := time.After(50 * time.Millisecond)
	var count int32
	for {
		select {
		case <-ticker.C:
			msg := &cluster.Message{
				Route:      route.Greet,
				IsCritical: true,
				Seq:        atomic.AddInt32(&count, 1),
				Data: &pojo.GreetReq{
					Message: fmt.Sprintf("高频消息+%d", atomic.LoadInt32(&count)),
				}}
			err := conn.Push(msg)
			if err != nil {
				log.Errorf("push message failed: %v", err)
			}
		case <-timeout:
			return
		}
	}
}

func TestForwardMessage(conn *client.Conn) {
	msg := &cluster.Message{
		Route:      route.ForwardMessage,
		IsCritical: true,
		Data: &pojo.GreetReq{
			Message: "转发消息",
		},
	}
	err := conn.Push(msg)
	if err != nil {
		log.Errorf("push message failed: %v", err)
	}

}

func TestStaefulRoute(conn *client.Conn) {
	ticker := time.NewTicker(time.Millisecond) // 每秒 1000 次
	defer ticker.Stop()
	timeout := time.After(50 * time.Millisecond)
	var count int32
	for {
		select {
		case <-ticker.C:
			msg := &cluster.Message{
				Route:      route.StatefulGreetRoute,
				IsCritical: true,
				Seq:        atomic.AddInt32(&count, 1),
				Data: &pojo.GreetReq{
					Message: fmt.Sprintf("hello，server+%d", atomic.LoadInt32(&count)),
				}}
			err := conn.Push(msg)
			if err != nil {
				log.Errorf("push message failed: %v", err)
				break
			}
		case <-timeout:
			return
		}
	}
}
