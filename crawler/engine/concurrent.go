package engine

/*
一.并发版爬虫架构
1.Worker
* 负责处理输入的Request，并将处理的结果Request+Items返回给Engine
* 将fetcher与parser合成一个模块:Worker

2.scheduler
* 负责将Request分发给worker

3.Engine
* 负责将worker处理的结果Request交给scheduler
* 负责将worker处理的结果Item输出

4.scheduler实现1
* 所有worker共用一个输入:会遇到channel阻塞引起的循环等待问题

5.scheduler实现2
* 为每个Request创建一个goroutine，每个goroutine分发相应的Request给worker

6.scheduler实现3
* request队列和worker chan队列
* 有request进来，则加入requestQ，缓存起来
* 有worker进来，则加入workerChanQ，缓存起来
* 若request队列和worker chan队列均有数据，则把request通过worker chan分发给worker

7.channel of channel
* channel也是一等公民，可以做函数参数、返回值、变量、类型等等均可
* goroutine对外的信息通道或接口就是chan，channel of chan相当于接口的chan

8.URL去重
* URL哈希表：有些URL过大，占据空间过大
* 计算URL的MD5等哈希，再存哈希表
* 使用bloom filter多重哈希结构
* 使用Redis等Key-Value存储系统实现分布式去重

9.如何存储Items
1> 抽象出Task的概念
* FetchTask、PersistTask共用一个Engine,scheduler
* 需要创建FetchTask、PersistTask
* 对本项目来讲显得过重
2> 为每个Item创建goroutine,提交给ItemSaver()
* go save()
* go func(){itemChan<-item}()

10.ItemSaver架构
* 对于engine来讲，有结果后马上分发给其它任务去处理，自身不适宜处理太重的任务
* ItemSaver的速度比Fetcher快
* 类似SimpleScheduler的方法即可

*/

type ConcurrentEngine struct {
	Scheduler        Scheduler //Scheduler
	WorkerCount      int       //worker的数量
	ItemChan         chan Item
	RequestProcessor Processor
}

type Processor func(request Request) (ParserResult, error)

//scheduler接口
type Scheduler interface {
	ReadyNotifier
	Submit(request Request)
	WorkerChan() chan Request
	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParserResult) //Worker->engine
	e.Scheduler.Run()

	for _, r := range seeds {
		//URL去重(dedup)
		if isDuplicate(r.Url) {
			continue
		}
		//engine通过函数调用的方式将request分发给scheduler
		e.Scheduler.Submit(r)
	}

	//worker初始化
	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	//engine主协程接收out通道中的处理结果
	for {
		result := <-out //engine: receive Worker chan parser result
		for _, item := range result.Items {
			if _, ok := interface{}(item).(Item); ok {
				//item存储
				go func() { e.ItemChan <- item }()
			}
		}

		for _, request := range result.Requests {
			//URL去重(dedup)
			if isDuplicate(request.Url) {
				continue
			}
			//engine通过函数调用的方式将request分发给scheduler
			e.Scheduler.Submit(request)
		}
	}
}

//注意三角channel发送与接收时循环等待的问题
func (e *ConcurrentEngine) createWorker(in chan Request, out chan ParserResult, ready ReadyNotifier) {
	go func() {
		for {
			//tell scheduler i'm ready
			ready.WorkerReady(in)
			request := <-in //Worker: receive scheduler request
			//result, err := Worker(request)
			result, err := e.RequestProcessor(request)

			if err != nil {
				continue
			}
			out <- result //Worker:send parser result
		}
	}()
}

//hash map将urls存储于内存中，缺陷是机器重启后数据会丢失
//但可以将该hash map于程序退出时或隔段时间存储于Redis中
var visitedUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}

	visitedUrls[url] = true
	return false
}
