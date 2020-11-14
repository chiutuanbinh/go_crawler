package main

import (
	"binhct/common/xtype"
	"crawler/mongodb"
	"crawler/util"
	"encoding/json"
	"html"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func vietnamnetCrawler(wg *sync.WaitGroup) {
	defer wg.Done()

	articleChan := make(chan *xtype.Article, 10)
	c := colly.NewCollector(
		colly.AllowedDomains("vietnamnet.vn"),
		colly.CacheDir("./vietnamet_cache"),
		colly.MaxDepth(10),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		// log.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		if strings.Contains(link, "html") {
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}

	})

	c.OnHTML(`div[id=ArticleContent]`, func(e *colly.HTMLElement) {
		if mongodb.Exist(e.Request.URL.String()) {
			return
		}
		log.Printf("article found %v\n", e.Request.URL)
		article := &xtype.Article{}
		content := html.UnescapeString(e.ChildText("p"))
		title := e.DOM.ParentsUntil("html").Parent().ChildrenFiltered("head").ChildrenFiltered("title").Text()
		x := e.DOM.ParentsUntil("html").Parent().ChildrenFiltered("head").ChildrenFiltered("script[type='application/ld+json']")
		x.Each(func(index int, item *goquery.Selection) {
			var dat map[string]interface{}
			json.Unmarshal([]byte(item.Text()), &dat)

			val, ok := dat["datePublished"]
			// log.Printf("tims %v %v\n", val, ok)
			if !ok {
				return
			}
			t, err := time.Parse(time.RFC3339, val.(string))
			var publishTs int64
			if err != nil {
				publishTs = 0
			} else {
				publishTs = t.Unix()
			}

			article.Meta.PublishTs = uint64(publishTs)

		})
		article.Title = title
		article.URI = e.Request.URL.String()
		article.Content.Parts = make([]interface{}, 0)
		for _, p := range strings.Split(content, "\n") {
			paragraph := xtype.Paragraph{}
			paragraph.Content = p
			article.Content.Parts = append(article.Content.Parts, paragraph)
		}
		article.ID = util.GenNextUUID()
		article.Publisher = "VIETNAMNET"
		articleChan <- article
	})

	go func() {
		for job := range articleChan {
			if mongodb.Exist(job.URI) {
				continue
			}
			mongodb.Add(job.ID, job)
		}
	}()
	c.Visit("https://vietnamnet.vn")
	// c.Visit("https://vietnamnet.vn/vn/giai-tri/the-gioi-sao/sao-viet-3-9-tuoi-38-khanh-thi-nong-bong-ben-chong-tre-va-hai-con-671118.html")

}
