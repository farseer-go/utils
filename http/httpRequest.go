package http

import (
	"encoding/json"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/valyala/fasthttp"
	"net/url"
	"time"
)

// 支持请求超时设置，单位：ms
func httpRequest(methodName string, requestUrl string, head map[string]any, body any, contentType string, requestTimeout int) (string, error) {
	sw := stopwatch.StartNew()

	client := fasthttp.Client{}

	// request
	request := fasthttp.AcquireRequest()

	// requestUrl
	request.SetRequestURI(requestUrl)

	switch b := body.(type) {
	case string:
		request.SetBodyString(b)
	case url.Values:
		request.SetBodyString(urlValuesToString(b, contentType))
	case map[string]string:
		request.SetBodyString(mapStringToString(b, contentType))
	case map[string]any:
		request.SetBodyString(mapAnyToString(b, contentType))
	default:
		// request.body
		bytesData, _ := json.Marshal(body)
		request.SetBody(bytesData)
	}

	// request.contentType
	if contentType != "" {
		request.Header.SetContentType(contentType)
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

	flog.ComponentInfof("httpRequest", "[%s] %s body:%v，耗时：%s", methodName, requestUrl, body, sw.GetMillisecondsText())

	if err != nil {
		return "", flog.Errorf("%s request error:%s", requestUrl, err.Error())
	}
	return string(response.Body()), nil
}

func urlValuesToString(body url.Values, contentType string) string {
	if contentType == "application/json" {
		bytesData, _ := json.Marshal(body)
		return string(bytesData)
	} else {
		return body.Encode()
	}
}

func mapStringToString(body map[string]string, contentType string) string {
	if contentType == "application/json" {
		bytesData, _ := json.Marshal(body)
		return string(bytesData)
	} else {
		val := make(url.Values)
		for k, v := range body {
			val.Add(k, v)
		}
		return val.Encode()
	}
}

func mapAnyToString(body map[string]any, contentType string) string {
	if contentType == "application/json" {
		bytesData, _ := json.Marshal(body)
		return string(bytesData)
	} else {
		val := make(url.Values)
		for k, v := range body {
			val.Add(k, v.(string))
		}
		return val.Encode()
	}
}
