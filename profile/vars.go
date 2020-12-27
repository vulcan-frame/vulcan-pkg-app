package profile

import (
	"net/url"
	"strings"
)

var (
	_profile      string
	_color        string
	_version      string
	_grpcEndpoint string
	_nodeName     string
	_zone         uint32
)

func Init(profile, color string, zone uint32, version string, nodeName string, gRPCEndpoint *url.URL) {
	_profile = profile
	_color = color
	_version = version
	_grpcEndpoint = strings.Replace(gRPCEndpoint.String(), "grpc://", "", -1)
	_nodeName = nodeName
	_zone = zone
}

func Profile() string {
	return _profile
}

func Color() string {
	return _color
}

func Version() string {
	return _version
}

func NodeName() string {
	return _nodeName
}

func GRPCEndpoint() string {
	return _grpcEndpoint
}

func Zone() uint32 {
	return _zone
}
