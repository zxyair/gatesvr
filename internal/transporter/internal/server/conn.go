package server

import (
	"context"
	"gatesvr/cluster"
	"gatesvr/core/buffer"
	"gatesvr/errors"
	"gatesvr/internal/transporter/internal/def"
	"gatesvr/internal/transporter/internal/protocol"
	"gatesvr/log"
	"gatesvr/utils/xtime"

	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Conn struct {
	ctx               context.Context    // 上下文
	cancel            context.CancelFunc // 取消函数
	server            *Server            // 连接管理
	rw                sync.RWMutex       // 锁
	conn              net.Conn           // TCP源连接
	state             int32              // 连接状态
	chData            chan chData        // 消息处理通道
	lastHeartbeatTime int64              // 上次心跳时间
	InsKind           cluster.Kind       // 集群类型
	InsID             string             // 集群ID
}

var i int32

func newConn(server *Server, conn net.Conn) *Conn {
	c := &Conn{}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.conn = conn
	c.server = server
	c.state = def.ConnOpened
	c.chData = make(chan chData, 10240)
	c.lastHeartbeatTime = xtime.Now().Unix()
	go c.read()

	go c.process()

	return c
}

// Send 发送消息
func (c *Conn) Send(buf buffer.Buffer) (err error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if err = c.checkState(); err != nil {
		return err
	}

	buf.Range(func(node *buffer.NocopyNode) bool {
		if _, err = c.conn.Write(node.Bytes()); err != nil {
			return false
		}
		return true
	})

	buf.Release()

	return
}

// 检测连接状态
func (c *Conn) checkState() error {
	if atomic.LoadInt32(&c.state) == def.ConnClosed {
		return errors.ErrConnectionClosed
	} else {
		return nil
	}
}

// 关闭连接
func (c *Conn) close(isNeedRecycle ...bool) error {
	if !atomic.CompareAndSwapInt32(&c.state, def.ConnOpened, def.ConnClosed) {
		return errors.ErrConnectionClosed
	}

	c.rw.Lock()
	defer c.rw.Unlock()

	c.cancel()

	close(c.chData)

	if len(isNeedRecycle) > 0 && isNeedRecycle[0] {
		c.server.recycle(c.conn)
	}

	return c.conn.Close()
}

// 读取消息
func (c *Conn) read() {
	conn := c.conn

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			isHeartbeat, route, _, data, err := protocol.ReadMessage(conn)
			if err != nil {
				_ = c.close(true)
				return
			}

			c.rw.RLock()

			if atomic.LoadInt32(&c.state) == def.ConnClosed {
				c.rw.RUnlock()
				return
			}

			c.chData <- chData{
				isHeartbeat: isHeartbeat,
				route:       route,
				data:        data,
			}
			c.rw.RUnlock()
		}
	}
}

// 处理数据
func (c *Conn) process() {
	ticker := time.NewTicker(def.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			deadline := xtime.Now().Add(-2 * def.HeartbeatInterval).Unix()
			if atomic.LoadInt64(&c.lastHeartbeatTime) < deadline {
				_ = c.close(true)
				return
			}
		case ch, ok := <-c.chData:
			if !ok {
				return
			}

			atomic.StoreInt64(&c.lastHeartbeatTime, xtime.Now().Unix())

			if ch.isHeartbeat {
				c.heartbeat()
			} else {
				handler, ok := c.server.handlers[ch.route]
				if !ok {
					continue
				}
				if err := handler(c, ch.data); err != nil && !errors.Is(err, errors.ErrNotFoundUserLocation) {
					log.Warnf("process route %d message failed: %v", ch.route, err)
				}
			}
		}
	}
}

// 响应心跳消息
func (c *Conn) heartbeat() {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if _, err := c.conn.Write(protocol.Heartbeat()); err != nil {
		log.Warnf("write heartbeat message error: %v", err)
	}
}
