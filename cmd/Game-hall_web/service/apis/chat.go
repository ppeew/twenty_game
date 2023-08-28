package apis

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"hall_web/service/domains"
	"net/http"
	"sync"
	"time"
)

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
