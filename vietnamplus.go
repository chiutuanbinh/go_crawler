package main

import (
	"crawler/mongodb"
	"crawler/util"
	"crawler/xtype"
	"encoding/json"
	"html"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func vietnamplusCrawler(wg *sync.WaitGroup) {
	defer wg.Done()

	articleChan := make(chan *xtype.Article, 10)
	c := colly.NewCollector(
		colly.AllowedDomains("www.vietnamplus.vn"),
		colly.CacheDir("./vietnamplus_cache"),
		colly.MaxDepth(10),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		// log.Printf("Link found: %q -> %s\n", e.Text, e.Request.AbsoluteURL(link))
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		if strings.HasSuffix(link, "vnp") {
			// log.Printf("Link found: %q -> %s\n", e.Text, e.Request.AbsoluteURL(link))
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}

	})

	c.OnHTML(".cms-body", func(e *colly.HTMLElement) {
		if mongodb.Exist(e.Request.URL.String()) {
			return
		}
		log.Printf("article found %v\n", e.Request.URL)
		article := &xtype.Article{}
		content := html.UnescapeString(e.ChildText("p"))
		title, exist := e.DOM.ParentsUntil("div[data-tile]").Parent().Attr("data-tile")
		if !exist {
			title = ""
		}
		x := e.DOM.ParentsUntil("html").Parent().ChildrenFiltered("head").ChildrenFiltered("script[type='application/ld+json']")
		x.Each(func(index int, item *goquery.Selection) {
			var dat map[string]interface{}
			json.Unmarshal([]byte(item.Text()), &dat)

			val, ok := dat["datePublished"]
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
			// log.Println(publishTs)
			article.Meta.PublishTs = uint64(publishTs)

		})
		article.Title = title
		article.URI = e.Request.URL.String()
		// article.Meta.PublishTs = uint64(publishTs)
		article.Content.Parts = make([]interface{}, 0)
		for _, p := range strings.Split(content, "\n") {
			paragraph := xtype.Paragraph{}
			paragraph.Content = p
			article.Content.Parts = append(article.Content.Parts, paragraph)
		}
		article.ID = util.GenNextUUID()
		article.Publisher = "VIETNAMPLUS"
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

	c.Visit("https://www.vietnamplus.vn")
	log.Printf("TAH FUK")
	// c.Visit("https://vietnamnet.vn/vn/giai-tri/the-gioi-sao/sao-viet-3-9-tuoi-38-khanh-thi-nong-bong-ben-chong-tre-va-hai-con-671118.html")

}
