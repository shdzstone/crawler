package view

import (
	"os"
	"stone/go/crawler/engine"
	"stone/go/crawler/frontend/model"
	common "stone/go/crawler/model"
	"testing"
)

func TestSearchResultView_Render(T *testing.T) {
	view := CreateSearchResultView("template.html")

	out, err := os.Create("template_test.html")

	page := model.SearchResult{}
	page.Hits = 123
	item := engine.Item{
		Url:  "https://album.parser.com/u/1076867343",
		Type: "parser",
		Id:   "1076867343",
		Payload: common.Profile{
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

	for i := 0; i < 10; i++ {
		page.Items = append(page.Items, item)
	}
	err = view.Render(out, page)
	//err := template.Execute(os.Stdout, page)
	if err != nil {
		panic(err)
	}
}
