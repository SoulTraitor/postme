//go:build !windows

package services

func getSystemProxyURL() string {
	return ""
}
