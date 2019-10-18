package main

import (
	"time"

	"github.com/berlingoqc/find-download-link/indexer"
	"github.com/berlingoqc/find-download-link/indexer/pb"
)

func main() {

	crawler := pb.GetCrawler()

	browsing := "hdmovie"

	criteria := crawler.WebSite.Browsings[browsing].Criteria
	ch := make(chan indexer.Record, 5)
	if err := startDbRoutine(criteria,ch); err != nil {
		panic(err)
	}

	index := 3
	for {
		page := crawler.GetPage(browsing, index)
		for _, item := range page.Items {
			magnet := item.GetMagnet()
			ch <- indexer.Record{
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
}

func startDbRoutine(criteria indexer.TorrentCriteria, ch chan indexer.Record) error {
	db ,err := indexer.GetDownloadDB()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			record := <- ch
			entityName, tags := indexer.ExtractNameAndTag(record.Detail.Name,criteria.RequiredKeyword)
			record.Detail.Flags = tags;
			if _, err := db.AddRecordEntity(entityName,[]indexer.Record{record}); err != nil {
				println("ERROR ",err)
				continue
			}
			println("Add record for "+entityName)
		}

	}()
	return err
}
