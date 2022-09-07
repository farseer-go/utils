package http

import (
	"encoding/json"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/valyala/fasthttp"
	"time"
)

// 支持请求超时设置，单位：ms
func httpRequest(methodName string, url string, head map[string]any, body any, contentType string, requestTimeout int) (string, error) {
	client := fasthttp.Client{}

	// request
	request := fasthttp.AcquireRequest()

	// url
	request.SetRequestURI(url)

	// request.body
	bytesData, _ := json.Marshal(body)
	request.SetBody(bytesData)

	// request.contentType
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}

	if head != nil || len(head) > 0 {
		for k, v := range head {
			request.Header.Set(k, parse.Convert(v, ""))
		}
	}

	// Method
	request.Header.SetMethod(methodName)

	response := fasthttp.AcquireResponse()
	timeout := time.Duration(requestTimeout) * time.Millisecond
	err := client.DoTimeout(request, response, timeout)

	if err != nil {
		flog.Errorf("%s request error:%s", url, err.Error())
		return "", err
	}
	return string(response.Body()), nil
}
