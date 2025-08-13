package logics

import (
	"context"
	"fmt"
	"gatesvr/cluster"
	"gatesvr/cluster/node"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/nodestart/pojo"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/route"
	"gatesvr/log"
	"gatesvr/packet"
	"gatesvr/utils/codes"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var reqPool = sync.Pool{New: func() any {
	return &pojo.GreetReq{}
}}

var resPool = sync.Pool{New: func() any {
	return &pojo.GreetRes{}
}}

func GreetHandler(ctx node.Context) {
	req := &pojo.GreetReq{}
	res := &pojo.GreetRes{}
	defer func() {
		if err := ctx.Reply(&cluster.Message{
			Route:      route.Greet,
			Seq:        ctx.Seq(),
			IsCritical: false,
			Data:       res,
		}); err != nil {
			log.Debugf("client disconnected, stop pushing: %v", err)
			return
		}
	}()

	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		res.Code = codes.InternalError.Code()
		return
	}
	//log.Debugf("node对请求反序列化后的消息内容为: %v", req)

	log.Debugf("node收到request message: %v", req)

	// 输出上下文值
	if val := ctx.GetValue("example_key"); val != nil {
		fmt.Printf("  Custom Value: %v\n", val)
	}
	res.Code = codes.OK.Code()
	res.Message = fmt.Sprintf("server reply +%d", ctx.Seq())
	//log.Debugf("node返回原始响应为: %+v", res)
	//res.Message = fmt.Sprintf("I'm tcp server, and the current time is: %s", xtime.Now().Format(xtime.DateTime))
}

func PressureTestHandler(ctx node.Context) {
	req := reqPool.Get().(*pojo.GreetReq)
	res := resPool.Get().(*pojo.GreetRes)
	defer reqPool.Put(req)
	defer resPool.Put(res)
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	}()

	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		return
	}
	res.Message = req.Message
}

func NodeContinuePush(proxy *node.Proxy, uid string, message string) {

	ticker := time.NewTicker(time.Millisecond) // 每秒 1000 次
	defer ticker.Stop()
	timeout := time.After(50 * time.Millisecond)
	var count int32
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)

	for {
		select {
		case <-ticker.C:
			msg := &packet.Notification{}
			msg.Code = codes.OK.Code()
			msg.Message = fmt.Sprintf("server %s push message %s %d to uid %d", proxy.GetID(), message, atomic.AddInt32(&count, 1), uidInt64)
			err := proxy.Push(context.Background(), &cluster.PushArgs{
				Kind:   2,
				Target: uidInt64,
				Message: &cluster.Message{
					Route:      route.ReceiveNotifications,
					Seq:        1001,
					IsCritical: true,
					Data:       msg,
				},
			})
			if err != nil {
				log.Errorf("push failed: %v", err)
				return
			}
		case <-timeout:
			return
		}
	}

}
func ForwardMessage(ctx node.Context) {
	req := &pojo.GreetReq{}
	res := &packet.Notification{}

	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		res.Code = codes.InternalError.Code()
		res.Message = err.Error()
		return
	}
	log.Debugf("node收到广播通知为: %s", req.Message)
	res.Code = codes.OK.Code()
	res.Message = fmt.Sprintf("server broadcast message: %v", req.Message)
	ctx.Proxy().Broadcast(ctx.Context(), &cluster.BroadcastArgs{
		Kind: 1,
		Message: &cluster.Message{
			Route:      0,
			Seq:        1010,
			IsCritical: true,
			Data:       res,
		},
	})
}

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

//func CheckConnection(ctx node.Context) {
//	req := &pojo.GreetReq{}
//	res := &pojo.GreetRes{}
//	defer func() {
//		if err := ctx.Reply(&cluster.Message{
//			Route:      0,
//			Seq:        ctx.Seq(),
//			IsCritical: true,
//			Data:       res,
//		}); err != nil {
//			log.Debugf("client disconnected, stop pushing: %v", err)
//			return
//		}
//	}()
//
//	if err := ctx.Parse(req); err != nil {
//		log.Errorf("parse request message failed: %v", err)
//		res.Code = codes.InternalError.Code()
//		return
//	}
//
//	// 输出上下文值
//	if val := ctx.GetValue("example_key"); val != nil {
//		fmt.Printf("  Custom Value: %v\n", val)
//	}
//	res.Code = codes.OK.Code()
//	res.Message = fmt.Sprintf("conn + %d 连接成功", ctx.CID())
//}
