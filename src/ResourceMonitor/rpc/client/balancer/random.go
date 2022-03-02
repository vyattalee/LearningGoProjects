package balancer

import (
	"math/rand"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const Random = "random"

func init() {
	balancer.Register(newRandom())
	rand.Seed(time.Now().UnixNano())
}

func newRandom() balancer.Builder {
	return base.NewBalancerBuilder(Random, &randomPickerBuilder{}, base.Config{})
}

type randomPickerBuilder struct{}

func (*randomPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var conns []balancer.SubConn
	//conns := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for conn := range info.ReadySCs {
		conns = append(conns, conn)
	}

	return &randomPicker{
		subConns: conns,
	}
}

//func (*randomPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
//	if len(readySCs) == 0 {
//		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
//	}
//
//	var conns []balancer.SubConn
//	for _, conn := range readySCs {
//		conns = append(conns, conn)
//	}
//
//	return &randomPicker{
//		subConns: conns,
//	}
//}

type randomPicker struct {
	subConns []balancer.SubConn
}

//func (p *randomPicker) Pick(ctx context.Context, info balancer.PickInfo) (conn balancer.SubConn, done func(balancer.DoneInfo), err error) {
//	conn = p.subConns[rand.Int()%len(p.subConns)]
//	return
//}

func (p *randomPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	//p.mu.Lock()
	//sc := p.subConns[p.next]
	//p.next = (p.next + 1) % len(p.subConns)
	//p.mu.Unlock()
	//return balancer.PickResult{SubConn: sc}, nil
	conn := p.subConns[rand.Int()%len(p.subConns)]
	return balancer.PickResult{SubConn: conn}, nil
}
