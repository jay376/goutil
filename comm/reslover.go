package comm

import (
	"fmt"
	"sync"

	"google.golang.org/grpc/resolver"
)

var defaultScheme = "infra"
var rwMu sync.RWMutex
var addrStore map[string][]string

func init() {
	addrStore = make(map[string][]string)
	resolver.Register(&StaticResloverBuilder{})
}

// StaticResloverBuilder ...
type StaticResloverBuilder struct {
}

// Scheme ...
func (b *StaticResloverBuilder) Scheme() string {
	return defaultScheme
}

// Build resolver
func (b *StaticResloverBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOption) (resolver.Resolver, error) {
	r := &StaticResolver{
		target: target,
		cc:     cc,
	}
	r.start()
	return r, nil
}

// StaticResolver ...
type StaticResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
}

// ResolveNow ...
func (r *StaticResolver) ResolveNow(o resolver.ResolveNowOption) {}

// Close ...
func (r *StaticResolver) Close() {}

func (r *StaticResolver) start() {
	s := resolver.State{
		Addresses: make([]resolver.Address, 0),
	}
	rwMu.RLock()
	defer rwMu.RUnlock()

	for _, addr := range addrStore[r.target.Endpoint] {
		s.Addresses = append(s.Addresses, resolver.Address{
			Addr: addr,
		})
	}
	r.cc.UpdateState(s)
}

// SetEndpointAddr ...
func SetEndpointAddr(endpoint string, addr []string) {
	rwMu.Lock()
	addrStore[endpoint] = addr
	rwMu.Unlock()
}

// GetStaticTarget get spserver nameing
func GetStaticTarget(endpoint string) string {
	return fmt.Sprintf("%s:///%s", defaultScheme, endpoint)
}

// // GetServiceTarget get service nameing
// func GetDynamicTarget(serviceName string) string {
// 	return fmt.Sprintf("%s:///%s", defaultScheme, serviceName)
// }

// // DynamicResolver should implements subscribe spserver service, so set addr with cc.UpdateState
// type DynamicResolver struct {
// 	target resolver.Target
// 	cc     resolver.ClientConn
// }

// // ResolveNow ...
// func (r *DynamicResolver) ResolveNow(o resolver.ResolveNowOption) {}

// // Close ...
// func (r *DynamicResolver) Close() {}
