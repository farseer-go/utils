package ws

import (
	"encoding/json"
	"fmt"
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/parse"
	"golang.org/x/net/websocket"
	"net/http"
)

// Client websocket 客户端
type Client struct {
	config        *websocket.Config // 客户端配置
	conn          *websocket.Conn   // 客户端连接
	msgBufferSize int               // 接收消息时的缓冲区大小
}

// NewClient 实例化对象
func NewClient(addr string, msgBufferSize int) (*Client, error) {
	config, err := websocket.NewConfig(addr, addr)
	config.Header = make(http.Header)

	if err != nil {
		return nil, err
	}

	if msgBufferSize == 0 {
		msgBufferSize = 1024
	}

	return &Client{
		config:        config,
		msgBufferSize: msgBufferSize,
	}, nil
}

// SetHeader 设置header
func (receiver *Client) SetHeader(key, value string) {
	receiver.config.Header.Set(key, value)
}

// SetHeaderMap 设置header
func (receiver *Client) SetHeaderMap(m map[string]any) {
	for k, v := range m {
		receiver.config.Header.Set(k, parse.ToString(v))
	}
}

// Connect 连接
func (receiver *Client) Connect() error {
	var err error
	receiver.conn, err = websocket.DialConfig(receiver.config)
	return err
}

// ReceiverJson 接收消息（并反序列化成val）
func (receiver *Client) ReceiverJson(val any) error {
	retMsg := make([]byte, receiver.msgBufferSize)
	n, err := receiver.conn.Read(retMsg)
	if err != nil {
		return err
	}
	return json.Unmarshal(retMsg[:n], val)
}

// Receiver 接收消息
func (receiver *Client) Receiver() (string, error) {
	var retMsg = make([]byte, receiver.msgBufferSize)
	n, err := receiver.conn.Read(retMsg)
	if err != nil {
		return "", err
	}
	return string(retMsg[:n]), err
}

// Send 发送消息，如果msg不是go的基础类型，则会自动序列化成json
func (receiver *Client) Send(msg any) error {
	switch fastReflect.PointerOf(msg).Type {
	case fastReflect.GoBasicType:
		_, err := receiver.conn.Write([]byte(parse.ToString(msg)))
		return err
	default:
		marshalBytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("发送数据时，出现反序列失败：%s", err.Error())
		}
		_, err = receiver.conn.Write(marshalBytes)
		return err
	}
}

// Close 关闭连接
func (receiver *Client) Close() {
	_ = receiver.conn.Close()
}
