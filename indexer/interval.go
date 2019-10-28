package indexer

// IntervalCrawled ...
type IntervalCrawled struct {
	ID       int64 `json:"id"`
	Interval []int `json:"interval"`
}

// GetCrawledInterval ...
func GetCrawledInterval(db *DownloadDB, crawler string) (map[string][]IntervalCrawled, error) {
	rows, err := db.GetCrawlingRunInfosForCrawler(crawler)
	if err != nil {
		return nil, err
	}
	ret := make(map[string][]IntervalCrawled)
	for _, row := range rows {
		if _, ok := ret[row.Browsing]; !ok {
			ret[row.Browsing] = make([]IntervalCrawled, 0)
		}
		ret[row.Browsing] = append(ret[row.Browsing], IntervalCrawled{
			ID:       row.ID,
			Interval: []int{row.Starting, row.Current},
		})
	}
	return ret, nil
}
