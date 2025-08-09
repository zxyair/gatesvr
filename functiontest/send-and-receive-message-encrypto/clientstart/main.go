package main

import (
	"fmt"
	"gatesvr"
	"gatesvr/cluster"
	"gatesvr/cluster/client"
	"gatesvr/core/hash"
	"gatesvr/crypto/rsa"
	"gatesvr/encoding/json"
	"gatesvr/log"
	"gatesvr/network/tcp"
)

// 路由号
const greet = 1
const (
	publicKey  = "./pem/key.pub.pem"
	privateKey = "./pem/key.pem"
)

func main() {
	// 创建容器
	container := gatesvr.NewContainer()
	encryptor := rsa.NewEncryptor(
		rsa.WithEncryptorHash(hash.SHA256),
		rsa.WithEncryptorPadding(rsa.OAEP),
		rsa.WithEncryptorPublicKey(publicKey),
		rsa.WithEncryptorPrivateKey(privateKey),
	)
	// 创建客户端组件
	component := client.NewClient(
		client.WithClient(tcp.NewClient()),
		client.WithCodec(json.DefaultCodec),
		client.WithEncryptor(encryptor),
	)
	// 初始化监听
	initListen(component.Proxy())
	// 添加客户端组件
	container.Add(component)
	// 启动容器
	container.Serve()
}

func initListen(proxy *client.Proxy) {
	// 监听组件启动
	proxy.AddHookListener(cluster.Start, startHandler)
	// 监听连接建立
	proxy.AddEventListener(cluster.Connect, connectHandler)
	// 监听消息回复
	proxy.AddRouteHandler(greet, greetHandler)
}

// 组件启动处理器
func startHandler(proxy *client.Proxy) {
	if _, err := proxy.Dial(); err != nil {
		log.Errorf("connect server failed: %v", err)
		return
	}
}

// 连接建立处理器
func connectHandler(conn *client.Conn) {
	pushMessage(conn)
}

// 消息回复处理器
func greetHandler(ctx *client.Context) {
	res := &greetRes{}

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

// 请求
type greetReq struct {
	Message string `json:"message"`
}

// 响应
type greetRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 推送消息
func pushMessage(conn *client.Conn) {
	msg := &cluster.Message{
		Route: 1,
		Data: &greetReq{
			Message: fmt.Sprintf("hello"),
		}}
	log.Debugf("client推送消原始消息为: Route: %d, Data: %+v", msg.Route, msg.Data)
	err := conn.Push(msg)
	if err != nil {
		log.Errorf("push message failed: %v", err)
	}

}
