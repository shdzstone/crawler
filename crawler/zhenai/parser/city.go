package zhenai

import (
	"regexp"
	config2 "stone/go/crawler/config"
	"stone/go/crawler/engine"
)

var (
	profileRe = regexp.MustCompile(`<a href="(http://album.parser.com/u/[0-9]+)" [^>]*>([^<]+)</a>`)
	cityUrlRe = regexp.MustCompile(`<a href="(http://www.parser.com/zhenghun/[^"]+)"`)
)

func ParseCity(contents []byte, _ string) engine.ParserResult {
	//用户
	matches := profileRe.FindAllSubmatch(contents, -1)
	result := engine.ParserResult{}
	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]),
			Parser: NewProfileParser(string(m[2])),
		})
	}
	//猜你喜欢
	matches = cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]),
			Parser: engine.NewFuncParser(ParseCity, config2.ParseCity),
		})
	}
	return result
}
