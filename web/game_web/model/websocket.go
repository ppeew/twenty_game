package model

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type WSConn struct {
	conn      *websocket.Conn
	userID    uint32        //标识是哪个用户的连接
	inChan    chan []byte   //读客户端发来数据
	outChan   chan []byte   //向客户端写入数据
	closeChan chan struct{} //标记ws是否关闭，关闭chan后，消费者依然可以读,通过这个标志说明是否关闭
	isClose   bool          // 通道closeChan是否已经关闭,chan不能多次关闭，所有需要保证只能关闭一次
	once      sync.Once     //保证closeChan只会关闭一次
}

// 创建websocket实例
func InitWebSocket(conn *websocket.Conn, userID uint32) (ws *WSConn) {
	ws = &WSConn{
		userID:    userID,
		inChan:    make(chan []byte, 5),
		outChan:   make(chan []byte, 30),
		closeChan: make(chan struct{}),
		conn:      conn,
	}
	// 必要协程：读取客户端数据协程/发送数据协程
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
			//zap.S().Info("[writeMsgLoop]:用户websocket断开")
			ws.CloseConn()
			break
		}
	}
}

// 从inChan读
func (ws *WSConn) InChanRead() chan []byte {
	return ws.inChan
}

// 写出outChan
func (ws *WSConn) OutChanWrite(data []byte) error {
	if ws.isClose {
		return errors.New("连接已断开")
	}
	ws.outChan <- data
	return nil
}

// 关闭websocket
func (ws *WSConn) CloseConn() {
	ws.isClose = true
	ws.once.Do(func() {
		close(ws.closeChan)
	})
	_ = ws.conn.Close()
}

// 查询ws状态
func (ws *WSConn) IsDisConn() chan struct{} {
	return ws.closeChan
}
