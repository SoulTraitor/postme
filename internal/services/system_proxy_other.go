//go:build !windows

package services

// getSystemProxyURL returns platform system proxy settings when available.
//
// macOS/Linux currently rely on http.ProxyFromEnvironment in http_client.go.
func getSystemProxyURL() string {
	return ""
}
