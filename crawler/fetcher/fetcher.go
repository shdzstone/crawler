package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"stone/go/crawler/config"
	"strings"
	"time"

	"golang.org/x/text/encoding/unicode"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

/*
1.error返回
* error.New()
* fmt.Errorf()

2.http库使用
* http.NewRequest()新建请求
* http.DefaultClient.Do()新建客户端
* httputil
*/

//fetcher：根据url获取对应的数据

var (
	rateLimiter    = time.Tick(time.Second / config.Qps)
	verboseLogging = false

	userAgent = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
		/*
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 OPR/26.0.1656.60",
			"Opera/8.0 (Windows NT 5.1; U; en)",
			"Mozilla/5.0 (Windows NT 5.1; U; en; rv:1.8.1) Gecko/20061208 Firefox/2.0.0 Opera 9.50",
			"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; en) Opera 9.50",
			"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:34.0) Gecko/20100101 Firefox/34.0",
			"Mozilla/5.0 (X11; U; Linux x86_64; zh-CN; rv:1.9.2.10) Gecko/20100922 Ubuntu/10.10 (maverick) Firefox/3.6.10",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.57.2 (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11",
			"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.133 Safari/534.16",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.11 TaoBrowser/2.0 Safari/536.11",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.71 Safari/537.1 LBBROWSER",
			"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; LBBROWSER)",
			"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E; LBBROWSER)",
			"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; QQBrowser/7.0.3698.400)",
			"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E)",
			"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.84 Safari/535.11 SE 2.X MetaSr 1.0",
			"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SV1; QQDownload 732; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.4.3.4000 Chrome/30.0.1599.101 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.122 UBrowser/4.0.3214.0 Safari/537.36",
		*/
	}
	cookie = "FSSBBIl1UgzbN7NO=5gCvqsLd.y7ESHwTweFWlM1Tj7eEHah1g8Jy_hKx1ma4c22d8Gyq8ay3vQ8hRk2NHC4k7fPyOkH_6_fscQXYkoG; sid=60397cbc-4b85-4287-a960-62be42f3d4fa; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1604056747; ec=W7bffEnW-1604056748766-f85c0b65c0bc81241206384; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1604130022; _exid=6mjrsk%2F8YXlYTgHgOIaaw%2B0Lb%2FJXHwIIqpidUc9PUf0LuES6mYnCOJApEO7DRKt%2BZIVK4X6EcFHxCFb6SJmZRw%3D%3D; _efmdata=3CMKFMaY57HyYXGg24S36v6IdoEVAgrFuzMUkO%2BseVvaVnb0%2BpGT3pQtStoMlqgv1d5FSLMRJG3QLb3nSiQogh78Bvm0TYfAYalOxW1HBtA%3D; FSSBBIl1UgzbN7NP=5Uy0YDT5wmpLqqqmTIRNTiGTJCl8x4o166wBe7XrLlWP5subGIauE8dmqXOhO4UFnTAvz3yS3QGC0IkyrGOtFQPdghYQUm19BLu4uJz6XFDWopDchZDMmzSUqUSMUPhMh6HLDELP3kDnpwoa4JZ6VGSj8ikbdQ0Q49YMMRFox627nrEIclEQlmkiFvIDP_xdmwCZvF5mC9y0OirGJIEyrRvXomUuk69ATDmMDIcOn_r0ctMOgvIFtQZ_s4C1qsIjlZ"
)

func Fetch(url string) ([]byte, error) {
	//<-rateLimiter

	//反爬
	newUrl := strings.Replace(url, "http://", "https://", 1)

	request, _ := http.NewRequest(http.MethodGet, newUrl, nil)
	//珍爱网做了简单的反爬措施，使用user-agent和cookie的技术突破
	request.Header.Add("user-agent", userAgent[rand.Intn(len(userAgent))])
	request.Header.Add("cookie", cookie)

	resp, _ := http.DefaultClient.Do(request)
	defer resp.Body.Close()
	log.Printf("Fetching url:%v", request.URL)

	//httputil.DumpResponse()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code %d", resp.StatusCode)
	}

	//如果页面传来的不是utf8，需要转为utf8格式
	bodyReader := bufio.NewReader(resp.Body)
	encoding := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(resp.Body, encoding.NewDecoder())
	//utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())

	return ioutil.ReadAll(utf8Reader)
}

//使用charset.DetermineEncoding()判断字符编码
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(128)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	encoding, _, _ := charset.DetermineEncoding(bytes, "")

	return encoding
}

func SetVerboseLogging() {
	verboseLogging = true
}
