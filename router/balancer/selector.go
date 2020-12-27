package balancer

import (
	"context"
	"strconv"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/node/direct"
	"github.com/pkg/errors"
	vctx "github.com/vulcan-frame/vulcan-pkg-app/context"
	"github.com/vulcan-frame/vulcan-pkg-app/router/routetable"
)

// New random a selector.
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

type Option func(o *options)

type options struct {
	balancerType BalancerType
	routeTable   routetable.RouteTable
}

func WithRouteTable(rt routetable.RouteTable) Option {
	return func(o *options) {
		o.routeTable = rt
	}
}
func WithBalancerType(balancerType BalancerType) Option {
	return func(o *options) {
		o.balancerType = balancerType
	}
}

type Builder struct {
	balancerType BalancerType
	routeTable   routetable.RouteTable
}

// NewBuilder returns a selector builder with wrr balancer
func NewBuilder(opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &selector.DefaultBuilder{
		Balancer: &Builder{
			balancerType: option.balancerType,
			routeTable:   option.routeTable,
		},
		Node: &direct.Builder{},
	}
}

func (b *Builder) Build() selector.Balancer {
	return &Balancer{
		balancerType:  b.balancerType,
		currentWeight: make(map[string]float64),
		routeTable:    b.routeTable,
	}
}

var _ selector.Balancer = (*Balancer)(nil)

type Balancer struct {
	balancerType  BalancerType
	mu            sync.Mutex
	currentWeight map[string]float64
	routeTable    routetable.RouteTable
}

// Pick is pick a weighted node
func (p *Balancer) Pick(ctx context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}

	oid, err := getOIDFromCtx(ctx)
	if err != nil {
		return nil, nil, err
	}
	color := getColorFromCtx(ctx)

	// select node by oid from routeTable
	addr, err := p.routeTable.LoadAndExpire(ctx, color, oid)
	if err != nil {
		return nil, nil, err
	}
	for _, node := range nodes {
		if node.Address() == addr {
			return node, nil, nil
		}
	}

	// select a new node by weight from nodes
	// the algorithm is the implement of nginx wrr, copied from https://github.com/go-kratos/kratos/blob/main/selector/wrr/wrr.go
	var (
		totalWeight  float64
		selected     selector.WeightedNode
		selectWeight float64
	)

	p.mu.Lock()
	for _, node := range nodes {
		totalWeight += node.Weight()
		cwt := p.currentWeight[node.Address()]
		cwt += node.Weight()
		p.currentWeight[node.Address()] = cwt
		if selected == nil || selectWeight < cwt {
			selectWeight = cwt
			selected = node
		}
	}
	p.currentWeight[selected.Address()] = selectWeight - totalWeight
	p.mu.Unlock()

	d := selected.Pick()

	if p.balancerType != BalancerTypeMaster {
		return selected, d, nil
	}

	// update route table if the connector is master
	// the route table may be set by other connections, so we need to judge it as empty before setting
	ok, addr, err := p.routeTable.SetNx(ctx, color, oid, selected.Address())
	if err != nil {
		return nil, nil, err
	}
	if ok {
		// the route table is set by this connection
		return selected, d, nil
	}

	log.Warnf("routeTable is set by other connections. oid=%d color=%s oldConn=%s newConn=%s", oid, color, addr, selected.Address())
	for _, node := range nodes {
		if node.Address() == addr {
			return node, nil, nil
		}
	}
	return nil, nil, errors.Errorf("the existed connection in routeTable is not found. oid=%d color=%s oldConn=%s", oid, color, addr)
}

func getOIDFromCtx(ctx context.Context) (oid int64, err error) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		err = errors.Errorf("metadata not in context")
		return
	}
	if oid, err = strconv.ParseInt(md.Get(vctx.CtxOID), 10, 64); err != nil {
		err = errors.Wrapf(err, "oid not int64")
	}
	return
}
