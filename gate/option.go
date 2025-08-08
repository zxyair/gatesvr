package gate

import (
	"context"
	"gatesvr/circuitbreaker"
	"gatesvr/compress"
	"gatesvr/crypto"
	"gatesvr/encoding"
	"gatesvr/etc"
	"gatesvr/limite"
	"gatesvr/locate"
	"gatesvr/locate/redis"
	"gatesvr/network"
	"gatesvr/registry"
	"gatesvr/utils/xuuid"
	"time"
)

const (
	defaultName    = "gate"          // 默认名称
	defaultAddr    = ":0"            // 连接器监听地址
	defaultTimeout = 3 * time.Second // 默认超时时间
	defaultWeight  = 1
	defaultCodec   = "json" // 默认编解码器名称// 默认权重
)

const (
	defaultIDKey      = "etc.cluster.gate.id"
	defaultNameKey    = "etc.cluster.gate.name"
	defaultAddrKey    = "etc.cluster.gate.addr"
	defaultTimeoutKey = "etc.cluster.gate.timeout"
	defaultWeightKey  = "etc.cluster.gate.weight"
)

type options struct {
	ctx            context.Context                // 上下文
	id             string                         // 实例ID
	name           string                         // 实例名称
	addr           string                         // 监听地址
	timeout        time.Duration                  // RPC调用超时时间
	weight         int                            // 权重
	server         network.Server                 // 网关服务器
	locator        locate.Locator                 // 用户定位器
	registry       registry.Registry              // 服务注册器
	encryptor      crypto.Encryptor               // 消息加密器
	compressor     compress.Compressor            // 消息压缩器
	limiter        limite.Limiter                 // 限流器
	circutibreaker *circuitbreaker.CircuitBreaker // 熔断器
	codec          encoding.Codec                 // 编解码器
}
type Option func(o *options)

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    defaultName,
		addr:    defaultAddr,
		timeout: defaultTimeout,
		weight:  defaultWeight,
		codec:   encoding.Invoke(defaultCodec),
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

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithServer 设置服务器
func WithServer(server network.Server) Option {
	return func(o *options) { o.server = server }
}

// WithTimeout 设置RPC调用超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithLocator 设置用户定位器
func WithLocator(locator *redis.Locator) Option {
	return func(o *options) { o.locator = locator }
}

// WithEncryptor 设置消息加密器
func WithEncryptor(encryptor crypto.Encryptor) Option {
	return func(o *options) { o.encryptor = encryptor }
}

// WithRegistry 设置服务注册器
func WithRegistry(r registry.Registry) Option {
	return func(o *options) { o.registry = r }
}

func WithCompressor(compressor compress.Compressor) Option {
	return func(o *options) { o.compressor = compressor }
}
func WithLimiter(limiter limite.Limiter) Option {
	return func(o *options) { o.limiter = limiter }
}
func WithCircuitBreaker(cb *circuitbreaker.CircuitBreaker) Option {
	return func(o *options) { o.circutibreaker = cb }
}

// WithWeight 设置权重
func WithWeight(weight int) Option {
	return func(o *options) { o.weight = weight }
}
func WithCodec(codec encoding.Codec) Option {
	return func(o *options) { o.codec = codec }
}
