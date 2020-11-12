package scheduler

import (
	"stone/go/crawler/engine"
)

//scheduler
type SimpleScheduler struct {
	//当前worker channel中可用的chan
	workerChan chan engine.Request
}

func (s *SimpleScheduler) WorkerChan() chan engine.Request {
	return s.workerChan
}

func (s *SimpleScheduler) WorkerReady(requests chan engine.Request) {

}

func (s *SimpleScheduler) Run() {
	s.workerChan = make(chan engine.Request)
}

//将Request分发给worker chan
func (s *SimpleScheduler) Submit(request engine.Request) {
	//scheduler: receive request and send request down to worker chan
	go func() { s.workerChan <- request }()
	//注：为每个分发操作开一个goroutine，否则因为未初始化当前通道的接收协程， 会阻塞主协程阻塞，从而不能继续向下继续执行
	//非缓冲通道阻塞机制：向非缓冲通道发送数据时，发送方和接收者必须同时都在运行，否则会阻塞当前协程。
}
