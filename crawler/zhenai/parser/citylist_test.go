package zhenai

import (
	"io/ioutil"
	"testing"
)

func TestParseCityList(t *testing.T) {
	//测试数据
	contents, err := ioutil.ReadFile("citylist_test_data.html")
	if err != nil {
		panic(err)
	}

	results := ParseCityList(contents, "")

	//verify
	const resultSize = 470
	expectedUrls := []string{
		"", "", "",
	}
	//expectedCities := []string{
	//	"", "", "",
	//}

	if len(results.Requests) != resultSize {
		t.Errorf("result Should have %d"+"requests;but had %d ", resultSize, len(results.Requests))
	}
	for i, url := range expectedUrls {
		if results.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s;but "+"was %s", i, url, results.Requests[i].Url)
		}
	}
	//
	//if len(results.Items) != resultSize {
	//	t.Errorf("result Should have %d"+"requests;but had %d ", resultSize, len(results.Items))
	//}
	//for i, city := range expectedCities {
	//	if results.Items[i].(string) != city {
	//		t.Errorf("expected city #%d: %s;but "+"was %s", i, city, results.Items[i].(string))
	//	}
	//}

}
