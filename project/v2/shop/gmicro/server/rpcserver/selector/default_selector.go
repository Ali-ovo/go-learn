package selector

import (
	"context"
	"sync/atomic"
)

// Default is composite selector.
type Default struct {
	NodeBuilder WeightedNodeBuilder
	Balancer    Balancer
	// 原子性
	nodes atomic.Value
}

// Select is select one node.
// 选择出一个节点
func (d *Default) Select(ctx context.Context) (selected Node, done DoneFunc, err error) {
	var (
		candidates []WeightedNode
	)
	nodes, ok := d.nodes.Load().([]WeightedNode) // 读取数据
	if !ok {
		return nil, nil, ErrNoAvailable
	}
	candidates = nodes

	if len(candidates) == 0 {
		return nil, nil, ErrNoAvailable
	}
	// 这里调用 如 p2c的Pick 或 random的Pick
	// done 是 gRPC 执行完成一次 访问后 回调函数 可以访问调用的结果和元数据
	wn, done, err := d.Balancer.Pick(ctx, candidates)
	if err != nil {
		return nil, nil, err
	}
	// TODO p ? 作用: context.Context 中 存储 node 信息
	p, ok := FromPeerContext(ctx)
	if ok {
		p.Node = wn.Raw()
	}
	return wn.Raw(), done, nil
}

// Apply update nodes info.
// nodes 相关信息 和 子连接 的数组
func (d *Default) Apply(nodes []Node) {
	weightedNodes := make([]WeightedNode, 0, len(nodes))
	for _, n := range nodes {
		weightedNodes = append(weightedNodes, d.NodeBuilder.Build(n)) // 重新分装 node
	}
	// TODO: Do not delete unchanged nodes
	// ServiceInstance 和 Addr 存储在 d.nodes 中
	d.nodes.Store(weightedNodes)
}

// DefaultBuilder is de
type DefaultBuilder struct {
	Node     WeightedNodeBuilder
	Balancer BalancerBuilder
}

// Build create builder
func (db *DefaultBuilder) Build() Selector {
	return &Default{
		NodeBuilder: db.Node,
		// 这里执行 random 相关 Build 逻辑
		Balancer: db.Balancer.Build(), // 生成 random.Balancer{} 空结构体 可以调用 random 的Pick
	}
}
