package scheduler

import (
	"stone/go/crawler/engine"
)

/*
1.将简易的调度器做成队列，用以控制worker任务的数量
* 将从engine接收到的request放到一个request队列
* 将并发的worker协程中准备好接收数据的协程的通道放到一个worker chan队列中
* 通过workerChan通道统一的分发给各个worker chan，以实现简单的任务调度

2.通道使用难点
* 将从engine接收到的变量推入一个request队列
* 将可用（可发送数据）的worker chan放到一个worker chan队列
* 通过select-case代码块的阻塞特性及队列的使用，同时实现接收变量、接收可用的通道、将可用的变量通过"通道的通道"分发给可用的通道
*/

type QueuedScheduler struct {
	//scheduler:用以接收从engine主协程发送过来的request
	requestChan chan engine.Request
	//scheduler:用以向worker chan发送request，channel of channel：把worker的channel贯在一个channel中
	workerChan chan chan engine.Request
}

func (s *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

//engine发送request，则scheduler接收
func (s *QueuedScheduler) Submit(r engine.Request) {
	s.requestChan <- r
}

//worker goroutine状态准备好，则把该worker的接收chan与scheduler的发送chan对接
//相当于两个管子套在一起，即worker的管子套进scheduler的workerChan中
func (s *QueuedScheduler) WorkerReady(w chan engine.Request) {
	s.workerChan <- w
}

func (s *QueuedScheduler) Run() {
	s.requestChan = make(chan engine.Request)
	s.workerChan = make(chan chan engine.Request)
	//scheduler goroutine
	go func() {
		var requestQ []engine.Request     //request 队列
		var workerQ []chan engine.Request //worker chan队列
		//循环接收、发送共享的变量，以实现简单的调度
		for {
			//此处为常用的设计模式：将接收到的变量存到一个队列，将可发送的通道存到一个队列
			//通过select-case的阻塞特性，将接收到的变量，通过通道的通道，分发给可用的协程的通道
			//经典的简易调度中心设计模式
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			//判断是否同时有request和worker chan
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeRequest = requestQ[0]
				activeWorker = workerQ[0]
			}
			//select：没有default，阻塞直至有一个case通过
			select {
			//接收requestChan的值
			case r := <-s.requestChan:
				//send r to a ?worker
				requestQ = append(requestQ, r)
			//接收workerChan的值
			case w := <-s.workerChan:
				//send ?next_request to w
				workerQ = append(workerQ, w)
			//workerChan有值，且其值为chan engine.Requests，requestChan有值，其值为engine.Requests
			case activeWorker <- activeRequest:
				requestQ = requestQ[1:]
				workerQ = workerQ[1:]
			}
		}
	}()
}
