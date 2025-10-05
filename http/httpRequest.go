package http

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/fs/trace"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// RequestProxy 支持请求超时设置，单位：ms
func RequestProxy(methodName string, requestUrl string, head map[string]any, body any, contentType string, requestTimeout int, proxyAddr string) (string, int, map[string]string, error) {
	return tryRequestProxy(methodName, requestUrl, head, body, contentType, requestTimeout, proxyAddr, 1)
}

func RequestProxyConfigure(methodName string, requestUrl string, head map[string]any, body any, contentType string, requestTimeout int) (string, int, map[string]string, error) {
	return tryRequestProxy(methodName, requestUrl, head, body, contentType, requestTimeout, configure.GetString("Proxy"), 1)
}

// tryRequestProxy 支持请求超时设置，单位：ms
func tryRequestProxy(methodName string, requestUrl string, head map[string]any, body any, contentType string, requestTimeout int, proxyAddr string, tryCount int) (string, int, map[string]string, error) {
	if tryCount > 3 {
		return "", 0, nil, fmt.Errorf("已超过最大尝试次数")
	}

	traceDetail := container.Resolve[trace.IManager]().TraceHttp(methodName, requestUrl)

	// request
	request := fasthttp.AcquireRequest()

	var bodyVal string

	if body != nil {
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
		case *bytes.Buffer:
			bodyVal = (body.(*bytes.Buffer)).String()
		default:
			bytesData, _ := snc.Marshal(body)
			bodyVal = string(bytesData)
		}

		if strings.ToUpper(methodName) == "GET" {
			reqUrl, _ := url.Parse(requestUrl)
			if len(reqUrl.RawQuery) > 0 {
				reqUrl.RawQuery += "&" + bodyVal
			} else {
				reqUrl.RawQuery = bodyVal
			}
			bodyVal = ""
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
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		if head == nil {
			head = make(map[string]any)
		}
		head["Trace-Id"] = traceContext.TraceId
		head["Trace-Level"] = traceContext.TraceLevel
		head["Trace-App-Name"] = core.AppName
		head["Accept-Encoding"] = ""
	}

	if head != nil || len(head) > 0 {
		for k, v := range head {
			request.Header.Set(k, parse.Convert(v, ""))
		}
	}

	// 支持压缩的格式
	if head == nil || head["Accept-Encoding"] == "" {
		request.Header.Set("Accept-Encoding", "gzip, deflate")
	}

	// Method
	request.Header.SetMethod(methodName)
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)
	//defer request.SetConnectionClose()

	fastHttpClient := fasthttp.Client{
		TLSConfig: &tls.Config{
			// 指定不校验 SSL/TLS 证书
			InsecureSkipVerify: true,
		},
		RetryIf: func(request *fasthttp.Request) bool { return false },
	}
	// 设置代理
	if proxyAddr != "" {
		proxyAddr = ConvertSocks5(proxyAddr)
		fastHttpClient.Dial = fasthttpproxy.FasthttpSocksDialer(proxyAddr)
	}
	var err error
	if requestTimeout > 0 {
		err = fastHttpClient.DoTimeout(request, response, time.Duration(requestTimeout)*time.Millisecond)
	} else {
		err = fastHttpClient.Do(request, response)
	}

	if err != nil {
		return "", 0, nil, err
	}

	// 响应头部
	responseHeader := make(map[string]string)
	for _, keyBytes := range response.Header.PeekKeys() {
		responseHeader[string(keyBytes)] = string(response.Header.PeekBytes(keyBytes))
	}
	// cookies
	response.Header.VisitAllCookie(func(key, value []byte) {
		responseHeader[string(key)] = string(value)
	})
	responseHeader["Content-Type"] = string(response.Header.ContentType())

	// 返回的body内容
	var responseBytes []byte
	// 解压缩
	responseContentEncoding := string(response.Header.ContentEncoding())
	switch responseContentEncoding {
	case "gzip":
		bodyReader, _ := gzip.NewReader(bytes.NewReader(response.Body()))
		responseBytes, _ = ioutil.ReadAll(bodyReader)
	case "deflate":
		bodyReader := flate.NewReader(bytes.NewReader(response.Body()))
		responseBytes, _ = ioutil.ReadAll(bodyReader)
	default:
		responseBytes = response.Body()
	}

	// 找到对应的响应编码
	e, name, certain := charset.DetermineEncoding(responseBytes, contentType)

	//charset := ""
	//for _, ctypes := range strings.Split(responseHeader["Content-Type"], ";") {
	//	ctype := strings.Split(ctypes, "=")
	//	if strings.TrimSpace(ctype[0]) == "charset" {
	//		charset = strings.TrimSpace(ctype[1])
	//		break
	//	}
	//}

	var bodyContent string
	switch name { // strings.ToLower(charset)
	case "big5":
		responseBytes, _ = traditionalchinese.Big5.NewDecoder().Bytes(responseBytes)
	default:
		if !certain || name != "utf-8" {
			// 使用新的decder来解析网页内容
			bodyReader := transform.NewReader(bytes.NewReader(responseBytes), e.NewDecoder())
			responseBytes, _ = io.ReadAll(bodyReader)
		}
	}

	bodyContent = string(responseBytes)

	// 链路追踪设置出入参
	traceDetail.TraceDetailHttp.SetHttpRequest(requestUrl, head, responseHeader, bodyVal, bodyContent, response.StatusCode())
	defer func() { traceDetail.End(err) }()

	if err != nil {
		return "", 0, nil, err
	}

	// 30X跳转
	if location := responseHeader["Location"]; bodyContent == "" && location != "" {
		if !strings.HasPrefix(location, "http") {
			uri := request.URI()
			uri.SetPath(filepath.Dir(string(uri.Path())) + "/" + location)
			location = uri.String()
		}
		switch response.StatusCode() {
		case 301, 302:
			return tryRequestProxy("GET", location, head, body, contentType, requestTimeout, proxyAddr, tryCount+1)
		case 307, 308:
			return tryRequestProxy(methodName, location, head, body, contentType, requestTimeout, proxyAddr, tryCount+1)
		}
	}

	return bodyContent, response.StatusCode(), responseHeader, nil
}

func urlValuesToString(body url.Values, contentType string) string {
	if contentType == "application/json" {
		bytesData, _ := snc.Marshal(body)
		return string(bytesData)
	} else {
		return body.Encode()
	}
}

func mapStringToString(body map[string]string, contentType string) string {
	if contentType == "application/json" {
		bytesData, _ := snc.Marshal(body)
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
		bytesData, _ := snc.Marshal(body)
		return string(bytesData)
	} else {
		val := make(url.Values)
		for k, v := range body {
			val.Add(k, parse.Convert(v, ""))
		}
		return val.Encode()
	}
}
