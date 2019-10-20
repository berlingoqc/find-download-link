package indexer

import (
	"errors"
	"time"
)

var (
	crawlers map[string]TorrentWebSiteCrawler = make(map[string]TorrentWebSiteCrawler)
)

func GetCrawler(name string) TorrentWebSiteCrawler {
	return crawlers[name]
}

func AddCrawler(name string, crawler TorrentWebSiteCrawler) {
	crawlers[name] = crawler
}

type CrawlingRunInfo struct {
	Crawler    string
	Browsing   string
	Starting   int
	Ending     int
	ItemsAdded []string
	SignalCh   chan interface{}
}

var mapCrawlingInfo map[string]*CrawlingRunInfo

func StartCrawlRoutine(crawlerName, browsing string, starting, ending int) (*CrawlingRunInfo, error) {

	name := crawlerName + ":" + browsing

	if _, ok := mapCrawlingInfo[name]; !ok {

	} else {
		return nil, errors.New("")
	}

	crawlingInfo := &CrawlingRunInfo{
		Crawler:    crawlerName,
		Browsing:   browsing,
		Starting:   starting,
		Ending:     ending,
		ItemsAdded: []string{},
		SignalCh:   make(chan interface{}),
	}

	return crawlingInfo.SignalCh, crawlBrowsing(crawlingInfo)
}

func StopCrawlRoutine() {
	if crawlingInfo != nil {
	}
	panic("Crawling routine is not started")
}

func crawlBrowsing(runInfo *CrawlingRunInfo) error {
	println("CRAWLER ", runInfo.Crawler, " BROWSING ", runInfo.Browsing, " STARTING ", runInfo.Starting, " ENDING ", runInfo.Ending)
	var ok bool
	var crawler TorrentWebSiteCrawler
	if crawler, ok = crawlers[runInfo.Crawler]; !ok {
		return errors.New("Not crawler found for " + runInfo.Crawler)
	}

	criteria := crawler.GetSettings().Browsings[runInfo.Browsing].Criteria
	ch := make(chan Record, 5)
	if err := startDbRoutine(criteria, ch); err != nil {
		return err
	}

	timeoutTime := 2

	go func() {
		index := runInfo.Starting
		for {
			page := crawler.GetPage(runInfo.Browsing, index)
			if len(page.Items) == 0 {
				println("No item return in the page, waiting timeout of ", timeoutTime)
				time.Sleep(time.Duration(timeoutTime) * time.Minute)
			}
			for _, item := range page.Items {
				magnet := item.GetMagnet()
				ch <- Record{
					Source:    runInfo.Crawler,
					Detail:    item.GetDetail(),
					ExtractOn: time.Now().Unix(),
					Link:      *magnet,
				}
			}
			time.Sleep(15 * time.Second)
			index++
			if index == runInfo.Ending {
				break
			}
		}
		runInfo.SignalCh <- 0
	}()
	return nil
}

func startDbRoutine(criteria TorrentCriteria, ch chan Record) error {
	db, err := GetDownloadDB()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			record := <-ch
			entityName, tags := ExtractNameAndTag(record.Detail.Name, criteria.RequiredKeyword)
			record.Detail.Flags = tags
			if _, err := db.AddRecordEntity(entityName, []Record{record}); err != nil {
				println("ERROR ", err)
				continue
			}
			println("Add record for " + entityName)
		}

	}()
	return err
}
