package main

import (
	"stone/go/crawler/engine"
	"stone/go/crawler/model"
	"stone/go/crawler_distributed/config"
	"stone/go/crawler_distributed/rpcsupport"
	"testing"
	"time"
)

/*
1.things to do
* start ItemSaverServer
* start ItemSaverClient
* call saver
*/
func TestItemSaver(t *testing.T) {
	go serveRpc(config.ItemSaverPort, "ping")
	time.Sleep(time.Second) //需要休眠以等待rpc server处理rpc调用

	client, err := rpcsupport.NewClient(config.ItemSaverPort)
	if err != nil {
		panic(err)
	}

	item := engine.Item{
		Url:  "https://album.parser.com/u/1076867343",
		Type: "parser",
		Id:   "1076867343",
		Payload: model.Profile{
			Age:        "34",
			Height:     "162",
			Weight:     "57",
			Income:     "3001-5000元",
			Gender:     "女",
			Name:       "安静的雪",
			Xingzuo:    "牡羊座",
			Occupation: "人事/行政",
			Marriage:   "离异",
			House:      "已购房",
			Hukou:      "山东菏泽",
			Education:  "大学本科",
			Car:        "未购车",
		},
	}

	result := ""
	err = client.Call("ItemSaverService.Save", item, &result)
	if err != nil || result != "ok" {
		t.Errorf("result:%s; err:%s", result, err)
	}
}
