package gate

import (
	"context"
	"gatesvr/cluster"
	"gatesvr/errors"
	"gatesvr/log"
	"gatesvr/packet"
	"gatesvr/session"
	"gatesvr/utils/xcall"
)

type provider struct {
	gate *Gate
}

// Bind 绑定用户与网关间的关系
func (p *provider) Bind(ctx context.Context, cid, uid int64) error {
	if cid <= 0 || uid <= 0 {
		return errors.ErrInvalidArgument
	}

	err := p.gate.session.Bind(cid, uid)
	if err != nil {
		return err
	}

	err = p.gate.proxy.bindGate(ctx, cid, uid)
	if err != nil {
		_, _ = p.gate.session.Unbind(uid)
	}

	return err
}

// Unbind 解绑用户与网关间的关系
func (p *provider) Unbind(ctx context.Context, uid int64) error {
	if uid == 0 {
		return errors.ErrInvalidArgument
	}

	cid, err := p.gate.session.Unbind(uid)
	if err != nil {
		return err
	}

	return p.gate.proxy.unbindGate(ctx, cid, uid)
}

// GetIP 获取客户端IP地址
func (p *provider) GetIP(ctx context.Context, kind session.Kind, target int64) (string, error) {
	return p.gate.session.RemoteIP(kind, target)
}

// IsOnline 检测是否在线
func (p *provider) IsOnline(ctx context.Context, kind session.Kind, target int64) (bool, error) {
	return p.gate.session.Has(kind, target)
}

// Stat 统计会话总数
func (p *provider) Stat(ctx context.Context, kind session.Kind) (int64, error) {
	return p.gate.session.Stat(kind)
}

// Disconnect 断开连接
func (p *provider) Disconnect(ctx context.Context, kind session.Kind, target int64, force bool) error {
	return p.gate.session.Close(kind, target, force)
}

// Push 发送消息
func (p *provider) Push(ctx context.Context, kind session.Kind, target int64, message []byte) error {

	messageEncry, err := p.processMessage(message)
	if err != nil {
		log.Errorf("processMessage failed: %v", err)
		return err
	}
	err = p.gate.session.Push(kind, target, messageEncry)

	if kind == session.User && errors.Is(err, errors.ErrNotFoundSession) {
		xcall.Go(func() {
			if err := p.gate.opts.locator.UnbindGate(ctx, target, p.gate.opts.id); err != nil {
				log.Errorf("unbind gate failed, uid = %d gid = %s err = %v", target, p.gate.opts.id, err)
			}
		})
	}

	return err
}

// Multicast 推送组播消息
func (p *provider) Multicast(ctx context.Context, kind session.Kind, targets []int64, message []byte) (int64, error) {
	messageEncry, err := p.processMessage(message)
	if err != nil {
		log.Errorf("processMessage failed: %v", err)
		return 0, err
	}
	return p.gate.session.Multicast(kind, targets, messageEncry)
}

// Broadcast 推送广播消息
func (p *provider) Broadcast(ctx context.Context, kind session.Kind, message []byte) (int64, error) {
	messageEncry, err := p.processMessage(message)
	if err != nil {
		log.Errorf("processMessage failed: %v", err)
		return 0, err
	}
	return p.gate.session.Broadcast(kind, messageEncry)
}

// GetState 获取状态
func (p *provider) GetState() (cluster.State, error) {
	return cluster.Work, nil
}

// SetState 设置状态
func (p *provider) SetState(state cluster.State) error {
	return nil
}

func (p *provider) processMessage(message []byte) ([]byte, error) {
	//拆包
	msg, err := packet.UnpackMessage(message)
	if err != nil {
		log.Errorf("unpack message failed: %v", err)
		return nil, err
	}
	//log.Debugf("gate 对响应消息拆包后消息为: %v", msg)
	//log.Debugf("gate 收到响应消息 %d", msg.Seq)

	//压缩
	if p.gate.opts.compressor != nil {
		msg.Buffer, err = p.gate.opts.compressor.Compress(msg.Buffer)
		//log.Debugf("gate 对响应消息压缩后消息为: %v,长度为%d", msg.Buffer, len(msg.Buffer))
		if err != nil {
			log.Errorf("compress failed: %v", err)
			return nil, err
		}
	}
	//log.Debugf("gate 对响应消息压缩后消息为: %v,长度为%d", compressedBuffer, len(compressedBuffer))

	//加密
	if p.gate.opts.encryptor != nil {
		encryptMsgBuffer, err := p.gate.opts.encryptor.Encrypt(msg.Buffer)
		if err != nil {
			return nil, err
		}
		//log.Debugf("gate 对响应消息加密后消息为: %v,消息长度为%d", encryptMsgBuffer, len(encryptMsgBuffer))
		msg.Buffer = encryptMsgBuffer
	}

	//打包
	messageEncry, err := packet.PackMessage(msg)
	//log.Debugf("gate 发送给client的响应消息为: %v", messageEncry)
	if err != nil {
		log.Errorf("pack message failed: %v", err)
		return nil, err
	}
	return messageEncry, nil
}
