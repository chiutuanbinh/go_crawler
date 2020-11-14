package main

import "sync"

func checkAndTrack(absoluteLink string) {

}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go vnexpressCrawler(&wg)
	// go vietnamnetCrawler(&wg)
	// go vietnamplusCrawler(&wg)
	wg.Wait()

}
