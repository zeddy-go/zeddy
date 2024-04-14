package service

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

func init() {
	balancer.Register(newBuilder())
}

// Name is the name of hash_round_robin balancer.
const Name = "hash_round_robin"

var grpcLogger = grpclog.Component("hashroundrobin")

// newBuilder creates a new roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &cPickerBuilder{}, base.Config{HealthCheck: true})
}

type cPickerBuilder struct {
}

func (c cPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	grpcLogger.Infof("hashPicker: Build called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &cPicker{
		subConns: scs,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		next: uint32(r.Intn(len(scs))),
	}
}

type cPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn
	next     uint32
}

func (p *cPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	md, _ := metadata.FromOutgoingContext(info.Ctx)
	hashValues := md.Get("Customer-Hash-Value")
	subConnsLen := uint32(len(p.subConns))
	var sc balancer.SubConn
	if len(hashValues) > 0 {
		hashValue := hashValues[0]
		i, err := strconv.ParseUint(hashValue, 10, 32)
		if err != nil {
			i = strSum(hashValue)
		}
		sc = p.subConns[uint32(i)%subConnsLen]
	} else {
		nextIndex := atomic.AddUint32(&p.next, 1)
		sc = p.subConns[nextIndex%subConnsLen]
	}

	return balancer.PickResult{SubConn: sc}, nil
}

func strSum(str string) (sum uint64) {
	for _, c := range str {
		sum += uint64(c)
	}
	return
}
