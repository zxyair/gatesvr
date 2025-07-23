package main

import (
	"gatesvr"
	"gatesvr/gate"
	"gatesvr/locate/redis"
	"gatesvr/network/tcp"
	"gatesvr/registry/etcd"
)

func main() {
	// 创建容器
	container := gatesvr.NewContainer()
	// 创建服务器

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
	)
	// 添加网关组件
	container.Add(component)
	// 启动容器
	container.Serve()
}
