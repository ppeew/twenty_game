package model

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type WSConn struct {
	UserID    uint32      //标识是哪个用户的连接
	inChan    chan []byte //读客户端发来数据
	outChan   chan []byte //向客户端写入数据
	closeChan chan []byte
	isClose   bool       // 通道closeChan是否已经关闭
	mutex     sync.Mutex //并发问题
	conn      *websocket.Conn
}

// 创建websocket实例
func InitWebSocket(conn *websocket.Conn, userID uint32) (ws *WSConn) {
	ws = &WSConn{
		UserID:    userID,
		inChan:    make(chan []byte, 1024),
		outChan:   make(chan []byte, 1024),
		closeChan: make(chan []byte, 1024),
		conn:      conn,
	}
	// 完善必要协程：读取客户端数据协程/发送数据协程
	go ws.readMsgLoop()
	go ws.writeMsgLoop()
	return
}

// 协程接受客户端msg
func (ws *WSConn) readMsgLoop() {
	for true {
		_, data, err := ws.conn.ReadMessage()
		if err != nil {
			//发生错误,关闭连接，停止协程
			ws.CloseConn()
			break
		}
		ws.inChan <- data
	}
}

// 协程发送客户端msg
func (ws *WSConn) writeMsgLoop() {
	for true {
		data := <-ws.outChan
		err := ws.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			//发生错误,关闭连接，停止协程
			ws.CloseConn()
			break
		}
	}
}

// 从inChan读
func (ws *WSConn) InChanRead() (data []byte, err error) {
	select {
	case data = <-ws.inChan:
	case <-ws.closeChan:
		err = errors.New("连接断开")
	}
	return
}

// 写出outChan
func (ws *WSConn) OutChanWrite(data []byte) (err error) {
	select {
	case ws.outChan <- data:
	case <-ws.closeChan:
		err = errors.New("连接断开")
	}
	return
}

// 关闭websocket
func (ws *WSConn) CloseConn() {
	ws.mutex.Lock()
	if !ws.isClose {
		close(ws.closeChan)
		ws.isClose = true
	}
	ws.mutex.Unlock()
	ws.conn.Close()
}

func (ws *WSConn) IsClose() bool {
	if ws.isClose {
		return true
	}
	return false
}
