package profile

import "strings"

// Profile is the running model of the service
const (
	ProfileDev  = "dev"
	ProfileTest = "test"
	ProfileProd = "prod"
)

func IsDev() bool {
	return IsDevStr(_profile)
}

func IsDevStr(profile string) bool {
	return strings.ToLower(profile) == ProfileDev
}

func IsTestStr(profile string) bool {
	return strings.ToLower(profile) == ProfileTest
}

func IsProdStr(profile string) bool {
	return strings.ToLower(profile) == ProfileProd
}
