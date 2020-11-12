package main

import (
	"flag"
	"fmt"
	"log"
	config2 "stone/go/crawler/config"
	"stone/go/crawler_distributed/persist"
	"stone/go/crawler_distributed/rpcsupport"

	"github.com/olivere/elastic/v7"
)

var port = flag.Int("port", 0, "the port for me to listen on")

/*
命令行参数启动方式：go run itemsaver.go -port=1234
*/
func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}

	//出错了强制退出，panic()的另一种写法
	log.Fatal(serveRpc(fmt.Sprintf("%d", *port), config2.CarElasticIndex))
}

//rpc服务器
//参数配置：主机名及elastic index名
func serveRpc(host string, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}
	//注：此处ItemSaverService的Save方法接收参数是指针类型，所以这里要传指针
	return rpcsupport.ServeRpc(host, &persist.ItemSaverService{
		Client: client,
		Index:  index,
	})
}
