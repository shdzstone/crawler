package zhenai

import (
	"log"
	"regexp"
	"stone/go/crawler/config"
	"stone/go/crawler/engine"
	"stone/go/crawler/model"

	"github.com/bitly/go-simplejson"
)

var re = regexp.MustCompile(`<script>window.__INITIAL_STATE__=(.+);\(function`)
var idUrlRe = regexp.MustCompile(`http://album.parser.com/u/([0-9]+)`)

func ParseProfile(contents []byte, url string, name string) engine.ParserResult {
	var profile model.Profile
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		info := match[1]
		profile = parseJson(info, name)
	}
	id := extractString([]byte(url), idUrlRe)

	result := engine.ParserResult{
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "parser",
				Id:      id,
				Payload: profile,
			},
		},
	}
	return result
}

//解析json数据
func parseJson(info []byte, name string) model.Profile {
	jsonInfo, err := simplejson.NewJson(info)
	if err != nil {
		log.Println("json解析失败...")
	}
	basicInfo, _ := jsonInfo.Get("objectInfo").Get("basicInfo").Array()
	//basicInfo是一个切片，类型是interface{}
	var profile model.Profile

	profile.Name = name
	profile.Gender, _ = jsonInfo.Get("objectInfo").Get("genderString").String()
	//id, _ := jsonInfo.Get("objectInfo").Get("memberID").String()

	//遍历basicInfo切片，使用断言来判断类型
	for k, v := range basicInfo {
		if e, ok := v.(string); ok {
			switch k {
			case 0:
				profile.Marriage = e
			case 1:
				//年龄:47岁，我们可以设置int类型，所以可以通过另一个json字段来获取
				profile.Age = e
			case 2:
				profile.Xingzuo = e
			case 3:
				profile.Height = e
			case 4:
				profile.Weight = e
			case 6:
				profile.Income = e
			case 7:
				profile.Occupation = e
			case 8:
				profile.Education = e
			}
		}
	}

	return profile
}

func extractString(contents []byte, re *regexp.Regexp) string {
	log.Printf("url:%s", contents)
	match := re.FindSubmatch(contents)

	for _, item := range match {
		log.Printf("match item:%s", item)
	}

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

//函数序列化结构体：目的是将userName打包
type ProfileParser struct {
	//用户名
	userName string
}

//解析函数
func (p *ProfileParser) Parse(contents []byte, url string) engine.ParserResult {
	return ParseProfile(contents, url, p.userName)
}

//序列化
func (p *ProfileParser) Serialize() (name string, args interface{}) {
	return config.ParseProfile, p.userName
}

//函数式编程：输出参数均为函数，把输入的参数包装下，包装成一个输出的函数来输出
//go语言中自定义的函数都是闭包
func NewProfileParser(name string) *ProfileParser {
	return &ProfileParser{
		userName: name,
	}
}
