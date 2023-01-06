package http

import (
	"encoding/json"
	_ "github.com/valyala/fasthttp"
)

// Get http get
func Get(url string, body any, requestTimeout int) (string, int, error) {
	return httpRequest("GET", url, nil, body, "application/x-www-form-urlencoded", requestTimeout)
}

// GetFormWithoutBody http get，application/x-www-form-urlencoded，
func GetFormWithoutBody(url string, requestTimeout int) (string, int, error) {
	return httpRequest("GET", url, nil, nil, "application/x-www-form-urlencoded", requestTimeout)
}

// GetJson Post方式将结果反序列化成TReturn
func GetJson[TReturn any](url string, body any, requestTimeout int) (TReturn, error) {
	rspJson, _, err := httpRequest("GET", url, nil, body, "application/json", requestTimeout)
	var val TReturn
	if err == nil {
		_ = json.Unmarshal([]byte(rspJson), &val)
	}
	return val, err
}
