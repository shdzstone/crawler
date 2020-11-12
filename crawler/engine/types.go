package engine

import config2 "stone/go/crawler/config"

type ParserFunc func(contents []byte, url string) ParserResult

//Parser接口:用以处理序列化与反序列化
type Parser interface {
	//调用函数本身
	Parse(contents []byte, url string) ParserResult
	//序列化：根据函数名称和参数，返回具体的parser函数
	//参数格式：{"ParseCityList",nil},{"ParseProfile",userName}
	Serialize() (funcName string, args interface{})
}

//请求封装
type Request struct {
	Url    string //解析出来的URL
	Parser Parser //将之前处理这个URL所需要的函数，转成一个接口
}

//解析结果封装
type ParserResult struct {
	Requests []Request
	Items    []Item
}

//单条解析信息
type Item struct {
	Url     string
	Type    string
	Id      string
	Payload interface{}
}

//空解析器
type NilParser struct{}

//空解析器解析
func (NilParser) Parse(_ []byte, _ string) ParserResult {
	return ParserResult{}
}

//空解析器序列化
func (NilParser) Serialize() (funcName string, args interface{}) {
	return config2.NilParser, nil
}

//本结构体负责函数与及序列化后的函数间的对应
//rpc客户端调用时，在网络上传输的是函数名，将具体的函数转换成函数名在网络上传输
//rpc服务端接到客户的调用后，通过具体的函数名，转换成具体的函数，并进行调用
//通过函数名实现函数的序列化及反序列化
type FuncParser struct {
	//具体的parser解析器
	parser ParserFunc
	//解析器的名字
	funcName string
}

//实际的解析函数
func (f *FuncParser) Parse(contents []byte, url string) ParserResult {
	return f.parser(contents, url)
}

//序列化
func (f *FuncParser) Serialize() (funcName string, args interface{}) {
	return f.funcName, nil
}

//工厂函数：创建实际的函数解析器
func NewFuncParser(p ParserFunc, funcName string) *FuncParser {
	return &FuncParser{
		parser:   p,
		funcName: funcName,
	}
}
