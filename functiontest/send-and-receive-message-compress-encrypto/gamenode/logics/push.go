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

func Push(proxy *node.Proxy, uid string, message string) {
	res := &packet.Notification{}
	res.Code = codes.OK.Code()
	res.Message = fmt.Sprintf("server push message: %v", message)
	//uid string转int64
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)

	err := proxy.Push(context.Background(), &cluster.PushArgs{
		Kind:   2,
		Target: uidInt64,
		Message: &cluster.Message{
			Route:      route.ReceiveNotifications,
			Seq:        1001,
			IsCritical: true,
			Data:       res,
		},
	})
	if err != nil {
		log.Errorf("push failed: %v", err)
		return
	} else {
		log.Infof("成功推送消息给网关")
	}

}
