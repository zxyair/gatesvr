package main

import (
	"fmt"
	"gatesvr"
	"gatesvr/cluster/node"
	"gatesvr/encoding/json"
	"gatesvr/locate/redis"
	"gatesvr/log"
	"gatesvr/registry/etcd"
	"gatesvr/utils/codes"
)

// 路由号
const greet = 1
const (
	publicKey  = "./pem/key.pub.pem"
	privateKey = "./pem/key.pem"
)

func main() {
	// 创建容器
	nodesvr := gatesvr.NewContainer()
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
	nodesvr.Add(component)
	// 启动容器
	nodesvr.Serve()
}

// 初始化应用
func initApp(proxy *node.Proxy) {
	proxy.Router().AddRouteHandler(greet, false, greetHandler)
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

// 路由处理器
func greetHandler(ctx node.Context) {
	req := &greetReq{}
	res := &greetRes{}
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	}()

	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		res.Code = codes.InternalError.Code()
		return
	}
	log.Debugf("node反序列化后的消息内容为: %v", req)

	log.Debugf("node收到request message: %v", req)

	//log.Info(req.Message)

	res.Code = codes.OK.Code()
	res.Message = fmt.Sprintf("server")
	//res.Message = fmt.Sprintf("I'm tcp server, and the current time is: %s", xtime.Now().Format(xtime.DateTime))
}
