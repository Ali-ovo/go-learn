package selector

import (
	"context"
	"time"
)

// Balancer is balancer interface
type Balancer interface {
	Pick(ctx context.Context, nodes []WeightedNode) (selected WeightedNode, done DoneFunc, err error)
}

// BalancerBuilder build balancer
type BalancerBuilder interface {
	Build() Balancer
}

// WeightedNode calculates scheduling weight in real time
type WeightedNode interface {
	Node

	// Raw returns the original node
	// Raw 返回 原始的 node
	Raw() Node

	// Weight is the runtime calculated weight
	// Weight 是 运行计算时的权重
	Weight() float64

	// Pick the node
	Pick() DoneFunc

	// PickElapsed is time elapsed since the latest pick
	PickElapsed() time.Duration
}

// WeightedNodeBuilder is WeightedNode Builder
type WeightedNodeBuilder interface {
	Build(Node) WeightedNode
}
