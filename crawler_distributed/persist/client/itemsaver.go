package client

import (
	"log"
	"stone/go/crawler/engine"
	"stone/go/crawler_distributed/config"
	"stone/go/crawler_distributed/rpcsupport"
)

//rpc客户端
func ItemSaver(host string) (chan engine.Item, error) {
	//创建RPC客户端
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		return nil, err
	}
	//创建管道
	out := make(chan engine.Item)
	//协程循环处理RPC客户端请求任务
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver:got item #%d:%v", itemCount, item)
			itemCount++

			//call rpc client to save item
			result := ""
			err = client.Call(config.ItemServerRpc, item, &result)
			if err != nil {
				log.Printf("Item Saver: error saving item %v: %v", item, err)
			}
		}
	}()
	return out, nil
}
