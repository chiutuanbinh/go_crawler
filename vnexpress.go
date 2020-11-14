package main

import (
	"binhct/common/message"
	"binhct/common/xtype"
	"crawler/mongodb"
	"crawler/util"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

func vnexpressCrawler(wg *sync.WaitGroup) {
	defer wg.Done()
	poster := message.CreatePoster("127.0.0.1:9092", "news")
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("vnexpress.net"),
		colly.CacheDir("./vnexpress_cache"),
		colly.MaxDepth(10),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		// log.Printf("Link found: %q -> %s\n", e.Text, e.Request.AbsoluteURL(link))
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		if strings.HasSuffix(link, "html") {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnHTML("article", func(e *colly.HTMLElement) {
		if mongodb.Exist(e.Request.URL.String()) {
			return
		}
		log.Printf("article found %v\n", e.Request.URL)
		article := &xtype.Article{}
		content := e.ChildText("p[class]")
		title := e.DOM.ParentsUntil("html").Parent().ChildrenFiltered("head").ChildrenFiltered("title").Text()
		val, ok := e.DOM.ParentsUntil("html").Parent().ChildrenFiltered("head").ChildrenFiltered("meta[name=pubdate]").Attr("content")

		var publishTs int64
		if !ok {
			publishTs = 0
		} else {
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				publishTs = 0
			} else {
				publishTs = t.Unix()
			}
		}
		// log.Printf("time stamp %v\n", publishTs)
		// log.Printf("TITLE %v\n", e.DOM.ParentsUntil("html").Parent().ChildrenFiltered("head").ChildrenFiltered("title").Text())
		article.Title = title
		article.URI = e.Request.URL.String()
		article.Meta.PublishTs = uint64(publishTs)
		article.Content.Parts = make([]interface{}, 0)
		for _, p := range strings.Split(content, "\n") {
			paragraph := xtype.Paragraph{}
			paragraph.Content = p
			article.Content.Parts = append(article.Content.Parts, paragraph)
		}
		article.ID = util.GenNextUUID()
		article.Publisher = "VNEXPRESS"
		// log.Printf("%+v\n", article.URI)
		mongodb.Add(article.ID, article)
		poster.Post("news", "id", article.ID)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://vnexpress.net")

}
