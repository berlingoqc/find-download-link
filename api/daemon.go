package api

import (
	"github.com/berlingoqc/find-download-link/indexer"
	// nlab
	_ "github.com/berlingoqc/find-download-link/indexer/pb"
)

// DaemonFindDownloadAPI ...
type DaemonFindDownloadAPI struct {
	db *indexer.DownloadDB
}

// GetActiveCrawler ...
func (d *DaemonFindDownloadAPI) GetActiveCrawler() []indexer.CrawlingRunInfo {
	return indexer.GetActiveCrawler()
}

// StopActiveCrawler ...
func (d *DaemonFindDownloadAPI) StopActiveCrawler(id float64) {
	if err := indexer.StopCrawlRoutine(int64(id)); err != nil {
		panic(err)
	}
}

// StartCrawler ...
func (d *DaemonFindDownloadAPI) StartCrawler(crawler, browsing string, start, end float64) indexer.CrawlingRunInfo {
	if info, err := indexer.StartCrawlRoutine(crawler, browsing, int(start), int(end)); err == nil {
		return *info
	} else {
		panic(err)
	}
}

// RemoveCrawler ...
func (d *DaemonFindDownloadAPI) RemoveCrawler(id float64) []indexer.CrawlingRunInfo {
	if err := indexer.RemoveCrawlRoutine(int64(id)); err != nil {
		panic(err)
	}
	return indexer.GetActiveCrawler()
}

// StartCrawlerAfter ...
func (d *DaemonFindDownloadAPI) StartCrawlerAfter(id float64) indexer.CrawlingRunInfo {
	if c, err := indexer.StartCrawlRoutineAfter(int64(id)); err == nil {
		return c
	} else {
		panic(err)
	}
}

// GetDaemonFindDownloadAPI ...
func GetDaemonFindDownloadAPI() (*DaemonFindDownloadAPI, error) {
	db, err := indexer.GetDownloadDB()
	if err != nil {
		return nil, err
	}
	return &DaemonFindDownloadAPI{
		db: db,
	}, nil

}
