package utils

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	Handler func(...interface{})
	Params  []interface{}
}

// 工作池
type Pool struct {
	capacity       uint64
	runningWorkers uint64
	status         int64
	chTask         chan *Task //生产者消费者消息队列
	sync.Mutex
}

func (p *Pool) incRunning() { // runningWorkers + 1
	atomic.AddUint64(&p.runningWorkers, 1)
}

func (p *Pool) decRunning() { // runningWorkers - 1
	atomic.AddUint64(&p.runningWorkers, ^uint64(1))
}

func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.runningWorkers)
}

func (p *Pool) GetCap() uint64 {
	return p.capacity
}

func (p *Pool) setStatus(status int64) bool {
	p.Lock()
	defer p.Unlock()

	if p.status == status {
		return false
	}
	p.status = status

	return true
}

const (
	RUNNING = 1
	STOPED  = 0
)

func NewPool(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, errors.New("不合法容量")
	}
	return &Pool{
		capacity: capacity,
		status:   RUNNING,
		// 初始化任务队列, 队列长度为容量
		chTask: make(chan *Task, capacity),
	}, nil
}

// 池中，当有新任务Task，尝试开启协程处理消息，处理完成后关闭协程，即worker数量是动态变化的，但是无法超出最大值
func (p *Pool) run() {
	p.incRunning()
	go func() {
		defer func() {
			p.decRunning()
			if r := recover(); r != nil {
				//恢复
			}
			p.checkWorker()
		}()
		select {
		case task, ok := <-p.chTask:
			if !ok {
				return
			}
			task.Handler(task.Params...)
		}
	}()
}

// 供外部调用，当池中工作者太少，会自动扩充
func (p *Pool) Put(task *Task) error {
	p.Lock()
	defer p.Unlock()
	if p.status == STOPED {
		return errors.New("池已经关闭")
	}
	if p.GetRunningWorkers() < p.GetCap() {
		p.run()
	}
	//告知协程去消费
	p.chTask <- task
	return nil
}

func (p *Pool) ClosePool() {
	p.setStatus(STOPED)
	for len(p.chTask) > 0 {
		time.Sleep(1e6)
	}
	close(p.chTask) // 关闭任务队列
}

func (p *Pool) checkWorker() {
	p.Lock()
	defer p.Unlock()

	// 当前没有 worker 且有任务存在，运行一个 worker 消费任务
	// 没有任务无需考虑 (当前 Put 不会阻塞，下次 Put 会启动 worker)
	if p.runningWorkers == 0 && len(p.chTask) > 0 {
		p.run()
	}
}

func main() {
	// 创建任务池
	pool, err := NewPool(10)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 20; i++ {
		// 任务放入池中
		pool.Put(&Task{
			Handler: func(v ...interface{}) {
				fmt.Println(v)
			},
			Params: []interface{}{i, 666},
		})
		time.Sleep(1e9) // 等待执行
	}
}
