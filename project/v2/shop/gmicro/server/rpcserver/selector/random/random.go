package random

import (
	"context"
	"math/rand"

	selector2 "shop/gmicro/server/rpcserver/selector"
	"shop/gmicro/server/rpcserver/selector/node/direct"
)

const (
	// Name is random balancer name
	Name = "random"
)

var _ selector2.Balancer = &Balancer{} // Name is balancer name

// Balancer is a random balancer.
type Balancer struct{}

// New a random selector.
func New() selector2.Selector {
	return NewBuilder().Build()
}

// Pick is pick a weighted node.
func (p *Balancer) Pick(_ context.Context, nodes []selector2.WeightedNode) (selector2.WeightedNode, selector2.DoneFunc, error) {
	// 随机核心相关逻辑
	if len(nodes) == 0 {
		return nil, nil, selector2.ErrNoAvailable
	}
	cur := rand.Intn(len(nodes))
	selected := nodes[cur]
	d := selected.Pick() // gRPC 回调函数
	return selected, d, nil
}

// NewBuilder returns a selector builder with random balancer
func NewBuilder() selector2.Builder {
	return &selector2.DefaultBuilder{
		Node:     &direct.Builder{}, // 创建空 Node 类型(最重要的是这个类型的方法)
		Balancer: &Builder{},
	}
}

// Builder is random builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() selector2.Balancer {
	return &Balancer{}
}
