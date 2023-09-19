package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func main() {
	link := FindLinkByShare(testShareText)
	log.Println("link:", link)

	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
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
	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("title:", e.Text)
	})
	c.OnResponse(func(resp *colly.Response) {
		fmt.Println("response received", resp.StatusCode)
		// goquery直接读取resp.Body的内容
		htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
		// 读取url再传给goquery，访问url读取内容，此处不建议使用
		// htmlDoc, err := goquery.NewDocument(resp.Request.URL.String())
		if err != nil {
			log.Fatal(err)
		}
		// 找到抓取项 <div class="hotnews" alog-group="focustop-hotnews"> 下所有的a解析
		htmlDoc.Find(".swiper-wrapper div").Each(func(i int, s *goquery.Selection) {
			//band, _ := s.Attr("href")
			//title := s.Text()
			//fmt.Printf("热点新闻 %d: %s - %s\n", i, title, band)
			//c.Visit(band)
			fmt.Println(s.Attr("style"))
		})
	})
	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})
	err = c.Visit(link)
}

// FindLinkByShare 根据小红书分享连接接取分享地址
var testShareText = "97 ちいかわ情报站发布了一篇小红书笔记，快来看吧！ 😆 uUmSHpfne4zXSfV 😆 http://xhslink.com/0nBwqu，复制本条信息，打开【小红书】App查看精彩内容！"

func FindLinkByShare(share string) string {
	compile, _ := regexp.Compile("(http(.*?))，")
	compileFindString := compile.FindString(testShareText)
	compileFindString = strings.TrimRight(compileFindString, "，")
	return compileFindString
}
