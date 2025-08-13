package logics

import (
	"context"
	"gatesvr/cluster"
	"gatesvr/cluster/node"
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
		return false, err
	}
	return onlineFlag, nil
}
func OnlineUserCount(proxy *node.Proxy) (int64, error) {
	usercount, err := proxy.Stat(context.Background(), 2)
	if err != nil {
		return 0, err
	} else {
		return usercount, nil
	}
}
