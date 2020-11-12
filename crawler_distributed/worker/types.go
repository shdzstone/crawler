package worker

import (
	"errors"
	"fmt"
	"log"
	config2 "stone/go/crawler/config"
	"stone/go/crawler/engine"
	xcar "stone/go/crawler/xcar/parser"
	zhenai "stone/go/crawler/zhenai/parser"
)

/*
1.go之序列化
在 Go 中并不是所有的类型都能进行序列化：
* JSON object key 只支持 string
* Channel、complex、function 等 type 无法进行序列化
* 数据中如果存在循环引用，则不能进行序列化，因为序列化时会进行递归
* Pointer 序列化之后是其指向的值或者是 nil
 还需要注意的是：只有 struct 中支持导出的 field 才能被 JSON package 序列化，即首字母大写的 field。

2.解析器的序列化与反序列化
* 解析器原本的定义为函数
* 需要处理函数的序列化与反序列化

3.函数名与参数转换为具体的函数
* 复杂点：将每个函数的名字与具体的函数变量注册到一个全局的map中，再根据具体的函数名找到具体的函数
* 简单点：简单的使用switch-case来将函数名与函数值进行匹配
*/

//序列化后的函数
type SerializedParser struct {
	Name string //function name函数名
	Args interface{}
}

//辅助结构体，用以在rpc间传输engine.Requests
type Request struct {
	Url    string
	Parser SerializedParser
}

//辅助结构体：用以在rpc间传输engine.parserResult
type ParseResult struct {
	Items    []engine.Item
	Requests []Request
}

//Request序列化辅助函数：将具体的Request转化为可以rpc间传输的Request结构体
func SerializeRequest(r engine.Request) Request {
	name, args := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: name,
			Args: args,
		},
	}
}

//ParseResult序列化辅助函数：
func SerializeResult(r engine.ParserResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		result.Requests = append(result.Requests, SerializeRequest(req))
	}
	return result
}

//request反序列化辅助函数：将rpc传输过来的Request转换为engine.Request
func DeserializeRequest(r Request) (engine.Request, error) {
	parser, err := deserializeParser(r.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	}, nil
}

//result反序列化辅助函数：将rpc传输过来的result转换为engine.ParseResult
func DeserializeResult(r ParseResult) engine.ParserResult {
	result := engine.ParserResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		engineReq, err := DeserializeRequest(req)
		if err != nil {
			log.Printf("error deserializing request:%v", err)
			continue
		}
		result.Requests = append(result.Requests, engineReq)
	}
	return result
}

//函数反序列化辅助函数：将rpc间传输的函数参数转换为具体的函数
//函数式编程：函数类型是类型、函数名为变量
func deserializeParser(p SerializedParser) (engine.Parser, error) {
	switch p.Name {
	case config2.NilParser:
		return engine.NilParser{}, nil
	case config2.ParseCityList:
		return engine.NewFuncParser(zhenai.ParseCityList, config2.ParseCityList), nil
	case config2.ParseCity:
		return engine.NewFuncParser(zhenai.ParseCity, config2.ParseCity), nil
	case config2.ParseProfile:
		if userName, ok := p.Args.(string); ok {
			return zhenai.NewProfileParser(userName), nil
		} else {
			return nil, fmt.Errorf("invalid arg:%v", p.Args)
		}

	case config2.ParseCarDetail:
		return engine.NewFuncParser(xcar.ParseCarDetail, config2.ParseCarDetail), nil
	case config2.ParseCarList:
		return engine.NewFuncParser(xcar.ParseCarList, config2.ParseCarList), nil
	case config2.ParseCarModel:
		return engine.NewFuncParser(xcar.ParseCarModel, config2.ParseCarModel), nil
	default:
		return nil, errors.New("unknown parser name")
	}
}
