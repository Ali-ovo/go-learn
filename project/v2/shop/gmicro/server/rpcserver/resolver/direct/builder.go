package direct

import (
	"google.golang.org/grpc/resolver"
	"strings"
)

func init() {
	// 解析器
	resolver.Register(NewBuilder())
}

type directBuilder struct{}

// NewBuilder creates a directBuilder which is used to factory direct resolvers.
// example:
//
//	direct://<authority>/127.0.0.1:9000,127.0.0.2:9000
func NewBuilder() resolver.Builder {
	return &directBuilder{}
}

func (d directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := make([]resolver.Address, 0)
	for _, addr := range strings.Split(strings.TrimPrefix(target.URL.Path, "/"), ",") {
		addrs = append(addrs, resolver.Address{Addr: addr})
	}

	// grpc 建立连接的逻辑都在这里 UpdateState
	err := cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		return nil, err
	}
	return newDirectResolver(), nil
}

func (d directBuilder) Scheme() string {
	return "direct"
}
