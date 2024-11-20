package http

import (
	"github.com/bytedance/sonic"
	"github.com/farseer-go/fs/flog"
)

type client struct {
	url            string
	head           map[string]any
	body           any
	contentType    string
	requestTimeout int
	proxy          string
}

// NewClient 创建Client客户端
func NewClient(url string) *client {
	return &client{
		url:            url,
		head:           map[string]any{},
		contentType:    "application/json",
		requestTimeout: 500,
	}
}

// Head 设置头部
// head：头部
func (receiver *client) Head(head map[string]any) *client {
	for k, v := range head {
		receiver.head[k] = v
	}
	return receiver
}

// HeadAdd 添加头部
// key：key
// value：Value
func (receiver *client) HeadAdd(key string, value any) *client {
	receiver.head[key] = value
	return receiver
}

// Body 设置Body
// body：提交的内容
func (receiver *client) Body(body any) *client {
	receiver.body = body
	return receiver
}

// Proxy 设置Proxy
func (receiver *client) Proxy(proxy string) *client {
	receiver.proxy = proxy
	return receiver
}

// SetJsonType 设置application/json
func (receiver *client) SetJsonType() *client {
	receiver.contentType = "application/json"
	return receiver
}

// SetFormType 设置application/x-www-form-urlencoded
func (receiver *client) SetFormType() *client {
	receiver.contentType = "application/x-www-form-urlencoded"
	return receiver
}

// Timeout 设置超时
// requestTimeout：超时时间（ms）
func (receiver *client) Timeout(requestTimeout int) *client {
	receiver.requestTimeout = requestTimeout
	return receiver
}

// Post POST请求
func (receiver *client) Post() (string, int, map[string]string, error) {
	return RequestProxy("POST", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout, receiver.proxy)
}

// PostUnmarshal POST请求，并反序列成对象
func (receiver *client) PostUnmarshal(val any) (int, error) {
	rspJson, statusCode, _, err := RequestProxy("POST", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout, receiver.proxy)

	if statusCode >= 400 {
		flog.Warningf("%s %d http.PostUnmarshal", receiver.url, statusCode)
		return statusCode, err
	}

	if err == nil {
		err = sonic.Unmarshal([]byte(rspJson), &val)
		if err != nil {
			flog.Warningf("%s http.PostUnmarshal error:%s", receiver.url, err.Error())
			return statusCode, err
		}
	}
	return statusCode, err
}

// Get GET方法请求
func (receiver *client) Get() (string, int, map[string]string, error) {
	return RequestProxy("GET", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout, receiver.proxy)
}

// GetUnmarshal GET方法请求，并反序列成对象
func (receiver *client) GetUnmarshal(val any) (int, error) {
	rspJson, statusCode, _, err := RequestProxy("GET", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout, receiver.proxy)

	if statusCode >= 400 {
		flog.Warningf("%s %d http.GetUnmarshal", receiver.url, statusCode)
	}

	if err == nil {
		err = sonic.Unmarshal([]byte(rspJson), &val)
		if err != nil {
			flog.Warningf("%s http.GetUnmarshal error:%s", receiver.url, err.Error())
			return statusCode, err
		}
	}
	return statusCode, err
}

// Put PUT方法请求
func (receiver *client) Put() (string, int, map[string]string, error) {
	return RequestProxy("PUT", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout, receiver.proxy)
}

// PutUnmarshal PUT方法请求，并反序列成对象
func (receiver *client) PutUnmarshal(val any) (int, error) {
	rspJson, statusCode, _, err := RequestProxy("PUT", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout, receiver.proxy)
	if statusCode >= 400 {
		flog.Warningf("%s %d http.PutUnmarshal", receiver.url, statusCode)
	}

	if err == nil {
		err = sonic.Unmarshal([]byte(rspJson), &val)
		if err != nil {
			flog.Warningf("%s http.PutUnmarshal error:%s", receiver.url, err.Error())
			return statusCode, err
		}
	}
	return statusCode, err
}
