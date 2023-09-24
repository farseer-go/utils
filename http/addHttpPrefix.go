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

// ClearHttpPrefix 清除http前缀
func ClearHttpPrefix(url string) string {
	var domain string
	if strings.HasPrefix(url, "https") {
		domain = strings.TrimPrefix(url, "https://")
	} else {
		domain = strings.TrimPrefix(url, "http://")
	}

	return strings.TrimRight(domain, "/")
}

// GetDomain 从URL中获取Domain部份
func GetDomain(url string) string {
	domain := strings.TrimPrefix(url, "https://")
	domain = strings.TrimPrefix(url, "https//")
	return strings.Split(domain, "/")[0]
}
