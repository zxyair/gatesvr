package logics

import (
	"context"
	"gatesvr/cluster"
	"gatesvr/cluster/node"
	"gatesvr/log"
	"strconv"
)

func IsOnline(proxy *node.Proxy, uid string) (bool, error) {
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)
	gateID, err := proxy.LocateGate(context.Background(), uidInt64)
	onlineFlag, err := proxy.IsOnline(context.Background(), &cluster.IsOnlineArgs{
		GID:    gateID,
		Kind:   2,
		Target: uidInt64,
	})
	if err != nil {
		log.Debugf("uid %d Not Online Or Not Exist", uidInt64)
		return false, err
	}
	log.Debugf("uid %d Online: %t", uidInt64, onlineFlag)
	return onlineFlag, nil
}
