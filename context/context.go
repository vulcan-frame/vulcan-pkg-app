package context

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/pkg/errors"
)

// Context is the context of the game
// Use the custom type for your constants
const (
	CtxSID         = "x-md-global-sid"     // Server ID is the ID for each server in multi-server games, or 0 for single-server games
	CtxUID         = "x-md-global-uid"     // User ID is the ID of the player.It is unique in the game.
	CtxOID         = "x-md-global-oid"     // Object ID for route the message to specific node which has the corresponding module and ID
	CtxColor       = "x-md-global-color"   // Color for route the message to specific node group
	CtxStatus      = "x-md-global-status"  // Status is the status of this connection
	CtxReferer     = "x-md-global-referer" // example: gate:10.0.1.31 or player:10.0.2.31
	CtxClientIP    = "x-md-global-client-ip"
	CtxGateReferer = "x-md-global-gate-referer" // example: 10.0.1.31:9100#10001
)

var Keys = []string{CtxSID, CtxUID, CtxOID, CtxStatus, CtxColor, CtxReferer, CtxClientIP, CtxGateReferer}

func SetColor(ctx context.Context, color string) context.Context {
	return metadata.AppendToClientContext(ctx, string(CtxColor), color)
}

func Color(ctx context.Context) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		return md.Get(CtxColor)
	}
	return ""
}

func SetUID(ctx context.Context, id int64) context.Context {
	return metadata.AppendToClientContext(ctx, CtxUID, strconv.FormatInt(id, 10))
}

func UID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return 0, errors.New("metadata not in context")
	}

	str := md.Get(CtxUID)
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "uid must be int64, uid=%s", str)
	}
	return id, nil
}

func SetOID(ctx context.Context, id int64) context.Context {
	return metadata.AppendToClientContext(ctx, CtxOID, strconv.FormatInt(id, 10))
}

func OID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return 0, errors.New("metadata not in context")
	}

	str := md.Get(CtxOID)
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "oid must be int64, oid=%s", str)
	}
	return id, nil
}

func SetSID(ctx context.Context, id int64) context.Context {
	return metadata.AppendToClientContext(ctx, CtxSID, strconv.FormatInt(id, 10))
}

func SID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return 0, errors.New("metadata not in context")
	}

	str := md.Get(CtxSID)
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "sid must be int64, sid=%s", str)
	}
	return id, nil
}

func SetStatus(ctx context.Context, status int64) context.Context {
	if status == 0 {
		return ctx
	}
	return metadata.AppendToClientContext(ctx, CtxStatus, strconv.FormatInt(status, 10))
}

func Status(ctx context.Context) int64 {
	if md, ok := metadata.FromServerContext(ctx); ok {
		v := md.Get(CtxStatus)
		status, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Errorf("status must be int64, status=%s", v)
			return 0
		}
		return status
	}
	return 0
}

func SetClientIP(ctx context.Context, ip string) context.Context {
	if len(ip) == 0 {
		return ctx
	}
	return metadata.AppendToClientContext(ctx, CtxClientIP, strings.Split(ip, ":")[0])
}

func ClientIP(ctx context.Context) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		return md.Get(CtxClientIP)
	}
	return ""
}

func SetGateReferer(ctx context.Context, server string, wid uint64) context.Context {
	if len(server) == 0 {
		return ctx
	}
	return metadata.AppendToClientContext(ctx, CtxGateReferer, fmt.Sprintf("%s#%d", server, wid))
}

func GateReferer(ctx context.Context) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		return md.Get(CtxGateReferer)
	}
	return ""
}
