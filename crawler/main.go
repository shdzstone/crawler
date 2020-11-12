package main

import (
	"stone/go/crawler/config"
	"stone/go/crawler/engine"
	"stone/go/crawler/persist"
	"stone/go/crawler/scheduler"
	xcar "stone/go/crawler/xcar/parser"
)

/*
1.GBK->utf-8转换
* go get -u golang.org/x/text

2.网站编码监测并转换
* go get -u golang.org/x/net/html

3.获取城市名称和链接
* 使用css选择器：通过定位html元素来获取数据，go标准库不支持css选择器，但有第三方库支持
* 使用xPath：通过定位html无素来获取数据
* regexp：正则表达式，通过字符串匹配来获取数据，更加通用

4.解析器Parser
* 输入：utf-8编码的文本
* 输出：Requests{URL,对应Parser}列表，Item列表

5.单任务版爬虫
* 获取并打印所有城市第一页用户的详细信息

6.获取网页内容
* 使用http.Get()获取内容
* 使用Encoding来转换编码：GBK->UTF-8
* 使用charset.DetermineEncoding来判断编码

7.爬虫总体算法
* 城市列表解析器+城市解析器+用户解析器
*/

func main() {

	//*
	//配置ES
	itemChan, err := persist.ItemSaver(config.CarElasticIndex)
	if err != nil {
		panic(err)
	}
	//配置并发引擎
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      10,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}

	e.Run(engine.Request{
		Url:    "http://newcar.xcar.com.cn",
		Parser: engine.NewFuncParser(xcar.ParseCarList, config.ParseCarList),
	})
	//*/

	//persist.DeleteAllElasticIndex()

	/*
		simpleEngine := engine.SimpleEngine{}
		simpleEngine.Run(engine.Requests{
			Url:    "http://newcar.xcar.com.cn",
			Parser: engine.NewFuncParser(parser.ParseCarList, config.ParseCarList),
			//Url:    "http://www.zhenai.com/zhenghun",
			//Parser: engine.NewFuncParser(parser.ParseCityList, config2.ParseCityList),
		})
		//*/
}
