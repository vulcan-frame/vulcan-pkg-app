package router

import "time"

const (
	AppTunnelChangeTimeout = time.Second * 3
	HolderCacheTimeout     = time.Second * 5
	AsyncRouteTableTimeout = time.Second * 1
)
