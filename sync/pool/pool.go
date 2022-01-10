package pool

import (
	"sync"
)

/*

	基于协程实现的协程池 2020年10月27日20:08:29
 	Goroutine Pool

*/

// Pool 池接受来自客户端的任务，它限制总数
// 通过回收goroutines到一个给定的数字的goroutines。
type Pool struct {
	// 容量
	capacity int

	actionJob chan *Task

	transmitChan chan *Task

	// 同步操作的锁
	lock sync.Mutex
}

type Task struct {
	fc func() error
}

func NewPool(cap int) *Pool {
	return &Pool{
		capacity:     cap,
		actionJob:    make(chan *Task),
		transmitChan: make(chan *Task),
		lock:         sync.Mutex{},
	}
}

func NewTask(f func() error) *Task {
	return &Task{fc: f}
}

func (t *Task) execute() {
	_ = t.fc()
}

func (p *Pool) work() {
	for aj := range p.actionJob {
		go aj.execute()
	}
}

func (p *Pool) goRoutines(t *Task) {
	go func() {
		p.actionJob <- t
	}()
}

func (p *Pool) Run() {
	for i := 0; i <= p.capacity; i++ {
		go p.work()
	}
}
