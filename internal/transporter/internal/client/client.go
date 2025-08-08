package client

import (
	"context"
	"gatesvr/circuitbreaker"
	"gatesvr/core/buffer"
	"gatesvr/errors"
	"gatesvr/log"

	"sync"
	"sync/atomic"
	"time"
)

const (
	ordered   = 20 // 有序消息连接数
	unordered = 10 // 无序消息连接数
)

const (
	defaultTimeout = 3 * time.Second // 调用超时时间
)

type chWrite struct {
	ctx  context.Context // 上下文
	seq  uint64          // 序列号
	buf  buffer.Buffer   // 数据Buffer
	call chan []byte     // 回调数据
}

type Client struct {
	opts           *Options       // 配置
	chWrite        chan *chWrite  // 写入队列
	connections    []*Conn        // 连接
	wg             sync.WaitGroup // 等待组
	closed         atomic.Bool    // 已关闭
	circuitbreaker *circuitbreaker.CircuitBreaker
}

func NewClient(opts *Options) *Client {
	c := &Client{}
	c.opts = opts
	c.chWrite = make(chan *chWrite, 10240)
	c.connections = make([]*Conn, 0, ordered+unordered)
	c.circuitbreaker = circuitbreaker.NewCircuitBreaker(3, 0.5, 3*time.Second)
	c.init()

	return c
}

// Call 调用
func (c *Client) Call(ctx context.Context, seq uint64, buf buffer.Buffer, idx ...int64) ([]byte, error) {
	if c.closed.Load() {
		return nil, errors.ErrClientClosed
	}

	call := make(chan []byte)

	conn := c.load(idx...)
	if err := conn.send(&chWrite{
		ctx:  ctx,
		seq:  seq,
		buf:  buf,
		call: call,
	}); err != nil {
		return nil, err
	}

	ctx1, cancel1 := context.WithTimeout(ctx, defaultTimeout)
	defer cancel1()

	select {
	case <-ctx.Done():
		log.Debugf("ctx client call timeout, conn: %v, seq: %d, buf: %v", conn, seq, buf)
		conn.cancel(seq)
		return nil, ctx.Err()
	case <-ctx1.Done():
		log.Debugf("ctx1 client call timeout, conn: %v, seq: %d, buf: %v", conn, seq, buf)
		conn.cancel(seq)
		return nil, ctx1.Err()
	case data := <-call:
		return data, nil
	}
}

// Send 发送
func (c *Client) Send(ctx context.Context, buf buffer.Buffer, idx ...int64) error {
	if c.closed.Load() {
		return errors.ErrClientClosed
	}

	conn := c.load(idx...)
	if !c.circuitbreaker.AllowRequest() {
		return errors.ErrServerCircuitBreaker
	}
	if err := conn.send(&chWrite{
		ctx:  ctx,
		buf:  buf,
		call: nil,
	}); err != nil {
		c.circuitbreaker.RecordFail()
		return err
	} else {
		c.circuitbreaker.RecordSuccess()
	}
	return nil
}

// 获取连接
func (c *Client) load(idx ...int64) *Conn {
	if len(idx) > 0 {
		return c.connections[idx[0]%ordered]
	} else {
		return c.connections[ordered]
	}
}

// 新建连接
func (c *Client) init() {
	c.wg.Add(ordered + unordered)

	go c.wait()

	for i := 0; i < ordered; i++ {
		c.connections = append(c.connections, newConn(c))
	}

	for i := 0; i < unordered; i++ {
		c.connections = append(c.connections, newConn(c, c.chWrite))
	}
}

// 连接断开
func (c *Client) done() {
	c.wg.Done()
}

// 等待客户端连接全部关闭
func (c *Client) wait() {
	c.wg.Wait()
	c.closed.Store(true)
	c.connections = nil

	time.AfterFunc(time.Second, func() {
		close(c.chWrite)
	})

	if c.opts.CloseHandler != nil {
		c.opts.CloseHandler()
	}
}
