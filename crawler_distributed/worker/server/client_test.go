package main

import (
	"fmt"
	config2 "stone/go/crawler/config"
	"stone/go/crawler_distributed/config"
	"stone/go/crawler_distributed/rpcsupport"
	"stone/go/crawler_distributed/worker"
	"testing"
	"time"
)

func TestCrawService(t *testing.T) {
	const host = "9000"
	go rpcsupport.ServeRpc(host, worker.CrawlService{})
	time.Sleep(time.Second)

	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(err)
	}

	req := worker.Request{
		Url: "https://newcar.xcar.com.cn/m57484/",
		Parser: worker.SerializedParser{
			Name: config2.ParseCarModel,
			Args: "安静的雪",
		},
	}
	var result worker.ParseResult
	err = client.Call(config.CrawlServiceRpc, req, &result)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("xxxx", result)
	}

	//TODO: Verify results
}
