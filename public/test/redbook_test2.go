package main

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"regexp"
	"strings"
	"time"
)

func main() {
	html := GetHttpHTML("https://www.xiaohongshu.com/explore/65313791000000002500b8dd?app_platform=ios&app_version=8.10.1&author_share=2&share_from_user_hidden=true&type=normal&xhsshare=CopyLink&appuid=607a9f4e00000000010017e1&apptime=1698223402",
		"#noteContainer > div.media-container > div > div > div.swiper.swiper-initialized.swiper-horizontal.swiper-pointer-events.note-slider.narrower > div > div.swiper-slide.swiper-slide-active")
	log.Println(html)
}

func GetHttpHTML(url string, selector string) map[string]int {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), //设置成无浏览器弹出模式
		chromedp.Flag("blink-settings", "imageEnable=false"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"),
	}
	c, _ := chromedp.NewExecAllocator(context.Background(), options...)
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	_ = chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	timeOutCtx, cancel := context.WithTimeout(chromeCtx, 3*time.Second)
	defer cancel()
	log.Println("执行1")
	var htmlContent string
	err := chromedp.Run(timeOutCtx,
		chromedp.Navigate(url),
		//需要爬取的网页的url
		//chromedp.WaitVisible(`#content > div > section.fp-tournament-award-badge-carousel_awardBadgeCarouselSection__w_Ys5 > div > div > div.col-12.fp-tournament-award-badge-carousel_awardCarouselColumn__fQJLf.g-0 > div > div > div > div > div > div > div.slick-slide.slick-active.slick-current > div > div > div`),
		chromedp.WaitVisible(selector),
		//等待某个特定的元素出现
		chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
		//生成最终的html文件并保存在htmlContent文件中
	)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(htmlContent)
	htmlDoc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	log.Println("执行2")
	//awardTmp := [][]string{}
	//doc.Find(`div[class="fp-tournament-award-badge_awardContent__dUtoO"]`).Each(func(i int, selection *goquery.Selection) {
	//	award := selection.Find(`p[class="fp-tournament-award-badge_awardName__JpsZZ"]`).Text()
	//	//goquery通过Find()查找到我们选择的位置，Each()的功能与遍历相似返回所有的结果，Text()返回文本内容
	//	name := selection.Find(`h4[class=" fp-tournament-award-badge_awardWinner__P_z2d"]`).Text()
	//	country := selection.Find(`p[class="fp-tournament-award-badge_awardWinnerCountry__EmjVU"]`).Text()
	//
	//	awardTmp = append(awardTmp, []string{award, name, country})
	//})

	//if err != nil {
	//	log.Printf("解析%s网址内容时发生错误:%s", link, err.Error())
	//	api.Base.Failed(ctx, "解析网址内容时发生异常")
	//	return
	//}
	var backgroundsMap = make(map[string]int)
	//var bgs = make([]string, 0)

	htmlDoc.Find(".swiper-slide").Each(func(i int, s *goquery.Selection) {

		style, ok := s.Attr("style")
		//log.Println(style)
		findBgCompile, _ := regexp.Compile(`(http(.*?));`)
		if ok {
			//筛选出真正的URL
			compileFindString := findBgCompile.FindString(style)
			log.Println("待处理", compileFindString)
			compileFindString = strings.TrimRight(compileFindString, "\");")
			if len(compileFindString) > 0 {
				backgroundsMap[compileFindString] = 1
				log.Println("GetUrls:", compileFindString)
			}
		}
	})
	return backgroundsMap
}
