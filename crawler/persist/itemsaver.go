package persist

import (
	"context"
	"errors"
	"log"
	"stone/go/crawler/engine"

	"github.com/olivere/elastic/v7"
)

func ItemSaver(index string) (chan engine.Item, error) {
	out := make(chan engine.Item)

	client, err := elastic.NewClient(
		//SetURL可以省略，省略后默认找本机的9200端口
		//elastic.SetURL("http_publish_address"),
		//Must turn off sniff in docker
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	go func() {
		itemCount := 0
		for {
			item := <-out
			//log.Printf("Item Saver:got item #%d:%v", itemCount, item)
			itemCount++

			err := Save(client, index, item)
			if err != nil {
				log.Printf("Item#%d saver error:%v: %v", itemCount, err, item)
			} else {
				log.Printf("Item#%d save success:%v ", itemCount, item)
			}
		}
	}()
	return out, nil
}

//调用ES的rest接口两种方法
//1.http.Post()
//2.ES go客户端
//2.1官方：github.com/elastic/go-elasticsearch
//2.2个人:github.com/olivere/elastic
//3.安装：go get github.com/olivere/elastic/v7
func Save(client *elastic.Client, index string, item engine.Item) error {

	if item.Type == "" {
		return errors.New("must apply type")
	}
	indexService := client.Index().
		Index(index).
		Type(item.Type).
		BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}
	_, err := indexService.Do(context.Background())
	if err != nil {
		return err
	}
	//fmt.Printf("%+v", resp)
	return nil
}
