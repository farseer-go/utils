package http

import (
	"encoding/json"
	"github.com/farseer-go/fs/flog"
)

// Post http post，支持请求超时设置，单位：ms
func Post(url string, head map[string]any, body any, contentType string, requestTimeout int) (string, int, error) {
	rspBody, statusCode, _, err := RequestProxyConfigure("POST", url, head, body, contentType, requestTimeout)
	return rspBody, statusCode, err
}

// PostForm http post，application/x-www-form-urlencoded
func PostForm(url string, head map[string]any, body any, requestTimeout int) (string, int, error) {
	rspBody, statusCode, _, err := RequestProxyConfigure("POST", url, head, body, "application/x-www-form-urlencoded", requestTimeout)
	return rspBody, statusCode, err
}

// PostFormWithoutBody http post，application/x-www-form-urlencoded
func PostFormWithoutBody(url string, head map[string]any, requestTimeout int) (string, int, error) {
	rspBody, statusCode, _, err := RequestProxyConfigure("POST", url, head, nil, "application/x-www-form-urlencoded", requestTimeout)
	return rspBody, statusCode, err
}

// PostJson Post方式将结果反序列化成TReturn
func PostJson[TReturn any](url string, head map[string]any, body any, requestTimeout int) (TReturn, int, error) {
	var val TReturn
	rspJson, statusCode, _, err := RequestProxyConfigure("POST", url, head, body, "application/json", requestTimeout)
	if err == nil {
		err = json.Unmarshal([]byte(rspJson), &val)
		if err != nil {
			_ = flog.Errorf("%s json.Unmarshal error:%s", url, err.Error())
		}
	}
	return val, statusCode, err
}
