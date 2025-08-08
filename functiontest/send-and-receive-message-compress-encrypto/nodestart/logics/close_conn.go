package logics

import (
	"context"
	"gatesvr/cluster"
	"gatesvr/cluster/node"
	"gatesvr/log"
	"strconv"
)

func ForceClose(proxy *node.Proxy, uid string) {

	//uid string转int64
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)
	gid, err := proxy.LocateGate(context.Background(), uidInt64)
	if err != nil {
		log.Errorf("push failed: %v", err)
		return
	}
	err = proxy.Disconnect(context.Background(), &cluster.DisconnectArgs{
		GID:    gid,
		Kind:   2,
		Target: uidInt64,
		Force:  true,
	})
	if err != nil {
		log.Errorf("push failed: %v", err)
		return
	}

}
func GracefulClose(proxy *node.Proxy, uid string) {

	//uid string转int64
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)
	gid, err := proxy.LocateGate(context.Background(), uidInt64)
	if err != nil {
		log.Errorf("push failed: %v", err)
		return
	}
	err = proxy.Disconnect(context.Background(), &cluster.DisconnectArgs{
		GID:    gid,
		Kind:   2,
		Target: uidInt64,
		Force:  false,
	})
	if err != nil {
		log.Errorf("push failed: %v", err)
		return
	}

}
