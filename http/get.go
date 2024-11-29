package http

import (
	"github.com/farseer-go/fs/snc"
	_ "github.com/valyala/fasthttp"
)

// Get http get
func Get(url string, body any, requestTimeout int) (string, int, error) {
	rspBody, statusCode, _, err := RequestProxyConfigure("GET", url, nil, body, "application/x-www-form-urlencoded", requestTimeout)
	return rspBody, statusCode, err
}

// GetFormWithoutBody http get，application/x-www-form-urlencoded，
func GetFormWithoutBody(url string, requestTimeout int) (string, int, error) {
	rspBody, statusCode, _, err := RequestProxyConfigure("GET", url, nil, nil, "application/x-www-form-urlencoded", requestTimeout)
	return rspBody, statusCode, err
}

// GetJson Post方式将结果反序列化成TReturn
func GetJson[TReturn any](url string, body any, requestTimeout int) (TReturn, error) {
	rspJson, _, _, err := RequestProxyConfigure("GET", url, nil, body, "application/json", requestTimeout)
	var val TReturn
	if err == nil {
		_ = snc.Unmarshal([]byte(rspJson), &val)
	}
	return val, err
}
