package persist

import (
	"context"
	"log"
	"stone/go/crawler/config"

	"github.com/olivere/elastic/v7"
)

func DeleteAllElasticIndex() {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	resp, err := client.DeleteIndex(config.CarElasticIndex).Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("Elastic delete result:%v", resp)
}
