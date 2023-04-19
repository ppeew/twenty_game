package model

// 类似发布者订阅者模型
type Publisher struct {
	MsgChan     chan string            //接受信息管道
	Subscribers map[uint32]*Subscriber //用户id到订阅者的映射
}

type Subscriber struct {
	ch chan string

	//订阅者的websocket连接信息
	WS *WSConn
}

func NewPublisher() *Publisher {
	return &Publisher{Subscribers: make(map[uint32]*Subscriber)}
}

func NewSubscriber(buffer int, ws *WSConn) *Subscriber {
	return &Subscriber{
		ch: make(chan string, buffer),
		WS: ws,
		//UserID:  userID,
		//IsReady: isReady,
	}
}

func (p *Publisher) AddSubscriber(userID uint32, sub *Subscriber) {
	if p.Subscribers[userID] == nil {
		p.Subscribers[userID] = new(Subscriber)
	}
	p.Subscribers[userID] = sub
}

func (p *Publisher) SendTopicMsg(sub *Subscriber, v string) {
	sub.ch <- v
}

func (p *Publisher) Publish(s string) {
	for _, subscriber := range p.Subscribers {
		go p.SendTopicMsg(subscriber, s)
	}
}
