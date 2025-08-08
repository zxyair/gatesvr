package logics

import (
	"context"
	"fmt"
	"gatesvr/cluster"
	"gatesvr/cluster/node"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/route"
	"gatesvr/log"
	"gatesvr/packet"
	"gatesvr/utils/codes"
	"strconv"
	"sync/atomic"
	"time"
)

// 模拟游戏对局中node不断推送消息给client
func ReconnectDemo(proxy *node.Proxy, uid string, message string) {
	ticker := time.NewTicker(time.Millisecond * 1) // 每秒 1000 次
	defer ticker.Stop()
	timeout := time.After(30 * time.Second)
	var count int32
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)
	online, err := IsOnline(proxy, uid)
	if err != nil || !online {
		log.Debugf("uid %d Not Online Or Not Exist", uidInt64)
		return
	}
	for {
		select {
		case <-ticker.C:
			online, err := IsOnline(proxy, uid)
			if err != nil || !online {
				log.Debugf("uid %d Not Online Or Not Exist", uidInt64)
				return
			}
			msg := &packet.Notification{}
			msg.Code = codes.OK.Code()
			msg.Message = fmt.Sprintf("server %s push message %s %d to uid %d", proxy.GetID(), message, atomic.AddInt32(&count, 1), uidInt64)
			err = proxy.Push(context.Background(), &cluster.PushArgs{
				Kind:   2,
				Target: uidInt64,
				Message: &cluster.Message{
					Route:      route.ReceiveNotifications,
					Seq:        atomic.LoadInt32(&count),
					IsCritical: true,
					Data:       msg,
				},
			})
			if err != nil {
				log.Errorf("push failed: %v", err)
				return
			}
			log.Debugf("push message to uid %d,message: %s", uidInt64, msg.Message)
		case <-timeout:
			return
		}
	}
}
func IsOnline(proxy *node.Proxy, uid string) (bool, error) {
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)
	gateID, err := proxy.LocateGate(context.Background(), uidInt64)
	onlineFlag, err := proxy.IsOnline(context.Background(), &cluster.IsOnlineArgs{
		GID:    gateID,
		Kind:   2,
		Target: uidInt64,
	})
	if err != nil {
		return false, err
	}
	return onlineFlag, nil
}
