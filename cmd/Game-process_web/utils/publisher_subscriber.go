package utils

import (
	"sync"
	"time"
)

type (
	subscriber chan interface{}         // 订阅者为一个管道
	topicFunc  func(v interface{}) bool // 主题为一个过滤器
)

// 发布者对象
type Publisher struct {
	m           sync.RWMutex
	buffer      int
	timeout     time.Duration
	subscribers map[subscriber]topicFunc
}

// 构建一个发布者对象，可以设置发布超时时间和缓存队列的长度
func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		timeout:     publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// 添加一个新的订阅者，订阅过滤筛选后的主题
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.m.Lock()
	p.subscribers[ch] = topic
	p.m.Unlock()
	return ch
}

// 订阅全部主题，没有过滤
func (p *Publisher) Subscriber() chan interface{} {
	return p.SubscribeTopic(nil)
}

// 退出订阅
func (p *Publisher) Evict(sub chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()
	delete(p.subscribers, sub)
	close(sub)
}

// 发送主题，可以容忍一定的超时
func (p *Publisher) sendTopic(sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	// 没有订阅该主题
	if topic != nil && !topic(v) {
		return
	}
	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}

// 发布一个主题
func (p *Publisher) Publish(v interface{}) {
	p.m.RLock()
	defer p.m.RUnlock()
	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// 关闭发布者对象，同时关闭所有的订阅者管道
func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()
	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}
