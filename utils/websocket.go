package utils

import (
	"encoding/base64"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

// http升级websocket协议的配置
var wsUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	maxReadTimeout  = 5 * time.Minute
	maxWriteTimeOut = 5 * time.Minute
)

// WsMessage websocket消息
type WsMessage struct {
	MessageType int
	Data        []byte
}

// WsConnection 封装websocket连接
type WsConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan   chan *WsMessage // 读取队列
	outChan  chan *WsMessage // 发送队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	CloseChan chan byte // 关闭通知
}

// 读取协程
func (wsConn *WsConnection) wsReadLoop() {
	var (
		msgType int
		data    []byte
		msg     *WsMessage
		err     error
	)
	_ = wsConn.wsSocket.SetReadDeadline(time.Now().Add(maxReadTimeout))
	for {
		// 读一个message
		if msgType, data, err = wsConn.wsSocket.ReadMessage(); err != nil {
			goto ERROR
		}
		msg = &WsMessage{
			msgType,
			data,
		}
		// 放入请求队列
		select {
		case wsConn.inChan <- msg:
		case <-wsConn.CloseChan:
			goto ERROR
		}
	}
ERROR:
	wsConn.WsClose()
}

// 发送协程
func (wsConn *WsConnection) wsWriteLoop() {
	var (
		msg *WsMessage
		err error
	)
	_ = wsConn.wsSocket.SetWriteDeadline(time.Now().Add(maxWriteTimeOut))
	for {
		select {
		// 取一个应答
		case msg = <-wsConn.outChan:
			// 写给websocket
			if err = wsConn.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				goto ERROR
			}
		case <-wsConn.CloseChan:
			goto ERROR
		}
	}
ERROR:
	wsConn.WsClose()
}

func (wsConn *WsConnection) WritePing(body []byte) error {
	return wsConn.wsSocket.WriteMessage(websocket.PingMessage, body)
}

func (wsConn *WsConnection) WritePong(body []byte) error {
	return wsConn.wsSocket.WriteMessage(websocket.PongMessage, body)
}

// InitWebsocket 初始化websocket
/************** 并发安全 API **************/
func InitWebsocket(resp http.ResponseWriter, req *http.Request) (wsConn *WsConnection, err error) {
	var (
		wsSocket *websocket.Conn
	)
	// 应答客户端告知升级连接为websocket
	if wsSocket, err = wsUpgrader.Upgrade(resp, req, nil); err != nil {
		return
	}
	wsConn = &WsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *WsMessage, 1024),
		outChan:   make(chan *WsMessage, 1024),
		CloseChan: make(chan byte),
		isClosed:  false,
	}
	//设置 websocket 协议层面对应的ping和pong 处理方法
	wsSocket.SetPingHandler(func(appData string) error {
		return wsConn.WritePing([]byte(appData))
	})
	wsSocket.SetPongHandler(func(appData string) error {
		return wsConn.WritePong([]byte(appData))
	})
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()

	return
}

// WsWrite 发送消息
func (wsConn *WsConnection) WsWrite(messageType int, data []byte) (err error) {
	select {
	case wsConn.outChan <- &WsMessage{messageType, data}:
	case <-wsConn.CloseChan:
		err = errors.New("WebSocket connection closed")
	}
	return
}

// WsRead 读取消息
func (wsConn *WsConnection) WsRead() (msg *WsMessage, err error) {
	select {
	case msg = <-wsConn.inChan:
		return
	case <-wsConn.CloseChan:
		err = errors.New("WebSocket connection closed")
	}
	return
}

// WsClose 关闭连接
func (wsConn *WsConnection) WsClose() {
	_ = wsConn.wsSocket.Close()
	wsConn.mutex.Lock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		//<-wsConn.closeChan
		close(wsConn.CloseChan)
	}
	wsConn.mutex.Unlock()
	return
}

// WsHandleError 错误检查
func WsHandleError(ws *WsConnection, err error) bool {
	if err != nil {
		dt := time.Now().Add(time.Second)
		if err := ws.wsSocket.WriteMessage(websocket.TextMessage, []byte(base64.StdEncoding.EncodeToString([]byte(err.Error())))); err != nil {
			logrus.Error(ws.wsSocket.RemoteAddr(), err)
		}
		if err := ws.wsSocket.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); err != nil {
			logrus.Error(ws.wsSocket.RemoteAddr(), err)
		}
		return true
	}
	return false
}

// XtermMessage web终端发来的包
type XtermMessage struct {
	Type  string `json:"type"`  // 类型:resize客户端调整终端, input客户端输入
	Input string `json:"input"` // msgtype=input情况下使用
	Rows  uint16 `json:"rows"`  // msgtype=resize情况下使用
	Cols  uint16 `json:"cols"`  // msgtype=resize情况下使用
}

const (
	WsMsgInput  = "input"
	WsMsgResize = "resize"
	Heartbeat   = "heartbeat"
)
