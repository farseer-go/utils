package http

import (
	"crypto/tls"
	"fmt"
	"github.com/farseer-go/fs/parse"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"os"
	"time"
)

// Download 下载文件到本地
func Download(url string, savePath string, head map[string]any, requestTimeout int, proxyAddr string) (map[string]string, error) {
	// request
	request := fasthttp.AcquireRequest()

	request.SetRequestURI(url)
	// Method
	request.Header.SetMethod("GET")
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)
	defer request.SetConnectionClose()

	fastHttpClient := fasthttp.Client{
		TLSConfig: &tls.Config{
			// 指定不校验 SSL/TLS 证书
			InsecureSkipVerify: true,
		},
		RetryIf:        func(request *fasthttp.Request) bool { return false },
		ReadBufferSize: 8192,
	}

	// 设置代理
	if proxyAddr != "" {
		fastHttpClient.Dial = fasthttpproxy.FasthttpSocksDialer(proxyAddr)
	}

	if head != nil || len(head) > 0 {
		for k, v := range head {
			request.Header.Set(k, parse.Convert(v, ""))
		}
	}

	// 请求
	var err error
	if requestTimeout > 0 {
		err = fastHttpClient.DoTimeout(request, response, time.Duration(requestTimeout)*time.Millisecond)
	} else {
		err = fastHttpClient.Do(request, response)
	}
	if err != nil {
		return nil, err
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

	statusCode := response.StatusCode()
	if statusCode == 301 || statusCode == 302 || statusCode == 303 {
		location := string(response.Header.Peek("Location"))
		return Download(location, savePath, head, requestTimeout, proxyAddr)
	}

	_ = response.CloseBodyStream()
	if body := response.Body(); statusCode == 200 && len(body) > 0 {
		f, err := os.Create(savePath)
		defer func() {
			_ = f.Close()
			_ = os.Chmod(savePath, 0765)
		}()

		// 创建文件
		if err != nil {
			return responseHeader, fmt.Errorf("文件保存出错，检查目录: %v", err)
		}
		// 保存到文件
		return responseHeader, response.BodyWriteTo(f)
	} else {
		return responseHeader, fmt.Errorf("下载文件失败：%d", statusCode)
	}
}
