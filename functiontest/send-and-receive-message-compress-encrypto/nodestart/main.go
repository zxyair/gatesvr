package main

import (
	"gatesvr"
	"gatesvr/cluster/node"
	"gatesvr/encoding/json"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/nodestart/helper"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/nodestart/logics"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/route"
	"gatesvr/locate/redis"
	"gatesvr/registry/etcd"
	"time"
)

// 路由号

const (
	publicKey  = "./pem/key.pub.pem"
	privateKey = "./pem/key.pem"
)

func main() {
	// 创建容器
	lobbysvr := gatesvr.NewContainer()
	//encryptor := rsa.NewEncryptor(
	//	rsa.WithEncryptorHash(hash.SHA256),
	//	rsa.WithEncryptorPadding(rsa.OAEP),
	//	rsa.WithEncryptorPublicKey(publicKey),
	//	rsa.WithEncryptorPrivateKey(privateKey))
	// 创建用户定位器
	locator := redis.NewLocator()
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
	lobbysvr.Add(component)
	// 启动容器
	lobbysvr.Serve()

}

// 初始化应用
func initApp(proxy *node.Proxy) {
	proxy.Router().AddRouteHandler(route.Greet, false, logics.GreetHandler)
	proxy.Router().AddRouteHandler(route.AuthritionCheck, false, logics.AuthritionCheckHandler)
	proxy.Router().AddRouteHandler(route.ForwardMessage, false, logics.ForwardMessage)
	proxy.Router().AddRouteHandler(route.PressureTest, false, logics.PressureTestHandler)
	//proxy.Router().AddRouteHandler(checkConn, false, logics.CheckConnection)
	go func() {
		time.Sleep(1 * time.Second)
		helper.HandleConsoleInput(proxy)
	}()

}
