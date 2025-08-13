package logics

import (
	"context"
	"fmt"
	"gatesvr/cluster/node"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/gamenode/pojo"
	"gatesvr/log"
	"gatesvr/utils/codes"
	"strconv"
)

func AuthritionCheckHandler(ctx node.Context) {
	req := &pojo.AuthuritionReq{}
	res := &pojo.AuthuritionRes{}
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	}()

	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		res.Code = codes.InternalError.Code()
		return
	}
	log.Debugf("uuid %d", req.Message)
	if err := ctx.BindGate(req.Message); err != nil {
		log.Errorf("bind gate failed: %v", err)
		res.Code = codes.InternalError.Code()
		res.Message = fmt.Sprintf("bind gate failed: %v", err)
		return
	} else {
		res.Code = codes.OK.Code()
		res.Message = fmt.Sprintf("authurition check success")
	}
}

func BindNode(proxy *node.Proxy, uid string) {
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)
	err := proxy.BindNode(context.Background(), uidInt64)
	if err != nil {
		log.Errorf("bind node failed: %v", err)
	} else {
		log.Debugf("bind node success")
	}
}
