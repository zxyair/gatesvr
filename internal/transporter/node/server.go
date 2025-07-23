package node

import (
	"context"
	"gatesvr/cluster"
	"gatesvr/errors"
	"gatesvr/internal/transporter/internal/codes"
	"gatesvr/internal/transporter/internal/protocol"
	"gatesvr/internal/transporter/internal/route"
	"gatesvr/internal/transporter/internal/server"
)

type Server struct {
	*server.Server
	provider Provider
}

func NewServer(addr string, provider Provider) (*Server, error) {
	serv, err := server.NewServer(&server.Options{Addr: addr})
	if err != nil {
		return nil, err
	}

	s := &Server{Server: serv, provider: provider}
	s.init()

	return s, nil
}

func (s *Server) init() {
	s.RegisterHandler(route.Trigger, s.trigger)
	s.RegisterHandler(route.Deliver, s.deliver)
	s.RegisterHandler(route.GetState, s.getState)
	s.RegisterHandler(route.SetState, s.setState)
}

// 触发事件
func (s *Server) trigger(conn *server.Conn, data []byte) error {
	seq, event, cid, uid, err := protocol.DecodeTriggerReq(data)
	if err != nil {
		return err
	}

	if conn.InsKind != cluster.Gate {
		return errors.ErrIllegalRequest
	}

	if err = s.provider.Trigger(context.Background(), conn.InsID, cid, uid, event); seq == 0 {
		if errors.Is(err, errors.ErrNotFoundSession) {
			return nil
		} else {
			return err
		}
	} else {
		return conn.Send(protocol.EncodeTriggerRes(seq, codes.ErrorToCode(err)))
	}
}

// 投递消息
func (s *Server) deliver(conn *server.Conn, data []byte) error {
	seq, cid, uid, message, err := protocol.DecodeDeliverReq(data)
	if err != nil {
		return err
	}

	var (
		gid string
		nid string
	)

	switch conn.InsKind {
	case cluster.Gate:
		gid = conn.InsID
	case cluster.Node:
		nid = conn.InsID
	default:
		return errors.ErrIllegalRequest
	}

	if err = s.provider.Deliver(context.Background(), gid, nid, cid, uid, message); seq == 0 {
		return err
	} else {
		return conn.Send(protocol.EncodeDeliverRes(seq, codes.ErrorToCode(err)))
	}
}

// 获取状态
func (s *Server) getState(conn *server.Conn, data []byte) error {
	seq, err := protocol.DecodeGetStateReq(data)
	if err != nil {
		return err
	}

	state, err := s.provider.GetState()

	return conn.Send(protocol.EncodeGetStateRes(seq, codes.ErrorToCode(err), state))
}

// 设置状态
func (s *Server) setState(conn *server.Conn, data []byte) error {
	seq, state, err := protocol.DecodeSetStateReq(data)
	if err != nil {
		return err
	}

	err = s.provider.SetState(state)

	return conn.Send(protocol.EncodeSetStateRes(seq, codes.ErrorToCode(err)))
}
