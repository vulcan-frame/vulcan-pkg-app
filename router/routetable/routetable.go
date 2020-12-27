package routetable

import (
	"context"
	"fmt"
	"time"
)

type RouteTable interface {
	ReadOnlyRouteTable

	LoadAndExpire(ctx context.Context, color string, oid int64) (addr string, err error)
	Store(ctx context.Context, color string, key int64, addr string) error
	GetSet(ctx context.Context, color string, key int64, addr string) (old string, err error)
	SetNx(ctx context.Context, color string, key int64, addr string) (ok bool, result string, err error)
	DelDelay(ctx context.Context, color string, key int64, delay time.Duration) error
	DelIfSame(ctx context.Context, color string, key int64, value string) error
	Del(ctx context.Context, color string, key int64) error
}

type ReadOnlyRouteTable interface {
	Load(ctx context.Context, color string, key int64) (addr string, err error)
}

type RouteTableData interface {
	Load(ctx context.Context, key string) (addr string, err error)
	LoadAndExpire(ctx context.Context, key string, dur time.Duration) (string, error)
	Set(ctx context.Context, key string, addr string, dur time.Duration) error
	GetSet(ctx context.Context, key string, addr string, dur time.Duration) (old string, err error)
	SetNx(ctx context.Context, key string, addr string, dur time.Duration) (ok bool, result string, err error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	DelIfSame(ctx context.Context, key string, value string) error
	Del(ctx context.Context, key string) error
}

func NewRouteTable(name string, rt RouteTableData, opts ...Option) RouteTable {
	return NewBaseRouteTable(rt, name, key, opts...)
}

func key(name, color string, oid int64) string {
	return fmt.Sprintf("r_%s_{%s}_{%d}", name, color, oid)
}
