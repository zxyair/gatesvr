package link

import (
	"gatesvr/cluster"
	"gatesvr/crypto"
	"gatesvr/encoding"
	"gatesvr/internal/dispatcher"
	"gatesvr/locate"
	"gatesvr/registry"
)

type Options struct {
	InsID           string                     // 实例ID
	InsKind         cluster.Kind               // 实例类型
	Codec           encoding.Codec             // 编解码器
	Locator         locate.Locator             // 定位器
	Registry        registry.Registry          // 注册器
	Encryptor       crypto.Encryptor           // 加密器
	BalanceStrategy dispatcher.BalanceStrategy // 负载均衡策略
}
