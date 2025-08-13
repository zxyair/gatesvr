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
	"time"
)

// 路由号

const (
	publicKey  = "./pem/key.pub.pem"
	privateKey = "./pem/key.pem"
)

func main() {
	// 创建容器
	gamesvr := gatesvr.NewContainer()
	// 创建用户定位器
	locator := redis.NewLocator()
	// 创建服务发现
	registry := etcd.NewRegistry()
	// 创建节点组件
	component := node.NewNode(
		node.WithLocator(locator),
		node.WithRegistry(registry),
		node.WithCodec(json.DefaultCodec),
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

	go func() {
		time.Sleep(1 * time.Second)
		helper.HandleConsoleInput(proxy)
	}()

}

//// 打印组件信息
//func printInfo() {
//	infos := make([]string, 0, 6)
//	infos = append(infos, fmt.Sprintf("ID: %s", ))
//	infos = append(infos, fmt.Sprintf("Name: %s", component..Name()))
//	infos = append(infos, fmt.Sprintf("Link: %s", g.linker.ExposeAddr()))
//	infos = append(infos, fmt.Sprintf("Server: [%s] %s", g.opts.server.Protocol(), net.FulfillAddr(g.opts.server.Addr())))
//	infos = append(infos, fmt.Sprintf("Registry: %s", g.opts.registry.Name()))
//	if g.opts.encryptor != nil {
//		infos = append(infos, fmt.Sprintf("Encryptor: %s", g.opts.encryptor.Name()))
//	} else {
//		infos = append(infos, "Encryptor: -")
//	}
//	if g.opts.compressor != nil {
//		infos = append(infos, fmt.Sprintf("Compressor: %s", g.opts.compressor.Name()))
//	} else {
//		infos = append(infos, "Compressor: -")
//	}
//	info.PrintBoxInfo("Gate", infos...)
//}
