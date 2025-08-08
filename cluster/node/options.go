package node

import (
	"context"
	"gatesvr/crypto"
	"gatesvr/encoding"
	"gatesvr/etc"
	"gatesvr/locate"
	"gatesvr/registry"
	"gatesvr/transport"
	"gatesvr/utils/xuuid"

	"time"
)

const (
	defaultName    = "node"          // 默认节点名称
	defaultAddr    = ":0"            // 连接器监听地址
	defaultCodec   = "json"          // 默认编解码器名称
	defaultTimeout = 3 * time.Second // 默认超时时间
	defaultWeight  = 1               // 默认权重
)

const (
	defaultIDKey      = "etc.cluster.node.id"
	defaultNameKey    = "etc.cluster.node.name"
	defaultAddrKey    = "etc.cluster.node.addr"
	defaultCodecKey   = "etc.cluster.node.codec"
	defaultTimeoutKey = "etc.cluster.node.timeout"
	defaultWeightKey  = "etc.cluster.node.weight"
)

// SchedulingModel 调度模型
type SchedulingModel string

type Option func(o *options)

type options struct {
	ctx         context.Context       // 上下文
	id          string                // 实例ID
	name        string                // 实例名称；相同实例名称的节点，用户只能绑定其中一个
	addr        string                // 监听地址
	codec       encoding.Codec        // 编解码器
	timeout     time.Duration         // RPC调用超时时间
	locator     locate.Locator        // 用户定位器
	registry    registry.Registry     // 服务注册器
	encryptor   crypto.Encryptor      // 消息加密器
	transporter transport.Transporter // 消息传输器
	weight      int                   // 权重
}

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    defaultName,
		addr:    defaultAddr,
		codec:   encoding.Invoke(defaultCodec),
		timeout: defaultTimeout,
		weight:  defaultWeight,
	}

	if id := etc.Get(defaultIDKey).String(); id != "" {
		opts.id = id
	} else {
		opts.id = xuuid.UUID()
	}

	if name := etc.Get(defaultNameKey).String(); name != "" {
		opts.name = name
	}

	if addr := etc.Get(defaultAddrKey).String(); addr != "" {
		opts.addr = addr
	}

	if codec := etc.Get(defaultCodecKey).String(); codec != "" {
		opts.codec = encoding.Invoke(codec)
	}

	if timeout := etc.Get(defaultTimeoutKey).Duration(); timeout > 0 {
		opts.timeout = timeout
	}

	if weight := etc.Get(defaultWeightKey).Int(); weight > 0 {
		opts.weight = weight
	}

	return opts
}

// WithID 设置实例ID
func WithID(id string) Option {
	return func(o *options) { o.id = id }
}

// WithName 设置实例名称
func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

// WithAddr 设置连接地址
func WithAddr(addr string) Option {
	return func(o *options) { o.addr = addr }
}

// WithCodec 设置编解码器
func WithCodec(codec encoding.Codec) Option {
	return func(o *options) { o.codec = codec }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithTimeout 设置RPC调用超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithLocator 设置定位器
func WithLocator(locator locate.Locator) Option {
	return func(o *options) { o.locator = locator }
}

// WithRegistry 设置服务注册器
func WithRegistry(r registry.Registry) Option {
	return func(o *options) { o.registry = r }
}

// WithEncryptor 设置消息加密器
func WithEncryptor(encryptor crypto.Encryptor) Option {
	return func(o *options) { o.encryptor = encryptor }
}

// WithTransporter 设置消息传输器
func WithTransporter(transporter transport.Transporter) Option {
	return func(o *options) { o.transporter = transporter }
}

// WithWeight 设置权重
func WithWeight(weight int) Option {
	return func(o *options) { o.weight = weight }
}
