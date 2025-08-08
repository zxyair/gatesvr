package main

import (
	"gatesvr"
	"gatesvr/compress/lz4Compressor"
	"gatesvr/core/hash"
	"gatesvr/crypto/rsa"
	"gatesvr/gate"
	"gatesvr/limite/tokenbucket"
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
	// 创建加密器
	encryptor := rsa.NewEncryptor(
		rsa.WithEncryptorHash(hash.SHA256),
		rsa.WithEncryptorPadding(rsa.OAEP),
		rsa.WithEncryptorPublicKey(publicKey),
		rsa.WithEncryptorPrivateKey(privateKey))
	// 创建服务器
	server := tcp.NewServer()
	// 创建用户定位器
	locator := redis.NewLocator(
		redis.WithPassword("123456"),
	)
	// 创建服务发现
	registry := etcd.NewRegistry()
	//创建压缩器
	compressor := lz4Compressor.NewCompressor()
	//创建限流器
	limiter := tokenbucket.NewTokenBucketRateLimtImpl(1, 1)
	// 创建网关组件
	component := gate.NewGate(
		gate.WithServer(server),
		gate.WithLocator(locator),
		gate.WithRegistry(registry),
		gate.WithEncryptor(encryptor),
		gate.WithCompressor(compressor),
		gate.WithLimiter(limiter),
	)
	// 添加网关组件
	gateSvr.Add(component)
	// 启动容器
	gateSvr.Serve()
}
