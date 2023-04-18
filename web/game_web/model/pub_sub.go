package model

// 发布者订阅者模型
type Publisher struct {
	Subscribers []*Subscriber
}

type Subscriber struct {
	ch chan interface{}
}

func NewPublisher() *Publisher {
	return &Publisher{Subscribers: nil}
}

func NewSubscriber(buffer int) *Subscriber {
	return &Subscriber{
		ch: make(chan interface{}, buffer),
	}
}

func (p *Publisher) AddSubscriber(sub *Subscriber) {
	p.Subscribers = append(p.Subscribers, sub)
}

func (p *Publisher) SendTopicMsg(sub *Subscriber, v interface{}) {
	sub.ch <- v
}

func (p *Publisher) Publish(v interface{}) {
	for _, subscriber := range p.Subscribers {
		go p.SendTopicMsg(subscriber, v)
	}
}
