package selector

import (
	"strconv"

	"shop/gmicro/registry"
)

// DefaultNode is selector node
type DefaultNode struct {
	scheme   string
	addr     string
	weight   *int64
	version  string
	name     string
	metadata map[string]string
}

// Scheme is node scheme
func (n *DefaultNode) Scheme() string {
	return n.scheme
}

// Address is node address
func (n *DefaultNode) Address() string {
	return n.addr
}

// ServiceName is node serviceName
func (n *DefaultNode) ServiceName() string {
	return n.name
}

// InitialWeight is node initialWeight
func (n *DefaultNode) InitialWeight() *int64 {
	return n.weight
}

// Version is node version
func (n *DefaultNode) Version() string {
	return n.version
}

// Metadata is node metadata
func (n *DefaultNode) Metadata() map[string]string {
	return n.metadata
}

// NewNode new node
func NewNode(scheme, addr string, ins *registry.ServiceInstance) Node {
	n := &DefaultNode{
		scheme: scheme,
		addr:   addr,
	}
	if ins != nil {
		n.name = ins.Name                          // 服务发现设置的名字
		n.version = ins.Version                    // 服务发现设置的版本号
		n.metadata = ins.Metadata                  // 元数据
		if str, ok := ins.Metadata["weight"]; ok { // 是否 元数据中携带 初始权重
			if weight, err := strconv.ParseInt(str, 10, 64); err == nil {
				n.weight = &weight
			}
		}
	}
	return n
}
