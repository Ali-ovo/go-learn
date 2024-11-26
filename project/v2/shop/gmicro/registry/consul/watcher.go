package consul

import (
	"context"

	"shop/gmicro/registry"
)

type watcher struct {
	event chan struct{}
	set   *serviceSet

	// for cancel
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *watcher) Next() (services []*registry.ServiceInstance, err error) {
	select {
	case <-w.ctx.Done():
		err = w.ctx.Err()
	case <-w.event: // w.event 是否存在值  不存在 hold 住
	}

	ss, ok := w.set.services.Load().([]*registry.ServiceInstance) // 原子性读 操作 读取 服务端相关信息

	if ok {
		services = append(services, ss...)
	}
	return
}

func (w *watcher) Stop() error {
	w.cancel()
	w.set.lock.Lock()
	defer w.set.lock.Unlock()
	delete(w.set.watcher, w)
	return nil
}
