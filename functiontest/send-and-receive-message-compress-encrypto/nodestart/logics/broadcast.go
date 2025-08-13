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
)

func Broadcast(proxy *node.Proxy, seq string, message string) {
	res := &packet.Notification{}
	res.Code = codes.OK.Code()
	res.Message = fmt.Sprintf("server broadcast message: %v", message)
	seqInt, _ := strconv.ParseInt(seq, 10, 32)
	err := proxy.Broadcast(context.Background(), &cluster.BroadcastArgs{
		Kind: 1,
		Message: &cluster.Message{
			Route:      route.ReceiveNotifications,
			Seq:        int32(seqInt),
			IsCritical: true,
			Data:       res,
		},
	})
	if err != nil {
		log.Errorf("broadcast failed: %v", err)
		return
	} else {
		log.Infof("成功推送广播消息给网关")
	}
}
