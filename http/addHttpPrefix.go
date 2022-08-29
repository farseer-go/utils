package http

import "strings"

// AddHttpPrefix 添加http前缀
func AddHttpPrefix(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}
	return "http://" + url
}

// AddHttpsPrefix 添加https前缀
func AddHttpsPrefix(url string) string {
	if strings.HasPrefix(url, "https") {
		return url
	}
	return "https://" + url
}
