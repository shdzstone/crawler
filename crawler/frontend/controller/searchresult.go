package controller

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"stone/go/crawler/config"
	"stone/go/crawler/engine"
	"stone/go/crawler/frontend/model"
	"stone/go/crawler/frontend/view"
	"strconv"
	"strings"

	"github.com/olivere/elastic/v7"
)

/*
1.html/template简介
* 模板引擎
* 服务端器最终生成网页
* 简单的静态页面，非ajax动态页面
* 适合做后台或者维护页面

2.html/template用法：
* 取值
* 选择
* 循环
* 函数调用

3.前端展示：
* 搜索字符串重写
* 翻页
* 使用http.FileServer来提供静态内容，css、js、图片、首页等

4.http标准库用法
* 使用http.Handle()或http.HandleFunc()添加路由对应的handle对象或处理函数
* 实现HandleFunc或handle对象
* 使用http.ListenAndServe()启动并监听相应端口

5.go标准库html/template工作流程
* 使用http库启用前端服务器server
* 使用Http.Handle()注册前端服务器server的路由处理器，即不同的controller
* controller首先获取该路由对应的页面的模板，生成相应的view
* controller再从elasticsearch或数据库获取数据，生成相应的model
* controller使用view.render将相应的model写入Handle()或HandleFunc()中的http.ResponseWriter并返回给客户端
* 客户端根据前端服务器返回的response html进行渲染页面
*/

//查询结果
type SearchResultHandler struct {
	view   view.SearchResultView
	client *elastic.Client
}

func CreateSearchResultHandler(template string) SearchResultHandler {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	return SearchResultHandler{
		view:   view.CreateSearchResultView(template),
		client: client,
	}
}

//TODO: 1
//rewrite query string
//support paging
//add start page

//localhost:8888/search?q=男 已购房&from=20
func (h SearchResultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("Url:%s", req.URL)
	q := strings.TrimSpace(req.FormValue("q"))
	from, err := strconv.Atoi(req.FormValue("from"))
	if err != nil {
		from = 0
	}

	//fmt.Fprintf(w, "q=%s,from=%d", q, from)

	page, err := h.getSearchResult(q, from)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = h.view.Render(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

const pageSize = 10

func (h SearchResultHandler) getSearchResult(q string, from int) (model.SearchResult, error) {
	var result model.SearchResult
	resp, err := h.client.
		Search(config.CarElasticIndex).
		Query(elastic.NewQueryStringQuery(rewriteQueryString(q))).
		From(from).
		Do(context.Background())
	if err != nil {
		return result, err
	}
	result.Hits = resp.TotalHits()
	result.Start = from
	result.Query = q

	result.Items = resp.Each(reflect.TypeOf(engine.Item{}))
	if result.Start == 0 {
		result.PreFrom = -1
	} else {
		result.PreFrom = (result.Start - 1) / pageSize * pageSize
	}
	result.NextFrom = result.Start + len(result.Items)

	result.PreFrom = result.Start - len(result.Items)
	result.NextFrom = result.Start + len(result.Items)

	log.Printf("q=%s,from=%d", q, from)

	//for _, item := range result.Items {
	//log.Printf("Type:%s", reflect.TypeOf(item))
	//log.Printf("name:%s", profile.Name)
	//}
	log.Println(result.PreFrom)
	return result, nil
}

func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	return re.ReplaceAllString(q, "Payload.$1:")
}
