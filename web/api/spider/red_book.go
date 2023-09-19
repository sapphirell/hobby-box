package spider

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sukitime.com/v2/web/api"
	"time"
)

func GetDownloadList(ctx *gin.Context) {
	link := FindLinkByShare(testShareText)
	u, err := url.Parse(link)
	if err != nil {
		api.Base.Failed(ctx, "无法解析的URL")
		return
	}

	c := colly.NewCollector()
	// 超时设定
	c.SetRequestTimeout(100 * time.Second)
	// 指定Agent信息
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		// Request头部设定
		r.Headers.Set("Host", u.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", u.Host)
		r.Headers.Set("Referer", link)
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")
	})
	c.OnResponse(func(resp *colly.Response) {
		fmt.Println("response received", resp.StatusCode)
		// goquery直接读取resp.Body的内容
		htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
		if err != nil {
			log.Printf("解析%s网址内容时发生错误:%s", link, err.Error())
			api.Base.Failed(ctx, "解析网址内容时发生异常")
			return
		}
		var backgroundsMap = make(map[string]int)
		var bgs = make([]string, 0)
		htmlDoc.Find(".swiper-wrapper div").Each(func(i int, s *goquery.Selection) {
			style, ok := s.Attr("style")
			findBgCompile, _ := regexp.Compile(`(https(.*?));`)
			if ok {
				//筛选出真正的URL
				compileFindString := findBgCompile.FindString(style)
				compileFindString = strings.TrimRight(compileFindString, ");")
				if len(compileFindString) > 0 {
					backgroundsMap[compileFindString] = 1
				}
			}
		})
		for v, _ := range backgroundsMap {
			bgs = append(bgs, v)
		}
		api.Base.Success(ctx, bgs)
		return
	})
	c.OnError(func(resp *colly.Response, errHttp error) {
		log.Printf("解析%s网址时发生错误:%s", link, err.Error())
		api.Base.Failed(ctx, "无法解析的网址内容..")
		return
	})
	err = c.Visit(link)
	return
}

// FindLinkByShare 根据小红书分享连接接取分享地址
// 测试用
var testShareText = "97 ちいかわ情报站发布了一篇小红书笔记，快来看吧！ 😆 uUmSHpfne4zXSfV 😆 http://xhslink.com/0nBwqu，复制本条信息，打开【小红书】App查看精彩内容！"

func FindLinkByShare(share string) string {
	compile, _ := regexp.Compile("(http(.*?))，")
	compileFindString := compile.FindString(testShareText)
	compileFindString = strings.TrimRight(compileFindString, "，")
	return compileFindString
}
