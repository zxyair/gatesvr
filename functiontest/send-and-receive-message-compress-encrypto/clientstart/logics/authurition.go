package logics

import (
	"gatesvr/cluster"
	"gatesvr/cluster/client"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/clientstart/pojo"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/route"
	"gatesvr/log"
)

// 鉴权
func Authorition(conn *client.Conn, uid int64) bool {
	msg := &cluster.Message{
		Route:      route.AuthritionCheck,
		IsCritical: true,
		Data: &pojo.AuthuritionReq{
			Message: uid,
		},
	}
	err := conn.Push(msg)
	if err != nil {
		log.Errorf("push message failed: %v", err)
		return false
	}
	return true
}
