package redis

import (
	"context"
	"gatesvr/etc"
	"github.com/go-redis/redis/v8"
)

const (
	defaultAddr       = "127.0.0.1:6379"
	defaultDB         = 0
	defaultMaxRetries = 3
	defaultPrefix     = "due"
)

const (
	defaultAddrsKey      = "etc.locate.redis.addrs"
	defaultDBKey         = "etc.locate.redis.db"
	defaultMaxRetriesKey = "etc.locate.redis.maxRetries"
	defaultPrefixKey     = "etc.locate.redis.prefix"
	defaultUsernameKey   = "etc.locate.redis.username"
	defaultPasswordKey   = "etc.locate.redis.password"
)

type Option func(o *Options)

type Options struct {
	Ctx context.Context

	// 客户端连接地址
	// 内建客户端配置，默认为[]string{"127.0.0.1:6379"}
	Addrs []string

	// 数据库号
	// 内建客户端配置，默认为0
	Db int

	// 用户名
	// 内建客户端配置，默认为空
	Username string

	// 密码
	// 内建客户端配置，默认为空
	Password string

	// 最大重试次数
	// 内建客户端配置，默认为3次
	MaxRetries int

	// 客户端
	// 外部客户端配置，存在外部客户端时，优先使用外部客户端，默认为nil
	client redis.UniversalClient

	// 前缀
	// key前缀，默认为due
	Prefix string
}

func defaultOptions() *Options {
	return &Options{
		Ctx:        context.Background(),
		Addrs:      etc.Get(defaultAddrsKey, []string{defaultAddr}).Strings(),
		Db:         etc.Get(defaultDBKey, defaultDB).Int(),
		MaxRetries: etc.Get(defaultMaxRetriesKey, defaultMaxRetries).Int(),
		Prefix:     etc.Get(defaultPrefixKey, defaultPrefix).String(),
		Username:   etc.Get(defaultUsernameKey).String(),
		Password:   etc.Get(defaultPasswordKey).String(),
	}
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *Options) { o.Ctx = ctx }
}

// WithAddrs 设置连接地址
func WithAddrs(addrs ...string) Option {
	return func(o *Options) { o.Addrs = addrs }
}

// WithDB 设置数据库号
func WithDB(db int) Option {
	return func(o *Options) { o.Db = db }
}

// WithUsername 设置用户名
func WithUsername(username string) Option {
	return func(o *Options) { o.Username = username }
}

// WithPassword 设置密码
func WithPassword(password string) Option {
	return func(o *Options) { o.Password = password }
}

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(maxRetries int) Option {
	return func(o *Options) { o.MaxRetries = maxRetries }
}

// WithClient 设置外部客户端
func WithClient(client redis.UniversalClient) Option {
	return func(o *Options) { o.client = client }
}

// WithPrefix 设置前缀
func WithPrefix(prefix string) Option {
	return func(o *Options) { o.Prefix = prefix }
}
