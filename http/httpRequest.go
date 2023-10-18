package http

import (
	"encoding/json"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/linkTrace"
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
	"time"
)

// 支持请求超时设置，单位：ms
func httpRequest(methodName string, requestUrl string, head map[string]any, body any, contentType string, requestTimeout int) (string, int, error) {
	traceDetailHttp := linkTrace.TraceHttp(methodName, requestUrl)

	client := fasthttp.Client{}
	client.RetryIf = func(request *fasthttp.Request) bool {
		return false
	}
	// request
	request := fasthttp.AcquireRequest()

	if body != nil {
		var bodyVal string
		switch b := body.(type) {
		case string:
			bodyVal = b
		case []byte:
			bodyVal = string(b)
		case url.Values:
			bodyVal = urlValuesToString(b, contentType)
		case map[string]string:
			bodyVal = mapStringToString(b, contentType)
		case map[string]any:
			bodyVal = mapAnyToString(b, contentType)
		default:
			bytesData, _ := json.Marshal(body)
			bodyVal = string(bytesData)
		}

		if strings.ToUpper(methodName) == "GET" {
			reqUrl, _ := url.Parse(requestUrl)
			if len(reqUrl.RawQuery) > 0 {
				reqUrl.RawQuery += "&" + bodyVal
			} else {
				reqUrl.RawQuery = bodyVal
			}
			requestUrl = reqUrl.String()
		} else {
			request.SetBodyString(bodyVal)
		}
	}

	request.SetRequestURI(requestUrl)

	// request.contentType
	if contentType != "" {
		request.Header.SetContentType(contentType)
	}

	// 链路追踪
	if trace := linkTrace.GetCurTrace(); trace != nil {
		if head == nil {
			head = make(map[string]any)
		}
		head["TraceId"] = trace.TraceId
	}
	if head != nil || len(head) > 0 {
		for k, v := range head {
			request.Header.Set(k, parse.Convert(v, ""))
		}
	}

	// Method
	request.Header.SetMethod(methodName)
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)
	defer request.SetConnectionClose()

	timeout := time.Duration(requestTimeout) * time.Millisecond
	err := client.DoTimeout(request, response, timeout)
	defer func() { traceDetailHttp.End(err) }()

	if err != nil {
		return "", 0, err
	}
	return string(response.Body()), response.StatusCode(), nil
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
			val.Add(k, parse.Convert(v, ""))
		}
		return val.Encode()
	}
}
