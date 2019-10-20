package api

import "github.com/berlingoqc/find-download-link/indexer"

type DaemonFindDownloadAPI struct {
	db *indexer.DownloadDB
}

func (d *DaemonFindDownloadAPI) GetActiveCrawler() []indexer.CrawlingRunInfo {
	return indexer.GetActiveCrawler()
}

func (d *DaemonFindDownloadAPI) StopActiveCrawler(crawler, browsing string) {
	if err := indexer.StopCrawlRoutine(crawler, browsing); err != nil {
		panic(err)
	}
}

func GetDaemonFindDownloadAPI() (*DaemonFindDownloadAPI, error) {
	db, err := indexer.GetDownloadDB()
	if err != nil {
		return nil, err
	}
	return &DaemonFindDownloadAPI{
		db: db,
	}, nil

}
