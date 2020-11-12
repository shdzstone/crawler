package client

import (
	"net/rpc"
	"stone/go/crawler/engine"
	"stone/go/crawler_distributed/config"
	"stone/go/crawler_distributed/worker"
)

//worker客户端
func CreateProcessor(clientChan chan *rpc.Client) engine.Processor {
	//函数式编程：闭包，共享client
	return func(req engine.Request) (engine.ParserResult, error) {
		sReq := worker.SerializeRequest(req)
		var sResult worker.ParseResult
		c := <-clientChan
		err := c.Call(config.CrawlServiceRpc, sReq, &sResult)
		if err != nil {
			return engine.ParserResult{}, err
		}
		return worker.DeserializeResult(sResult), nil
	}
}
