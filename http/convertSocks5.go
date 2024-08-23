package http

import "strings"

// ConvertSocks5 如果With.Proxy不是socks5协议，则自动添加socks5协议
func ConvertSocks5(proxy string) string {
	index := strings.Index(proxy, "://")
	if index == -1 {
		return "socks5://" + proxy
	}

	if !strings.HasPrefix(strings.ToLower(proxy), "socks5://") {
		return "socks5" + proxy[index:]
	}
	return proxy
}
