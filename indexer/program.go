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

func CrawlBrowsing(crawlerName, browsing string, startingIndex, endingIndex int) error {
	var ok bool
	var crawler TorrentWebSiteCrawler
	if crawler, ok = crawlers[crawlerName]; !ok {
		return errors.New("Not crawler found for " + crawlerName)
	}

	criteria := crawler.GetSettings().Browsings[browsing].Criteria
	ch := make(chan Record, 5)
	if err := startDbRoutine(criteria, ch); err != nil {
		panic(err)
	}

	index := startingIndex
	for {
		page := crawler.GetPage(browsing, index)
		for _, item := range page.Items {
			magnet := item.GetMagnet()
			ch <- Record{
				Source:    "piratebay",
				Detail:    item.GetDetail(),
				ExtractOn: time.Now().Unix(),
				Link:      *magnet,
			}
			time.Sleep(15 * time.Second)
		}
		index++
		if index == 30 {
			break
		}
	}
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
