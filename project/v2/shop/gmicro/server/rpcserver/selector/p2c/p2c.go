package p2c

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	selector2 "shop/gmicro/server/rpcserver/selector"
	"shop/gmicro/server/rpcserver/selector/node/ewma"
)

const (
	forcePick = time.Second * 3
	// Name is balancer name
	Name = "p2c"
)

var _ selector2.Balancer = &Balancer{}

// New creates a p2c selector.
func New() selector2.Selector {
	return NewBuilder().Build()
}

// Balancer is p2c selector.
type Balancer struct {
	mu     sync.Mutex // 锁
	r      *rand.Rand // 随机
	picked int64
}

// choose two distinct nodes.
func (s *Balancer) prePick(nodes []selector2.WeightedNode) (nodeA selector2.WeightedNode, nodeB selector2.WeightedNode) {
	s.mu.Lock()
	a := s.r.Intn(len(nodes))
	b := s.r.Intn(len(nodes) - 1)
	s.mu.Unlock()
	if b >= a {
		b = b + 1
	}
	nodeA, nodeB = nodes[a], nodes[b]
	return
}

// Pick pick a node.
func (s *Balancer) Pick(ctx context.Context, nodes []selector2.WeightedNode) (selector2.WeightedNode, selector2.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector2.ErrNoAvailable
	}
	if len(nodes) == 1 {
		done := nodes[0].Pick()
		return nodes[0], done, nil
	}

	var pc, upc selector2.WeightedNode
	nodeA, nodeB := s.prePick(nodes) // 获取 2个 Node
	// meta.Weight is the weight set by the service publisher in discovery
	if nodeB.Weight() > nodeA.Weight() { // 比较权重
		pc, upc = nodeB, nodeA
	} else {
		pc, upc = nodeA, nodeB
	}

	// If the failed node has never been selected once during forceGap, it is forced to be selected once
	// Take advantage of forced opportunities to trigger updates of success rate and delay
	// atomic.CompareAndSwapInt64(&s.picked, 0, 1) 进行比较 如果 s.picked 是默认值0 则 把1 赋值进去 并且返回 true
	if upc.PickElapsed() > forcePick && atomic.CompareAndSwapInt64(&s.picked, 0, 1) {
		pc = upc
		atomic.StoreInt64(&s.picked, 0) // 执行完成后 重新赋值 0
	}
	done := pc.Pick()
	return pc, done, nil
}

// NewBuilder returns a selector builder with p2c balancer
func NewBuilder() selector2.Builder {
	return &selector2.DefaultBuilder{
		Node:     &ewma.Builder{},
		Balancer: &Builder{},
	}
}

// Builder is p2c builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() selector2.Balancer {
	return &Balancer{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
}
