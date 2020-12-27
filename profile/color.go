package profile

import "strings"

// Color is the color of the service, the message will be routed to the corresponding color node
const (
	ColorLocal = "local"
)

func IsLocal() bool {
	return strings.ToLower(_color) == ColorLocal
}
