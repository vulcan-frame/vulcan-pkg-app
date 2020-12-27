package routetable

import (
	"context"
	"time"
)

const (
	defaultTTL = time.Hour * 24 * 7
)

var _ RouteTable = (*BaseRouteTable)(nil)

type getKeyFunc func(name, color string, oid int64) string

type Option func(*BaseRouteTable)

func WithTTL(dur time.Duration) Option {
	return func(r *BaseRouteTable) {
		r.ttl = dur
	}
}

type BaseRouteTable struct {
	RouteTableData

	name   string
	getKey getKeyFunc
	ttl    time.Duration
}

func NewBaseRouteTable(rtd RouteTableData, name string, getKey getKeyFunc, opts ...Option) *BaseRouteTable {
	rt := &BaseRouteTable{
		RouteTableData: rtd,
		name:           name,
		getKey:         getKey,
		ttl:            defaultTTL,
	}
	for _, opt := range opts {
		opt(rt)
	}
	return rt
}

func (r *BaseRouteTable) Store(ctx context.Context, color string, uid int64, addr string) error {
	return r.RouteTableData.Set(ctx, r.getKey(r.name, color, uid), addr, r.ttl)
}

func (r *BaseRouteTable) GetSet(ctx context.Context, color string, uid int64, addr string) (old string, err error) {
	return r.RouteTableData.GetSet(ctx, r.getKey(r.name, color, uid), addr, r.ttl)
}

func (r *BaseRouteTable) SetNx(ctx context.Context, color string, uid int64, addr string) (ok bool, result string, err error) {
	return r.RouteTableData.SetNx(ctx, r.getKey(r.name, color, uid), addr, r.ttl)
}

func (r *BaseRouteTable) Load(ctx context.Context, color string, uid int64) (addr string, err error) {
	return r.RouteTableData.Load(ctx, r.getKey(r.name, color, uid))
}

func (r *BaseRouteTable) LoadAndExpire(ctx context.Context, color string, uid int64) (addr string, err error) {
	return r.RouteTableData.LoadAndExpire(ctx, r.getKey(r.name, color, uid), r.ttl)
}

func (r *BaseRouteTable) Del(ctx context.Context, color string, uid int64) error {
	return r.RouteTableData.Del(ctx, r.getKey(r.name, color, uid))
}

func (r *BaseRouteTable) DelDelay(ctx context.Context, color string, uid int64, expiration time.Duration) error {
	return r.RouteTableData.Expire(ctx, r.getKey(r.name, color, uid), expiration)
}

func (r *BaseRouteTable) DelIfSame(ctx context.Context, color string, uid int64, value string) error {
	return r.RouteTableData.DelIfSame(ctx, r.getKey(r.name, color, uid), value)
}
