package gate

import (
	"context"
	"fmt"
	"gatesvr/cluster"
	"gatesvr/errors"
	"gatesvr/internal/link"
	"gatesvr/log"
	"gatesvr/mode"
	"gatesvr/packet"
	"gatesvr/utils/codes"
)

type proxy struct {
	gate       *Gate            // 网关服
	nodeLinker *link.NodeLinker // 节点链接器
}

func newProxy(gate *Gate) *proxy {
	return &proxy{gate: gate, nodeLinker: link.NewNodeLinker(gate.ctx, &link.Options{
		InsID:    gate.opts.id,
		InsKind:  cluster.Gate,
		Locator:  gate.opts.locator,
		Registry: gate.opts.registry,
	})}
}

// 绑定用户与网关间的关系
func (p *proxy) bindGate(ctx context.Context, cid, uid int64) error {
	err := p.gate.opts.locator.BindGate(ctx, uid, p.gate.opts.id)
	if err != nil {
		return err
	}

	p.trigger(ctx, cluster.Reconnect, cid, uid)

	return nil
}

// 解绑用户与网关间的关系
func (p *proxy) unbindGate(ctx context.Context, cid, uid int64) error {
	err := p.gate.opts.locator.UnbindGate(ctx, uid, p.gate.opts.id)
	if err != nil {
		log.Errorf("user unbind failed, gid: %s, cid: %d, uid: %d, err: %v", p.gate.opts.id, cid, uid, err)
	}

	return err
}

// 触发事件
func (p *proxy) trigger(ctx context.Context, event cluster.Event, cid, uid int64) {
	if mode.IsDebugMode() {
		//log.Debugf("trigger event, event: %v cid: %d uid: %d", event.String(), cid, uid)
	}

	if err := p.nodeLinker.Trigger(ctx, &link.TriggerArgs{
		Event: event,
		CID:   cid,
		UID:   uid,
	}); err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFoundEvent), errors.Is(err, errors.ErrNotFoundUserLocation):
			//log.Warnf("trigger event failed, cid: %d, uid: %d, event: %v, err: %v", cid, uid, event.String(), err)
		default:
			//log.Errorf("trigger event failed, cid: %d, uid: %d, event: %v, err: %v", cid, uid, event.String(), err)
		}
	}
}

// 投递消息
func (p *proxy) deliver(ctx context.Context, cid, uid int64, message []byte) {
	msg, err := packet.UnpackMessage(message)
	if err != nil {
		log.Errorf("unpack message failed: %v", err)
		return
	}
	if !msg.IsCritical && p.gate.opts.limiter != nil {
		if !p.gate.opts.limiter.GetToken() {
			//log.Errorf("token is not enough")
			log.Debugf("token is not enough")
			message := &packet.Notification{
				Code:    codes.TooManyRequests.Code(),
				Message: fmt.Sprintf("token is not enough, please try again later，seq: %d", msg.Seq),
			}
			p.processMessageToClient(cid, message)
			return
		}
	}
	//调用encryptor.Decrypt解密消息
	if p.gate.opts.encryptor != nil {
		decryptMsgBuffer, err := p.gate.opts.encryptor.Decrypt(msg.Buffer)
		if err != nil {
			log.Errorf("decrypt message failed: %v", err)
			return
		}
		msg.Buffer = decryptMsgBuffer
		//log.Debugf("gate解密请求后消息: %v,长度为%d", message, len(message))
		////解压缩
		if p.gate.opts.compressor != nil {
			compressedBuffer, err := p.gate.opts.compressor.Decompress(msg.Buffer)
			if err != nil {
				log.Errorf("compress message failed: %v", err)
				return
			}
			msg.Buffer = compressedBuffer
			//log.Debugf("gate解压缩请求消息: %v，长度为%d", msg.Buffer, len(msg.Buffer))
		}

		messageTemp, err := packet.PackMessage(&packet.Message{
			Seq:    msg.Seq,
			Route:  msg.Route,
			Buffer: msg.Buffer,
		})
		if err != nil {
			return
		}
		message = messageTemp

	}

	if err = p.nodeLinker.Deliver(ctx, &link.DeliverArgs{
		CID:     cid,
		UID:     uid,
		Route:   msg.Route,
		Message: message,
	}); err != nil {
		switch {
		case errors.Is(err, errors.ErrNotFoundRoute), errors.Is(err, errors.ErrNotFoundEndpoint):
			message := &packet.Notification{
				Code:    codes.NotFound.Code(),
				Message: fmt.Sprintf("deliver message failed, cid: %d uid: %d seq: %d route: %d err: %v", cid, uid, msg.Seq, msg.Route, err),
			}
			p.processMessageToClient(cid, message)
			log.Warnf("deliver message failed, cid: %d uid: %d seq: %d route: %d err: %v", cid, uid, msg.Seq, msg.Route, err)
		case errors.Is(err, errors.ErrNotFoundUserLocation):
			message := &packet.Notification{
				Code:    codes.StateError.Code(),
				Message: fmt.Sprintf("deliver message failed, cid: %d uid: %d seq: %d route: %d err: %v", cid, uid, msg.Seq, msg.Route, err),
			}
			p.processMessageToClient(cid, message)
			log.Warnf("deliver message failed, cid: %d uid: %d seq: %d route: %d err: %v", cid, uid, msg.Seq, msg.Route, err)
		default:
			log.Errorf("deliver message failed, cid: %d uid: %d seq: %d route: %d err: %v", cid, uid, msg.Seq, msg.Route, err)
		}
	}
}

// 开始监听
func (p *proxy) watch() {
	p.nodeLinker.WatchUserLocate()

	p.nodeLinker.WatchClusterInstance()
}
func (p *proxy) processMessageToClient(cid int64, message *packet.Notification) {
	buffer, err := p.gate.opts.codec.Marshal(message)
	res := &packet.Message{
		Route:  0,
		Buffer: buffer,
	}
	if err != nil {
		log.Errorf("marshal message failed: %v", err)
		return
	}
	if p.gate.opts.compressor != nil {
		res.Buffer, err = p.gate.opts.compressor.Compress(res.Buffer)
		if err != nil {
			return
		}

	}

	//加密
	if p.gate.opts.encryptor != nil {
		encryptMsgBuffer, err := p.gate.opts.encryptor.Encrypt(res.Buffer)
		if err != nil {
			return
		}
		res.Buffer = encryptMsgBuffer
	}
	messageEncry, err := packet.PackMessage(res)

	p.gate.session.Push(1, cid, messageEncry)
}
