package worker

import (
	"stone/go/crawler/engine"
)

//worker rpc服务定义
type CrawlService struct{}

func (w CrawlService) Process(req Request, result *ParseResult) error {
	engineReq, err := DeserializeRequest(req)
	if err != nil {
		return err
	}
	engineResult, err := engine.Worker(engineReq)
	if err != nil {
		return err
	}
	*result = SerializeResult(engineResult)
	return nil
}
