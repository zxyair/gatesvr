package logics

import (
	"fmt"
	"gatesvr/cluster"
	"gatesvr/cluster/node"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/gamenode/pojo"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/route"
	"gatesvr/log"
	"gatesvr/utils/codes"
)

func StatefulGreet(ctx node.Context) {
	req := &pojo.GreetReq{}
	res := &pojo.GreetRes{}
	defer func() {
		if err := ctx.Reply(&cluster.Message{
			Route:      route.Greet,
			Seq:        ctx.Seq(),
			IsCritical: true,
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

	log.Debugf("node收到request message: %v", req)

	// 输出上下文值
	if val := ctx.GetValue("example_key"); val != nil {
		fmt.Printf("  Custom Value: %v\n", val)
	}
	res.Code = codes.OK.Code()
	res.Message = fmt.Sprintf("server reply +%d", ctx.Seq())
}
