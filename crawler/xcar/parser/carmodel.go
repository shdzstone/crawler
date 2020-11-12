package xcar

import (
	"regexp"
	config2 "stone/go/crawler/config"
	"stone/go/crawler/engine"
)

var carDetailRe = regexp.MustCompile(`<a href="(/m\d+/)" target="_blank"`)

func ParseCarModel(contents []byte, _ string) engine.ParserResult {
	matches := carDetailRe.FindAllSubmatch(contents, -1)

	result := engine.ParserResult{}
	for _, m := range matches {
		result.Requests = append(
			result.Requests,
			engine.Request{
				Url: host + string(m[1]),
				Parser: engine.NewFuncParser(
					ParseCarDetail, config2.ParseCarDetail),
			})
	}

	return result
}
