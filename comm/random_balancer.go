package comm

import (
	"context"
	"math/rand"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

// Random ...
const Random = "RandomLB"

func init() {
	balancer.Register(base.NewBalancerBuilderWithConfig(Random, RandomPickerBuilder{}, base.Config{HealthCheck: true}))
}

// RandomPickerBuilder ...
type RandomPickerBuilder struct{}

// Build ...
func (b RandomPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	scs := make([]balancer.SubConn, 0, len(readySCs))
	for _, sc := range readySCs {
		scs = append(scs, sc)
	}
	return &RandomPicker{scs: scs}
}

// RandomPicker ...
type RandomPicker struct {
	sync.Mutex
	scs []balancer.SubConn
}

// Pick ...
func (r *RandomPicker) Pick(ctx context.Context, opts balancer.PickOptions) (conn balancer.SubConn, done func(balancer.DoneInfo), err error) {
	r.Lock()
	defer r.Unlock()
	conn = r.scs[rand.Intn(len(r.scs))]
	return
}
