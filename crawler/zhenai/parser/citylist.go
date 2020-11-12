package zhenai

import (
	"regexp"
	config2 "stone/go/crawler/config"
	"stone/go/crawler/engine"
)

const cityListRe = `<a href="(http://www.parser.com/zhenghun/[a-z0-9]+)"[^>]*>([^<]+)</a>`

//const cityListRe = `<script>window.__INITIAL_STATE__=(.+);\(function`
func ParseCityList(contents []byte, _ string) engine.ParserResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1)

	result := engine.ParserResult{}
	for i, m := range matches {
		if i < 1 {
			result.Requests = append(result.Requests, engine.Request{
				Url:    string(m[1]),
				Parser: engine.NewFuncParser(ParseCity, config2.ParseCity),
			})
		}
	}
	return result
}
