package main

import (
	"flag"
	"fmt"
	"log"
	"stone/go/crawler/fetcher"
	"stone/go/crawler_distributed/rpcsupport"
	"stone/go/crawler_distributed/worker"
)

var port = flag.Int("port", 0, "the port for me to listen on")

/*
命令行参数启动方式：go run server.go -port=9000
*/
//worker server
func main() {
	flag.Parse()
	fetcher.SetVerboseLogging()

	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}

	err := rpcsupport.ServeRpc(fmt.Sprintf("%d", *port), worker.CrawlService{})
	if err != nil {
		log.Fatal(err)
	}
}
