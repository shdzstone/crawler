package xcar

import (
	"regexp"
	config2 "stone/go/crawler/config"
	"stone/go/crawler/engine"
)

const host = "http://newcar.xcar.com.cn"

var carModelRe = regexp.MustCompile(`<a href="(/\d+/)" target="_blank" class="list_img">`)
var carListRe = regexp.MustCompile(`<a href="(//newcar.xcar.com.cn/car/[\d+-]+\d+/)"`)

func ParseCarList(contents []byte, _ string) engine.ParserResult {
	matches := carModelRe.FindAllSubmatch(contents, -1)

	result := engine.ParserResult{}
	for _, m := range matches {
		result.Requests = append(
			result.Requests,
			engine.Request{
				Url:    host + string(m[1]),
				Parser: engine.NewFuncParser(ParseCarModel, config2.ParseCarModel),
			})
	}

	matches = carListRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(
			result.Requests, engine.Request{
				Url: "http:" + string(m[1]),
				Parser: engine.NewFuncParser(
					ParseCarList, config2.ParseCarList),
			})
	}

	return result
}
