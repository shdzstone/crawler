package main

import (
	"errors"
	"flag"
	"log"
	"net/rpc"
	config2 "stone/go/crawler/config"
	"stone/go/crawler/engine"
	"stone/go/crawler/scheduler"
	xcar "stone/go/crawler/xcar/parser"
	"stone/go/crawler_distributed/persist/client"
	"stone/go/crawler_distributed/rpcsupport"
	worker "stone/go/crawler_distributed/worker/client"
	"strings"
)

/*
一、分布式版爬虫要解决的问题
1.限流问题
* 单节点能够承受的流量有限
* 将worker放到不同的节点
2.去重问题
* 单节点能承受的去重数据量有限
* 无法保存之前去重结果
* 基于key-value store(Redis)进行分布式去重
具体实现：
* 将去重功能单独拉出来做成一个服务
* engine每次调用worker前调用去重服务，验证成功则调用
* 也可以使用worker协程调用去重服务，以免去重服务卡顿造成服务卡顿
3.存储问题
* 存储部分的结构，技术栈和爬虫部分区别很大
* 进一步优化需要特殊的elasticsearch技术背景
* 固有分布式，因为人的知识领域是有边界的，不同公司的人搭出的分布式的架构就不一样
4.本课程架构
* 解决限流问题（理论上）：需要放在不同节点的docker container中去解决
* 解决存储问题
* 解决分布式去重：从docker拉一个redis过来，再在项目中拉一个redis的go client,将之前存储于内存map的数据，存储于redis中
二、分布式爬虫的难点
1.从channel到分布式
* goroutine->goroutine
* goroutine client->goroutine(rpc client)<->goroutine(rpc server)-> goroutine server
*/

var (
	itemSaverHost = flag.String("itemsaver_host", "", "itemsaver host")
	workerHosts   = flag.String("worker_hosts", "", "worker hosts (comma separated)")
)

/*
命令行参数启动方式：go run main.go --itemsaver_host="1234" --worker_hosts="9000,9001,9002"
*/
func main() {
	flag.Parse()

	//persist rpc客户端，使用通道以接收item
	itemChan, err := client.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}
	//初始化worker rpc客户端
	pool, err := createClientPool(strings.Split(*workerHosts, ","))
	if err != nil {
		panic(err)
	}
	processor := worker.CreateProcessor(pool)

	//引擎配置中心
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      40,
		ItemChan:         itemChan,
		RequestProcessor: processor,
	}

	e.Run(engine.Request{
		Url:    "http://newcar.xcar.com.cn",
		Parser: engine.NewFuncParser(xcar.ParseCarList, config2.ParseCarList),
		//Url:    "http://www.zhenai.com/zhenghun",
		//Parser: engine.NewFuncParser(parser.ParseCityList, config2.ParseCityList),
	})

}

func createClientPool(hosts []string) (chan *rpc.Client, error) {
	var clients []*rpc.Client
	for _, h := range hosts {
		client, err := rpcsupport.NewClient(h)
		if err == nil {
			clients = append(clients, client)
			log.Printf("Connected to %s", h)
		} else {
			log.Printf("Error connecting to %s: %v", h, err)
		}
	}

	if len(clients) == 0 {
		return nil, errors.New("no connections available")
	}
	out := make(chan *rpc.Client)
	go func() {
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out, nil
}
