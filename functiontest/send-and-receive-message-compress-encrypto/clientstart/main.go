package main

import (
	"gatesvr"
	"gatesvr/cluster"
	"gatesvr/cluster/client"
	"gatesvr/encoding/json"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/clientstart/benchtest"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/clientstart/helper"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/clientstart/pojo"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/route"
	"gatesvr/log"
	"gatesvr/network/tcp"
	"gatesvr/packet"
)

// 路由号

const (
	publicKey  = "./pem/key.pub.pem"
	privateKey = "./pem/key.pem"
)

func main() {
	// 创建容器
	container := gatesvr.NewContainer()
	//encryptor := rsa.NewEncryptor(
	//	rsa.WithEncryptorHash(hash.SHA256),
	//	rsa.WithEncryptorPadding(rsa.OAEP),
	//	rsa.WithEncryptorPublicKey(publicKey),
	//	rsa.WithEncryptorPrivateKey(privateKey),
	//)
	//compressor := lz4Compressor.NewCompressor()
	// 创建客户端组件
	component := client.NewClient(
		client.WithClient(tcp.NewClient()),
		client.WithCodec(json.DefaultCodec),
		//client.WithEncryptor(encryptor),
		//client.WithCompressor(compressor),
	)
	// 初始化监听
	initListen(component.Proxy())
	// 添加客户端组件
	container.Add(component)
	// 启动容器
	go container.Serve()

	// 防止主线程退出
	select {}
}

func initListen(proxy *client.Proxy) {
	// 监听组件启动
	proxy.AddHookListener(cluster.Start, startHandler)
	// 监听连接建立
	proxy.AddEventListener(cluster.Connect, connectHandler)
	// 监听消息回复
	proxy.AddRouteHandler(route.Greet, greetHandler)
	proxy.AddRouteHandler(route.AuthritionCheck, authuritionHandler)
	proxy.AddRouteHandler(route.ReceiveNotifications, receiveNotificationsHandler)
	proxy.AddRouteHandler(route.PressureTest, benchtest.PressurTestHandler)
}

// 组件启动处理器
func startHandler(proxy *client.Proxy) {
	go helper.HandleConsoleInput(proxy)
}

// 连接建立处理器
func connectHandler(conn *client.Conn) {
	// 启动控制台输入处理

}

// 消息回复处理器
func greetHandler(ctx *client.Context) {
	res := &pojo.GreetRes{}

	if err := ctx.Parse(res); err != nil {
		log.Errorf("invalid response message, err: %v", err)
		return
	}

	if res.Code != 0 {
		log.Errorf("nodestart response failed, code: %d", res.Code)
		return
	}

	log.Debugf("client收到响应：%+v", res)
	//time.AfterFunc(time.Second, func() {
	//	pushMessage(ctx.Conn())
	//})
	//pushMessage(ctx.Conn())
}

// 消息回复处理器
func authuritionHandler(ctx *client.Context) {
	res := &pojo.AuthuritionRes{}

	if err := ctx.Parse(res); err != nil {
		log.Errorf("invalid response message, err: %v", err)
		return
	}

	if res.Code != 0 {
		log.Errorf("nodestart response failed, code: %d", res.Code)
		return
	}

	log.Debugf("client收到响应：%+v", res)
}
func receiveNotificationsHandler(ctx *client.Context) {
	res := &packet.Notification{}

	if err := ctx.Parse(res); err != nil {
		log.Errorf("invalid response message, err: %v", err)
		return
	}

	if res.Code != 0 {
		log.Debugf("client收到失败通知：%+v", res)
		//log.Errorf("nodestart response failed, code: %d", res.Code)
		return
	}

	log.Debugf("client收到通知：%+v", res)
}
