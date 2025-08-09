package main

import (
	"gatesvr"
	"gatesvr/cluster/node"
	"gatesvr/encoding/json"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/gamenode/helper"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/gamenode/logics"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/route"
	"gatesvr/locate/redis"
	"gatesvr/registry/etcd"
)

// 路由号

const (
	publicKey  = "./pem/key.pub.pem"
	privateKey = "./pem/key.pem"
)

func main() {
	// 创建容器
	gamesvr := gatesvr.NewContainer()
	//encryptor := rsa.NewEncryptor(
	//	rsa.WithEncryptorHash(hash.SHA256),
	//	rsa.WithEncryptorPadding(rsa.OAEP),
	//	rsa.WithEncryptorPublicKey(publicKey),
	//	rsa.WithEncryptorPrivateKey(privateKey))
	// 创建用户定位器
	locator := redis.NewLocator(
		redis.WithPassword("123456"),
	)
	// 创建服务发现
	registry := etcd.NewRegistry()
	// 创建节点组件
	component := node.NewNode(
		node.WithLocator(locator),
		node.WithRegistry(registry),
		node.WithCodec(json.DefaultCodec),
		//node.WithEncryptor(encryptor),
	)
	// 初始化应用
	initApp(component.Proxy())
	// 添加节点组件
	gamesvr.Add(component)
	// 启动容器
	gamesvr.Serve()

}

// 初始化应用
func initApp(proxy *node.Proxy) {
	proxy.Router().AddRouteHandler(route.StatefulGreetRoute, true, logics.StatefulGreet)
	go helper.HandleConsoleInput(proxy)

}
