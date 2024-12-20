package app

import (
	"net/url"
	"os"
	"shop/gmicro/registry"
	"shop/gmicro/server"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type Option func(o *options)

type options struct {
	id        string
	name      string
	endpoints []*url.URL

	sigs []os.Signal

	// registrar 允许用户传入自己的注册 实现
	registrar registry.Registrar
	// registrarTimeout 注册超时退出
	registrarTimeout time.Duration
	// stopTimeout 注销超时退出
	stopTimeout time.Duration
	servers     []server.Server
	//// rpc 服务
	//rpcServer []*rpcserver.Server
	//// http 服务
	//restServer *restserver.Server
}

func DefaultOptions() options {
	o := options{
		sigs:             []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT},
		registrarTimeout: 10 * time.Second, // 注册服务 超时时间
		stopTimeout:      10 * time.Second, // 注销服务 超时时间
	}

	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	return o
}

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithSigs(sigs []os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}

//func WithRPCServer(server *rpcserver.Server) Option {
//	return func(o *options) {
//		o.rpcServer = server
//	}
//}
//
//func WithRestServer(server *restserver.Server) Option {
//	return func(o *options) {
//		o.restServer = server
//	}
//}

// WithServer with servers.
func WithServer(srv ...server.Server) Option {
	return func(o *options) { o.servers = srv }
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}
