package zhenai

import (
	"io/ioutil"
	"stone/go/crawler/engine"
	"stone/go/crawler/model"
	"testing"
)

func TestParseProfile(t *testing.T) {
	contents, err := ioutil.ReadFile("profile_test_data.html")
	if err != nil {
		panic(err)
	}
	results := ParseProfile(contents, "https://album.parser.com/u/1076867343", "安静的雪")

	if len(results.Items) != 1 {
		t.Errorf("result Should contain 1 result; but was %d ", len(results.Items))
	}

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

	actual := interface{}(results.Items[0]).(engine.Item)
	if actual != expected {
		t.Errorf("expected %v; but was %v",
			expected, actual)
	}
}
