package global

import (
	"errors"
	"fmt"
	"process_web/model"
	"process_web/model/response"
	"sync"

	"github.com/gorilla/websocket"
)

type WSConn struct {
	conn      *websocket.Conn
	userID    uint32                        //标识是哪个用户的连接
	inChan    chan model.Message            //读客户端发来数据
	outChan   chan response.MessageResponse //向客户端写入数据
	closeChan chan struct{}                 //标记ws是否关闭，关闭chan后，消费者依然可以读,通过这个标志说明是否关闭
	isClose   bool                          // 通道closeChan是否已经关闭,chan不能多次关闭，所有需要保证只能关闭一次
	once      sync.Once                     //保证closeChan只会关闭一次
}

// 创建websocket实例
func InitWebSocket(conn *websocket.Conn, userID uint32) (ws *WSConn) {
	ws = &WSConn{
		userID:    userID,
		inChan:    make(chan model.Message, 5),
		outChan:   make(chan response.MessageResponse, 30),
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
		//zap.S().Infof("[readMsgLoop]读到:%v", data)
		data := model.Message{}
		err := ws.conn.ReadJSON(&data)
		//zap.S().Infof("[readMsgLoop]:Data:%+v", data)
		if err != nil {
			//zap.S().Warnf("[readMsgLoop]:%s", err)
			//发生错误,关闭连接，停止协程
			ws.CloseConn()
			break
		}
		if data.Type == model.UserIntoMsg {
			continue
		}
		ws.inChan <- data
	}
}

// 协程发送客户端msg
func (ws *WSConn) writeMsgLoop() {
	for true {
		data := <-ws.outChan
		err := ws.conn.WriteJSON(data)
		if err != nil {
			//发生错误,关闭连接，停止协程
			ws.CloseConn()
			break
		}
	}
}

// 从inChan读
func (ws *WSConn) InChanRead() chan model.Message {
	return ws.inChan
}

// 写出outChan
func (ws *WSConn) OutChanWrite(data response.MessageResponse) error {
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

func SendErrToUser(ws *WSConn, handlerFunc string, error error) {
	if error != nil {
		errRsp := response.MessageResponse{
			MsgType: response.ErrResponseMsgType,
			ErrInfo: &response.ErrResponse{Error: fmt.Sprintf("[%s]:%s", handlerFunc, error)},
		}
		if ws != nil {
			_ = ws.OutChanWrite(errRsp)
		}
	}
}

func SendMsgToUser(ws *WSConn, data response.MessageResponse) {
	if ws != nil {
		_ = ws.OutChanWrite(data)
	}
}
