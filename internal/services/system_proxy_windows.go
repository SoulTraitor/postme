//go:build windows

package services

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

// getSystemProxyURL reads the Windows system proxy settings from registry.
func getSystemProxyURL() string {
	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer key.Close()

	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil || proxyEnable == 0 {
		return ""
	}

	proxyServer, _, err := key.GetStringValue("ProxyServer")
	if err != nil || proxyServer == "" {
		return ""
	}

	if !strings.HasPrefix(proxyServer, "http://") && !strings.HasPrefix(proxyServer, "https://") {
		proxyServer = "http://" + proxyServer
	}

	return proxyServer
}
