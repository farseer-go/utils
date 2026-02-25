package ws

import (
	ctx "context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
	"golang.org/x/net/websocket"
)

// Client websocket 客户端
type Client struct {
	config   *websocket.Config // 客户端配置
	conn     *websocket.Conn   // 客户端连接
	isClose  bool              // 是否断开连接
	Ctx      ctx.Context       // 用于通知应用端是否断开连接
	cancel   ctx.CancelFunc    // 用于通知Ctx，连接已断开
	AutoExit bool              // 当断开连接时，自动退出（抛出异常）
	mu       sync.Mutex        // 💡 解決併發 Send 的衝突
}

// NewClient 实例化对象
func NewClient(addr string, autoExit bool) (*Client, error) {
	config, err := websocket.NewConfig(addr, addr)
	if err != nil {
		return nil, err
	}
	config.Header = make(http.Header)
	// 强制设置拨号超时（例如 10 秒），防止网络不通时无限期等待
	config.Dialer = &net.Dialer{Timeout: 10 * time.Second}

	client := &Client{
		config:   config,
		AutoExit: autoExit,
	}

	client.Ctx, client.cancel = ctx.WithCancel(ctx.Background())
	return client, nil
}

// Connect 连接
func Connect(addr string, msgBufferSize int) (*Client, error) {
	client, err := NewClient(addr, true)
	if err != nil {
		return client, err
	}
	err = client.Connect()
	return client, err
}

// SetHeader 设置header
func (receiver *Client) SetHeader(key, value string) {
	receiver.config.Header.Set(key, value)
}

// SetReadDeadline 设置读取超时时间
func (receiver *Client) SetReadDeadline(timeout time.Duration) {
	receiver.conn.SetReadDeadline(time.Now().Add(timeout))
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
	if err != nil {
		receiver.cancel()
	} else {
		receiver.Ctx, receiver.cancel = ctx.WithCancel(ctx.Background())
	}
	receiver.isClose = err != nil || receiver.conn == nil
	return err
}

// Receiver 接收消息（并反序列化成val）
func (receiver *Client) Receiver(val any) error {
	if receiver.conn == nil {
		err := errors.New("连接未建立")
		receiver.errorIsClose(err)
		return err
	}

	err := websocket.JSON.Receive(receiver.conn, val)
	if err != nil {
		receiver.errorIsClose(err)
	}
	return err
}

// ReceiverMessage 接收消息
func (receiver *Client) ReceiverMessage() (string, error) {
	if receiver.conn == nil {
		err := errors.New("连接未建立")
		receiver.errorIsClose(err)
		return "", err
	}

	var msg string
	// 💡 websocket.Message 會自動幫你把一個完整的 WebSocket 幀讀進來
	err := websocket.Message.Receive(receiver.conn, &msg)
	if err != nil {
		receiver.errorIsClose(err)
		return "", err
	}
	return msg, nil
}

// Send 发送消息，如果msg不是go的基础类型，则会自动序列化成json
func (receiver *Client) Send(msg any) error {
	if receiver.conn == nil {
		err := errors.New("连接未建立")
		receiver.errorIsClose(err)
		return err
	}

	receiver.mu.Lock() // 💡 加鎖，確保同一時間只有一個協程在發包
	defer receiver.mu.Unlock()

	switch fastReflect.PointerOf(msg).Type {
	case fastReflect.GoBasicType:
		_, err := receiver.conn.Write([]byte(parse.ToString(msg)))
		if err != nil {
			receiver.errorIsClose(err)
		}
		return err
	default:
		marshalBytes, err := snc.Marshal(msg)
		if err != nil {
			return fmt.Errorf("发送数据时，出现反序列失败：%s", err.Error())
		}
		_, err = receiver.conn.Write(marshalBytes)
		if err != nil {
			receiver.errorIsClose(err)
		}
		return err
	}
}

// Close 关闭连接
func (receiver *Client) Close() {
	_ = receiver.conn.Close()
	receiver.cancel()
	receiver.isClose = true
}

// IsClose 是否已断开连接
func (receiver *Client) IsClose() bool {
	return receiver.isClose
}

// 根据错误信息，判断是否为断开连接导致的
func (receiver *Client) errorIsClose(err error) {
	var opError *net.OpError
	if receiver.conn == nil || errors.As(err, &opError) || err.Error() == "EOF" {
		receiver.cancel()
		receiver.isClose = true
		if receiver.AutoExit {
			exception.ThrowWebException(408, "服务端已关闭")
		}
	}
}
