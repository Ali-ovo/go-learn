package consul

import (
	"sync"
	"sync/atomic"

	"shop/gmicro/registry"
)

type serviceSet struct {
	serviceName string
	watcher     map[*watcher]struct{}
	services    *atomic.Value // 原子性操作
	lock        sync.RWMutex  // 读写锁
}

func (s *serviceSet) broadcast(ss []*registry.ServiceInstance) {
	//原子操作， 保证线程安全， 我们平时写struct的时候
	s.services.Store(ss) // 是 atomic.Value 一个方法 用于将新值存储到atomic.Value变量中
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.watcher {
		select {
		case k.event <- struct{}{}: // 写入 空结构
		default:
		}
	}
}
