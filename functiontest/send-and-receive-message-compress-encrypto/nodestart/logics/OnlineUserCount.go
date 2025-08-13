package logics

import (
	"context"
	"gatesvr/cluster/node"
)

func OnlineUserCount(proxy *node.Proxy) (int64, error) {
	usercount, err := proxy.Stat(context.Background(), 2)
	if err != nil {
		return -1, err
	} else {
		return usercount, nil
	}
}
