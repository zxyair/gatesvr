package client

import (
	"gatesvr/cluster"
	"gatesvr/core/value"
	"gatesvr/network"
	"gatesvr/packet"
	"net"
	"sync"
)

type Conn struct {
	conn   network.Conn
	client *Client
	attrs  sync.Map
}

// ID 获取连接ID
func (c *Conn) ID() int64 {
	return c.conn.ID()
}

// UID 获取用户ID
func (c *Conn) UID() int64 {
	return c.conn.UID()
}

// Bind 绑定用户ID
func (c *Conn) Bind(uid int64) {
	c.conn.Bind(uid)
}

// Unbind 解绑用户ID
func (c *Conn) Unbind() {
	c.conn.Unbind()
}

// SetAttr 设置属性值
func (c *Conn) SetAttr(key, value any) {
	c.attrs.Store(key, value)
}

// GetAttr 获取属性值
func (c *Conn) GetAttr(key any) value.Value {
	if val, ok := c.attrs.Load(key); ok {
		return value.NewValue(val)
	} else {
		return value.NewValue()
	}
}

// DelAttr 删除属性值
func (c *Conn) DelAttr(key any) {
	c.attrs.Delete(key)
}

// LocalIP 获取本地IP
func (c *Conn) LocalIP() (string, error) {
	return c.conn.LocalIP()
}

// LocalAddr 获取本地地址
func (c *Conn) LocalAddr() (net.Addr, error) {
	return c.conn.LocalAddr()
}

// RemoteIP 获取远端IP
func (c *Conn) RemoteIP() (string, error) {
	return c.conn.RemoteIP()
}

// RemoteAddr 获取远端地址
func (c *Conn) RemoteAddr() (net.Addr, error) {
	return c.conn.RemoteAddr()
}

// Push 推送消息
func (c *Conn) Push(message *cluster.Message) error {
	var (
		err    error
		buffer []byte
	)

	if message.Data != nil {
		if v, ok := message.Data.([]byte); ok {
			buffer = v
		} else {
			buffer, err = c.client.opts.codec.Marshal(message.Data)
			//log.Debugf("client推送消息序列化后为: %v,消息长度%d", buffer, len(buffer))

			if err != nil {
				return err
			}
		}
		if c.client.opts.compressor != nil {
			buffer, err = c.client.opts.compressor.Compress(buffer)
			//log.Debugf("client推送消息压缩后为: %v,消息长度：%d", buffer, len(buffer))
			if err != nil {
				return err
			}
		}

		//加密
		if c.client.opts.encryptor != nil {
			buffer, err = c.client.opts.encryptor.Encrypt(buffer)
			//log.Debugf("client推送消息加密后为: %v,消息长度：%d", buffer, len(buffer))
			if err != nil {
				return err
			}
		}
	}

	msg, err := packet.PackMessage(&packet.Message{
		Seq:        message.Seq,
		Route:      message.Route,
		IsCritical: message.IsCritical,
		Buffer:     buffer,
	})
	//log.Debugf("client推送消息打包后为: %v", msg)
	if err != nil {
		return err
	}

	return c.conn.Push(msg)
}

// Close 关闭连接
func (c *Conn) Close(force ...bool) error {
	return c.conn.Close(force...)
}
