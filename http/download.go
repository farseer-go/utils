package http

import (
	"crypto/tls"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"os"
	"time"
)

// Download 下载文件到本地
func Download(url string, savePath string, requestTimeout int, proxyAddr string) error {
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
		RetryIf: func(request *fasthttp.Request) bool { return false },
	}

	// 设置代理
	if proxyAddr != "" {
		fastHttpClient.Dial = fasthttpproxy.FasthttpSocksDialer(proxyAddr)
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")

	// 请求
	var err error
	if requestTimeout > 0 {
		err = fastHttpClient.DoTimeout(request, response, time.Duration(requestTimeout)*time.Millisecond)
	} else {
		err = fastHttpClient.Do(request, response)
	}
	if err != nil {
		return err
	}

	statusCode := response.StatusCode()
	if statusCode == 301 || statusCode == 302 || statusCode == 303 {
		location := string(response.Header.Peek("Location"))
		return Download(location, savePath, requestTimeout, proxyAddr)
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
			return fmt.Errorf("文件保存出错，检查目录: %v", err)
		}
		// 保存到文件
		return response.BodyWriteTo(f)
	} else {
		return fmt.Errorf("下载文件失败：%d", statusCode)
	}
}
