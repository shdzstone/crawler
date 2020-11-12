package persist

import (
	"context"
	"encoding/json"
	"stone/go/crawler/engine"
	"stone/go/crawler/model"
	"testing"

	"github.com/olivere/elastic/v7"
)

func TestSave(t *testing.T) {
	expected := engine.Item{
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

	//TODO: try to start up elastic search
	// here using docker go client
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	index := "dating_test"
	//save expected item
	err = Save(client, index, expected)
	if err != nil {
		panic(err)
	}

	//fetch saved item
	resp, err := client.Get().
		Index(index).
		Type(expected.Type).
		Id(expected.Id).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	t.Logf("%s", resp.Source)

	var actual engine.Item
	json.Unmarshal(resp.Source, &actual)
	actualProfile, _ := model.FromJsonObj(actual.Payload)
	actual.Payload = actualProfile
	//verify result
	if expected != actual {
		t.Logf("Got %v; Expected %v", actual, expected)
	}
}
