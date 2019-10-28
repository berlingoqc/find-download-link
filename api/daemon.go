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

// GetAvailableCrawler ...
func (d *DaemonFindDownloadAPI) GetAvailableCrawler() []string {
	return indexer.GetAllCrawlerName()
}

// GetAvailableBrowsingForCrawler ...
func (d *DaemonFindDownloadAPI) GetAvailableBrowsingForCrawler(name string) []string {
	if r, err := indexer.GetAvailableBrowsingForCrawler(name); err == nil {
		return r
	} else {
		panic(err)
	}
}

// GetAllBrowsing ...
func (d *DaemonFindDownloadAPI) GetAllBrowsing() []string {
	return indexer.GetAllAvailableBrowsing()
}

// RemoveCrawlingRun ...
func (d *DaemonFindDownloadAPI) RemoveCrawlingRun(id float64) float64 {
	if err := d.db.RemoveCrawlingRun(int64(id)); err != nil {
		panic(err)
	}
	indexer.FeedBack("findDownload", "onArchiveDelete", id)
	return id
}

// GetCrawlingRunInfos ...
func (d *DaemonFindDownloadAPI) GetCrawlingRunInfos() []indexer.CrawlingRunInfo {
	if r, err := d.db.GetCrawlingRunInfos(); err == nil {
		return r
	} else {
		panic(err)
	}
}

// GetCrawlingRunDetail ...
func (d *DaemonFindDownloadAPI) GetCrawlingRunDetail(id float64) indexer.CrawlingRunInfo {
	if c, err := d.db.GetCrawlingRunDetail(int64(id)); err == nil {
		return *c
	} else {
		panic(err)
	}
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
func (d *DaemonFindDownloadAPI) StartCrawlerAfter(id float64, ending float64) indexer.CrawlingRunInfo {
	if c, err := indexer.StartCrawlRoutineAfter(int64(id), int(ending)); err == nil {
		return *c
	} else {
		panic(err)
	}
}

// GetCrawledInterval ...
func (d *DaemonFindDownloadAPI) GetCrawledInterval(crawler string) map[string][]indexer.IntervalCrawled {
	ma, err := indexer.GetCrawledInterval(d.db, crawler)
	if err != nil {
		panic(err)
	}
	return ma
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
