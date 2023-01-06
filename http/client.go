package http

import (
	"encoding/json"
	"github.com/farseer-go/fs/flog"
)

type client struct {
	url            string
	head           map[string]any
	body           any
	contentType    string
	requestTimeout int
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
func (receiver *client) Post() (string, error) {
	return httpRequest("POST", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout)
}

// PostUnmarshal POST请求，并反序列成对象
func (receiver *client) PostUnmarshal(val any) error {
	rspJson, err := httpRequest("POST", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout)
	if err == nil {
		err = json.Unmarshal([]byte(rspJson), &val)
		if err != nil {
			return flog.Errorf("%s http.PostUnmarshal error:%s", receiver.url, err.Error())
		}
	}
	return err
}

// Get GET方法请求
func (receiver *client) Get() (string, error) {
	return httpRequest("GET", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout)
}

// GetUnmarshal GET方法请求，并反序列成对象
func (receiver *client) GetUnmarshal(val any) error {
	rspJson, err := httpRequest("GET", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout)
	if err == nil {
		err = json.Unmarshal([]byte(rspJson), &val)
		if err != nil {
			return flog.Errorf("%s http.PostUnmarshal error:%s", receiver.url, err.Error())
		}
	}
	return err
}

// Put PUT方法请求
func (receiver *client) Put() (string, error) {
	return httpRequest("PUT", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout)
}

// PutUnmarshal PUT方法请求，并反序列成对象
func (receiver *client) PutUnmarshal(val any) error {
	rspJson, err := httpRequest("PUT", receiver.url, receiver.head, receiver.body, receiver.contentType, receiver.requestTimeout)
	if err == nil {
		err = json.Unmarshal([]byte(rspJson), &val)
		if err != nil {
			return flog.Errorf("%s http.PostUnmarshal error:%s", receiver.url, err.Error())
		}
	}
	return err
}
