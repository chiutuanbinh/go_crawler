package main

import (
	"crawler/nlp"
	"log"
)

func checkAndTrack(absoluteLink string) {

}

func main() {
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go vnexpressCrawler(&wg)
	// go vietnamnetCrawler(&wg)
	// go vietnamplusCrawler(&wg)
	// wg.Wait()
	z := nlp.NLPExtract("ong Phan Ngọc Thọ, Chủ tịch tỉnh Thừa Thiên Huế, cho biết như trên, chiều 17/10.")
	log.Println(z)
}
