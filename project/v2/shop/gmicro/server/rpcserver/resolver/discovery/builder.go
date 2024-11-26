package discovery

import (
	"context"
	"errors"
	"strings"
	"time"

	"shop/gmicro/registry"

	"google.golang.org/grpc/resolver"
)

const name = "discovery"

// Option is builder option.
type Option func(o *builder)

// WithTimeout with timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(b *builder) {
		b.timeout = timeout
	}
}

// WithInsecure with isSecure option.
func WithInsecure(insecure bool) Option {
	return func(b *builder) {
		b.insecure = insecure
	}
}

type builder struct {
	discoverer registry.Discovery // 服务发现相关接口
	timeout    time.Duration
	insecure   bool // 是否安全
}

// NewBuilder creates a builder which is used to factory registry resolvers. 创建一个用于生成解析器的 注册器
func NewBuilder(d registry.Discovery, opts ...Option) resolver.Builder {
	b := &builder{
		discoverer: d,
		timeout:    time.Second * 10,
		insecure:   false,
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

// Build 构建 解析器
func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	/*
			使用 设计模式 - 观察者模式
			调用 Watch 由具体的实现方去实现 如 consul nacos 等  返回一个 watcher 给我
			启动协程 去监听 watcher 即可
			比如:
				我托 张三 帮我写个代码
		        我需要等待张三完成这个代码
				我怎么知道张三写完了呢
				张三 会把这代码上传到 github 上
				我只要 去 github 上 看即可
	*/
	var (
		err error
		w   registry.Watcher // 声明 服务监听器
	)
	done := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// 创建 服务监听器 在内部 维护更新数据
		w, err = b.discoverer.Watch(ctx, strings.TrimPrefix(target.URL.Path, "/"))
		close(done)
	}()
	select {
	case <-done: // done 有值 说明 上面的协程执行完成
	case <-time.After(b.timeout): // 说明超时了 默认 10s超时时间
		err = errors.New("discovery create watcher overtime")
	}
	if err != nil {
		cancel()
		return nil, err
	}
	r := &discoveryResolver{
		w:        w,          // 观察者 (监控器)
		cc:       cc,         // 客户端连接
		ctx:      ctx,        // context
		cancel:   cancel,     // cancel
		insecure: b.insecure, // 是否安全
	}
	go r.watch() // 不停的拉去最新数据
	return r, nil
}

// Scheme return scheme of discovery
func (*builder) Scheme() string {
	return name
}
