package client

import (
	"context"
	"fmt"
	"gatesvr/cluster"
	"gatesvr/component"
	"gatesvr/core/info"
	"gatesvr/errors"
	"gatesvr/log"
	"gatesvr/network"
	"gatesvr/packet"
	"gatesvr/utils/xcall"

	"sync"
	"sync/atomic"
)

type HookHandler func(proxy *Proxy)

type RouteHandler func(ctx *Context)

type EventHandler func(conn *Conn)

type Client struct {
	component.Base
	opts                *options
	ctx                 context.Context
	cancel              context.CancelFunc
	routes              map[int32][]RouteHandler
	events              map[cluster.Event][]EventHandler
	defaultRouteHandler RouteHandler
	proxy               *Proxy
	state               int32
	conns               sync.Map
	rw                  sync.RWMutex
	hooks               map[cluster.Hook][]HookHandler
}

func NewClient(opts ...Option) *Client {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	c := &Client{}
	c.opts = o
	c.proxy = newProxy(c)
	c.routes = make(map[int32][]RouteHandler)
	c.events = make(map[cluster.Event][]EventHandler)
	c.hooks = make(map[cluster.Hook][]HookHandler)
	c.ctx, c.cancel = context.WithCancel(o.ctx)
	c.state = int32(cluster.Shut)

	return c
}

// Name 组件名称
func (c *Client) Name() string {
	return c.opts.name
}

// Init 初始化节点
func (c *Client) Init() {
	if c.opts.client == nil {
		log.Fatal("client plugin is not injected")
	}

	if c.opts.codec == nil {
		log.Fatal("codec plugin is not injected")
	}

	c.runHookFunc(cluster.Init)
}

// Start 启动组件
func (c *Client) Start() {
	c.setState(cluster.Work)

	c.opts.client.OnDisconnect(c.handleDisconnect)
	c.opts.client.OnReceive(c.handleReceive)

	c.printInfo()

	c.runHookFunc(cluster.Start)
}

// Destroy 销毁组件
func (c *Client) Destroy() {
	c.setState(cluster.Shut)

	c.runHookFunc(cluster.Destroy)
}

// Proxy 获取节点代理
func (c *Client) Proxy() *Proxy {
	return c.proxy
}

// 处理断开连接
func (c *Client) handleDisconnect(conn network.Conn) {
	val, ok := c.conns.Load(conn)
	if !ok {
		return
	}

	c.conns.Delete(conn)

	handlers, ok := c.events[cluster.Disconnect]
	if !ok {
		return
	}

	for _, handler := range handlers {
		xcall.Call(func() {
			handler(val.(*Conn))
		})
	}
}

// 处理接收到的消息
func (c *Client) handleReceive(conn network.Conn, data []byte) {
	val, ok := c.conns.Load(conn)
	if !ok {
		return
	}

	message, err := packet.UnpackMessage(data)
	if err != nil {
		log.Errorf("unpack message failed: %v", err)
		return
	}

	handlers, ok := c.routes[message.Route]
	if ok {
		for _, handler := range handlers {
			xcall.Call(func() {
				handler(&Context{
					ctx:     context.Background(),
					conn:    val.(*Conn),
					message: message,
				})
			})
		}
	} else if c.defaultRouteHandler != nil {
		c.defaultRouteHandler(&Context{
			ctx:     context.Background(),
			conn:    val.(*Conn),
			message: message,
		})
	} else {
		log.Debugf("route handler is not registered, route: %v", message.Route)
	}
}

// 拨号
func (c *Client) dial(opts ...DialOption) (*Conn, error) {
	if c.getState() == cluster.Shut {
		return nil, errors.ErrClientShut
	}

	o := &dialOptions{attrs: make(map[string]any)}
	for _, opt := range opts {
		opt(o)
	}

	conn, err := c.opts.client.Dial(o.addr)
	if err != nil {
		return nil, err
	}

	cc := &Conn{conn: conn, client: c}

	for key, value := range o.attrs {
		cc.SetAttr(key, value)
	}

	c.conns.Store(conn, cc)

	if handlers, ok := c.events[cluster.Connect]; ok {
		for _, handler := range handlers {
			xcall.Call(func() {
				handler(cc)
			})
		}
	}

	return cc, nil
}

// 添加路由处理器
func (c *Client) addRouteHandler(route int32, handler RouteHandler) {
	if c.getState() == cluster.Shut {
		c.routes[route] = append(c.routes[route], handler)
	} else {
		log.Warnf("client is working, can't add route handler")
	}
}

// 默认路由处理器
func (c *Client) setDefaultRouteHandler(handler RouteHandler) {
	if c.getState() == cluster.Shut {
		c.defaultRouteHandler = handler
	} else {
		log.Warnf("client is working, can't set default route handler")
	}
}

// 添加事件处理器
func (c *Client) addEventListener(event cluster.Event, handler EventHandler) {
	if c.getState() == cluster.Shut {
		c.events[event] = append(c.events[event], handler)
	} else {
		log.Warnf("client is working, can't add event handler")
	}
}

// 添加钩子监听器
func (c *Client) addHookListener(hook cluster.Hook, handler HookHandler) {
	switch hook {
	case cluster.Destroy:
		c.rw.Lock()
		c.hooks[hook] = append(c.hooks[hook], handler)
		c.rw.Unlock()
	default:
		if c.getState() == cluster.Shut {
			c.hooks[hook] = append(c.hooks[hook], handler)
		} else {
			log.Warnf("server is working, can't add hook handler")
		}
	}
}

// 设置状态
func (c *Client) setState(state cluster.State) {
	atomic.StoreInt32(&c.state, int32(state))
}

// 获取状态
func (c *Client) getState() cluster.State {
	return cluster.State(atomic.LoadInt32(&c.state))
}

// 执行钩子函数
func (c *Client) runHookFunc(hook cluster.Hook) {
	c.rw.RLock()

	if handlers, ok := c.hooks[hook]; ok {
		wg := &sync.WaitGroup{}
		wg.Add(len(handlers))

		for i := range handlers {
			handler := handlers[i]
			xcall.Go(func() {
				handler(c.proxy)
				wg.Done()
			})
		}

		c.rw.RUnlock()

		wg.Wait()
	} else {
		c.rw.RUnlock()
	}
}

// 打印组件信息
func (c *Client) printInfo() {
	infos := make([]string, 0)
	infos = append(infos, fmt.Sprintf("Name: %s", c.Name()))
	infos = append(infos, fmt.Sprintf("Codec: %s", c.opts.codec.Name()))
	infos = append(infos, fmt.Sprintf("Protocol: %s", c.opts.client.Protocol()))

	if c.opts.encryptor != nil {
		infos = append(infos, fmt.Sprintf("Encryptor: %s", c.opts.encryptor.Name()))
	} else {
		infos = append(infos, "Encryptor: -")
	}

	info.PrintBoxInfo("Client", infos...)
}
