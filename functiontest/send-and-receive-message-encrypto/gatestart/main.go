package main

import (
	"gatesvr"
	"gatesvr/core/hash"
	"gatesvr/crypto/rsa"
	"gatesvr/gate"
	"gatesvr/locate/redis"
	"gatesvr/network/tcp"
	"gatesvr/registry/etcd"
)

const (
	publicKey  = "./pem/key.pub.pem"
	privateKey = "./pem/key.pem"
)

func main() {

	// 创建容器
	gateSvr := gatesvr.NewContainer()
	// 创建服务器
	// 创建加密器
	encryptor := rsa.NewEncryptor(
		rsa.WithEncryptorHash(hash.SHA256),
		rsa.WithEncryptorPadding(rsa.OAEP),
		rsa.WithEncryptorPublicKey(publicKey),
		rsa.WithEncryptorPrivateKey(privateKey))
	server := tcp.NewServer()
	// 创建用户定位器
	locator := redis.NewLocator()
	// 创建服务发现
	registry := etcd.NewRegistry()
	// 创建网关组件
	component := gate.NewGate(
		gate.WithServer(server),
		gate.WithLocator(locator),
		gate.WithRegistry(registry),
		gate.WithEncryptor(encryptor),
	)
	// 添加网关组件
	gateSvr.Add(component)
	// 启动容器
	gateSvr.Serve()
}
