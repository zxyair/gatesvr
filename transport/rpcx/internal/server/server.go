package server

import (
	"crypto/tls"
	"github.com/dobyte/due/v2/core/endpoint"
	xnet "github.com/dobyte/due/v2/core/net"
	"github.com/dobyte/due/v2/errors"
	"github.com/smallnest/rpcx/server"
)

const scheme = "rpcx"

type Server struct {
	listenAddr       string
	exposeAddr       string
	server           *server.Server
	endpoint         *endpoint.Endpoint
	disabledServices []string
}

type Options struct {
	Addr       string
	KeyFile    string
	CertFile   string
	ServerOpts []server.OptionFn
}

func NewServer(opts *Options) (*Server, error) {
	listenAddr, exposeAddr, err := xnet.ParseAddr(opts.Addr)
	if err != nil {
		return nil, err
	}

	isSecure := false
	serverOpts := make([]server.OptionFn, 0)
	serverOpts = append(serverOpts, opts.ServerOpts...)
	if opts.CertFile != "" && opts.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
		if err != nil {
			return nil, err
		}
		serverOpts = append(serverOpts, server.WithTLSConfig(&tls.Config{Certificates: []tls.Certificate{cert}}))
		isSecure = true
	}

	s := &Server{}
	s.listenAddr = listenAddr
	s.exposeAddr = exposeAddr
	s.server = server.NewServer(serverOpts...)
	s.endpoint = endpoint.NewEndpoint(scheme, exposeAddr, isSecure)

	return s, nil
}

// Addr 监听地址
func (s *Server) Addr() string {
	return s.listenAddr
}

// Scheme 协议
func (s *Server) Scheme() string {
	return scheme
}

// Endpoint 获取服务端口
func (s *Server) Endpoint() *endpoint.Endpoint {
	return s.endpoint
}

// Start 启动服务器
func (s *Server) Start() error {
	return s.server.Serve("tcp", s.listenAddr)
}

// Stop 停止服务器
func (s *Server) Stop() error {
	return s.server.Close()
}

// RegisterService 注册服务
func (s *Server) RegisterService(desc, ss interface{}) error {
	name, ok := desc.(string)
	if !ok {
		return errors.New("invalid dispatcher desc")
	}

	return s.server.RegisterName(name, ss, "")
}
