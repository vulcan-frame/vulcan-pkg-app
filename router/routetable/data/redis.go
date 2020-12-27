package raw

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	verrors "github.com/vulcan-frame/vulcan-pkg-app/errors"
	"github.com/vulcan-frame/vulcan-pkg-app/router/routetable"
)

const (
	defaultTimeout = 2 * time.Second
	errPrefix      = "redis routeTable"
)

var (
	delIfSameScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 1
end`)
)

type RedisCmdable interface {
	redis.Cmdable
}

var _ routetable.RouteTableData = (*RedisRouteTable)(nil)

type Option func(*RedisRouteTable)

func WithTimeout(dur time.Duration) Option {
	return func(r *RedisRouteTable) {
		r.timeout = dur
	}
}

type RedisRouteTable struct {
	rdb     RedisCmdable
	timeout time.Duration
}

func NewRedisRouteTable(rdb RedisCmdable, opts ...Option) *RedisRouteTable {
	rt := &RedisRouteTable{
		rdb:     rdb,
		timeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(rt)
	}
	return rt
}

func wrapErr(err error, operation string, args ...interface{}) error {
	if errors.Is(err, redis.Nil) {
		return errors.Wrapf(verrors.ErrRouteTableNotFound,
			"%s data not found", operation)
	}
	return errors.Wrapf(err, "%s %s failed %v", errPrefix, operation, args)
}

func (rt *RedisRouteTable) Set(ctx context.Context, key string, addr string, dur time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, rt.timeout)
	defer cancel()

	if err := rt.rdb.SetEx(ctx, key, addr, dur).Err(); err != nil {
		return wrapErr(err, "Set", "key", key, "addr", addr)
	}
	return nil
}

func (rt *RedisRouteTable) GetSet(ctx context.Context, key string, addr string, dur time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, rt.timeout)
	defer cancel()

	var old string
	cmders, err := rt.rdb.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.GetSet(ctx, key, addr)
		pipeliner.Expire(ctx, key, dur)
		return nil
	})

	if err := wrapErr(err, "GetSet", "key", key, "addr", addr); err != nil {
		return "", err
	}

	for _, cmder := range cmders {
		if cmd, ok := cmder.(*redis.StringCmd); ok && cmd.Name() == "getset" {
			old = cmd.Val()
			break
		}
	}
	return old, nil
}

// SetNx sets the value if not exists with expiration, returns:
// ok - true when key was set
// result - current value (new value when ok=true)
// err - operation error
func (rt *RedisRouteTable) SetNx(ctx context.Context, key string, addr string, dur time.Duration) (bool, string, error) {
	ctx, cancel := context.WithTimeout(ctx, rt.timeout)
	defer cancel()

	cmds, err := rt.rdb.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.SetNX(ctx, key, addr, dur)
		pipeliner.Get(ctx, key)
		return nil
	})
	if err != nil {
		return false, "", wrapErr(err, "SetNx", "key", key, "addr", addr)
	}

	if len(cmds) != 2 {
		return false, "", wrapErr(errors.New("redis pipeline failed"), "SetNx", "key", key)
	}

	var ok bool
	if setCmd, okCmd := cmds[0].(*redis.BoolCmd); okCmd {
		var errSet error
		ok, errSet = setCmd.Result()
		if errSet != nil {
			return false, "", wrapErr(errSet, "SetNx", "key", key)
		}
	} else {
		return false, "", wrapErr(errors.Errorf("unexpected SETNX response type: %T", cmds[0]), "SetNx", "key", key)
	}

	var currentValue string
	if getCmd, okCmd := cmds[1].(*redis.StringCmd); okCmd {
		val, errGet := getCmd.Result()
		if errGet != nil && !errors.Is(errGet, redis.Nil) {
			return false, "", wrapErr(errGet, "SetNx", "key", key)
		}
		currentValue = val
	} else {
		return false, "", wrapErr(errors.Errorf("unexpected GET response type: %T", cmds[1]), "SetNx", "key", key)
	}

	return ok, currentValue, nil
}

func (rt *RedisRouteTable) Load(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, rt.timeout)
	defer cancel()

	result, err := rt.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", wrapErr(verrors.ErrRouteTableNotFound, "Load", "key", key)
		}
		return result, wrapErr(err, "Load", "key", key)
	}
	return result, nil
}

func (rt *RedisRouteTable) LoadAndExpire(ctx context.Context, key string, dur time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, rt.timeout)
	defer cancel()

	result, err := rt.rdb.GetEx(ctx, key, dur).Result()
	if err != nil {
		return "", wrapErr(err, "LoadAndExpire", "key", key)
	}
	return result, nil
}

func (rt *RedisRouteTable) Del(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, rt.timeout)
	defer cancel()

	if err := rt.rdb.Del(ctx, key).Err(); err != nil {
		return wrapErr(err, "Del", "key", key)
	}
	return nil
}

func (rt *RedisRouteTable) DelIfSame(ctx context.Context, key string, value string) error {
	ctx, cancel := context.WithTimeout(ctx, rt.timeout)
	defer cancel()

	result, err := delIfSameScript.Run(ctx, rt.rdb, []string{key}, value).Int64()
	if err != nil {
		return wrapErr(err, "DelIfSame", "key", key, "value", value)
	}

	if result == 0 {
		return wrapErr(errors.New("redis script execute failed"), "DelIfSame", "key", key, "value", value)
	}
	return nil
}

func (rt *RedisRouteTable) Expire(ctx context.Context, key string, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, rt.timeout)
	defer cancel()

	if err := rt.rdb.Expire(ctx, key, expiration).Err(); err != nil {
		return wrapErr(err, "Expire", "key", key)
	}
	return nil
}
