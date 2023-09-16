package apis

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"hall_web/global"
	"hall_web/service/domains"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

/**
需求：
分布式大厅聊天
前端先获取聊天列表，再开启长轮询
每个用户发送信息时，会触发长轮询后续获取消息操作，使其它用户收到消息

设计思路：
1 redis储存消息列表，并使用发布订阅，需要websocket
暂时最佳，但需要多引入中间件
2 长轮询，用户发送消息后，触发webhooks钩子，请求所有节点的addChat接口来解决分布式问题
但服务重启或崩溃就会丢失内存消息数据
3 分布式redis缓存，前端轮询，隔段时间就请求一次消息列表
但实时性差了些，而且影响前端性能

目前采用第二种方法
*/

type MessageQueue struct {
	mutex *sync.Mutex // 保护message
	*notify
	*queue
}

type queue struct {
	head    int
	tail    int
	length  int
	message []*domains.Message
}

type notify struct {
	mutex   *sync.RWMutex
	cond    *sync.Cond
	canNext bool
}

func NewNotify() *notify {
	m := new(sync.RWMutex)
	return &notify{mutex: m, cond: sync.NewCond(m)}
}

func (n *notify) broadcast() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.cond.Broadcast()
}

func (n *notify) check() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.canNext
}

func (n *notify) modifyCondition(val bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.canNext = val
}

func init() {
	MQ = NewMessageQueue(100)
}

func NewMessageQueue(l int) *MessageQueue {
	return &MessageQueue{mutex: new(sync.Mutex), notify: NewNotify(),
		queue: &queue{message: make([]*domains.Message, l, l), head: -1}}
}

func (mq *MessageQueue) Push(m *domains.Message) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	// 环形队列
	mq.head++
	if mq.head == len(mq.message) {
		mq.head = 0
	}
	mq.message[mq.head] = m

	if mq.head == mq.tail && mq.length == len(mq.message) {
		mq.tail++
		if mq.tail == len(mq.message) {
			mq.tail = 0
		}
	}
	if mq.length < len(mq.message) {
		mq.length++
	}
}

var MQ *MessageQueue

// 目前还是本地接口，需要更改为线上的
func SendChat(ctx *gin.Context) {
	chat := &domains.Chat{}
	err := ctx.Bind(chat)
	if err != nil {
		zap.S().Infof("[SendChat]: jsonErr %v", err)
		return
	}

	form := url.Values{}
	form.Set("id", strconv.Itoa(chat.Id))
	form.Set("user_id", strconv.Itoa(chat.Userid))
	form.Set("nickName", chat.Nickname)
	form.Set("image", chat.Image)
	form.Set("time", chat.Time)
	form.Set("content", chat.Content)

	for _, service := range global.ConsulHallWebServices {
		targetUrl := fmt.Sprintf("http://%s:%d/v1/chat/addChat", service.ServiceAddress, service.ServicePort)
		zap.S().Infof("[SendChat]: targetUrl %s", targetUrl)

		// 可考虑重试机制
		client := new(http.Client)
		resp, err := client.Post(targetUrl, "application/x-www-form-urlencoded", bytes.NewReader([]byte(form.Encode())))
		if err != nil {
			zap.S().Infof("[SendChat]: post err: %v %v", resp, err)
		}
	}
}

func AddChat(ctx *gin.Context) {
	chat := &domains.Chat{}
	err := ctx.Bind(chat)
	if err != nil {
		zap.S().Infof("[AddChat]: jsonErr %v", err)
		return
	}

	if chat.Time == "" {
		chat.Time = time.Now().Format("2006-01-02 15:04")
	}

	message := &domains.Message{
		Id:      int(time.Now().Unix()),
		Type:    0,
		Content: chat,
	}
	MQ.Push(message)
	//MQ.modifyCondition(true)
	//defer MQ.modifyCondition(false)
	MQ.cond.Broadcast()

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": "发送消息成功",
	})
}

func ChatList(ctx *gin.Context) {
	res := make([]*domains.Message, MQ.length)

	index := MQ.tail
	for i := 0; i < MQ.length; i++ {
		res[i] = MQ.message[index]
		if index == MQ.head {
			break
		}

		index++
		if index == len(MQ.message) {
			index = 0
		}
	}
	ctx.JSON(http.StatusOK, res)
}

func Listen(ctx *gin.Context) {
	MQ.cond.L.Lock()
	defer MQ.cond.L.Unlock()

	MQ.cond.Wait()
	ctx.JSON(http.StatusOK, MQ.message[MQ.head])
}
