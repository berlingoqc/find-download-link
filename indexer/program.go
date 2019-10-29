package indexer

import (
	"errors"
	"strconv"
	"time"
)

var (
	// FeedBack is to be implement to send message by the application
	FeedBack func(ns, event string, data interface{})
	crawlers map[string]TorrentWebSiteCrawler = make(map[string]TorrentWebSiteCrawler)
)

// GetCrawler ...
func GetCrawler(name string) TorrentWebSiteCrawler {
	return crawlers[name]
}

// GetAllCrawlerName ...
func GetAllCrawlerName() []string {
	var n []string
	for k := range crawlers {
		n = append(n, k)
	}
	return n
}

// GetAvailableBrowsingForCrawler ...
func GetAvailableBrowsingForCrawler(name string) ([]string, error) {
	if c, ok := crawlers[name]; ok {
		var r []string
		for n := range c.GetSettings().Browsings {
			r = append(r, n)
		}
		return r, nil
	}
	return nil, errors.New("Clawler doesnt exists")
}

// GetAllAvailableBrowsing ...
func GetAllAvailableBrowsing() []string {
	var r []string
	for _, v := range crawlers {
		for k := range v.GetSettings().Browsings {
			present := false
			for _, r := range r {
				if r == k {
					present = true
					break
				}
			}
			if !present {
				r = append(r, k)
			}
		}
	}
	return r
}

// AddCrawler ...
func AddCrawler(name string, crawler TorrentWebSiteCrawler) {
	crawlers[name] = crawler
}

// CrawlingRunInfo ...
type CrawlingRunInfo struct {
	Crawler   string                    `json:"crawler"`
	Browsing  string                    `json:"browsing"`
	Starting  int                       `json:"starting"`
	Ending    int                       `json:"ending"`
	Current   int                       `json:"current"`
	Status    string                    `json:"status"`
	ID        int64                     `json:"id"`
	PageInfos []TorrentPageCrawlingInfo `json:"pageinfos"`
	Error     string                    `json:"error"`
	SignalCh  chan interface{}          `json:"-"`
}

// TorrentPageCrawlingInfo ...
type TorrentPageCrawlingInfo struct {
	RunInfo     *CrawlingRunInfo `json:"-"`
	ExtractedOn int64            `json:"extracted_on"`
	ItemsAdded  []string         `json:"itemsAdded"`
}

// TimeoutCrawling ...
type TimeoutCrawling struct {
	BetweenRequest         int `json:"betweenRequest"`
	BetweenPages           int `json:"betweenPages"`
	AfterRejection         int `json:"aferRejection"`
	MaxRetryAfterRejection int `json:"maxRetryAfterRejection"`
}

var timeout TimeoutCrawling
var mapCrawlingInfo map[int64]*CrawlingRunInfo

// StartCrawlRoutine ...
func StartCrawlRoutine(crawlerName, browsing string, starting, ending int) (*CrawlingRunInfo, error) {
	name := time.Now().Unix()
	if _, ok := mapCrawlingInfo[name]; !ok {
		crawlingInfo := &CrawlingRunInfo{
			Crawler:   crawlerName,
			Browsing:  browsing,
			Starting:  starting,
			Ending:    ending,
			ID:        name,
			PageInfos: []TorrentPageCrawlingInfo{},
			SignalCh:  make(chan interface{}),
		}
		return crawlingInfo, crawlBrowsing(name, crawlingInfo)
	}
	return nil, errors.New("no crawler found for " + crawlerName)
}

// GetActiveCrawler ...
func GetActiveCrawler() []CrawlingRunInfo {
	a := []CrawlingRunInfo{}
	for _, i := range mapCrawlingInfo {
		a = append(a, *i)
	}
	return a
}

// StopCrawlRoutine ...
func StopCrawlRoutine(name int64) error {
	if r, ok := mapCrawlingInfo[name]; ok {
		r.SignalCh <- 0
		r.SignalCh <- 0
		r.SignalCh <- 0
		return nil
	}
	return errors.New("Crawling routine is not started")
}

// RemoveCrawlRoutine ...
func RemoveCrawlRoutine(id int64) error {
	if c, ok := mapCrawlingInfo[id]; ok {
		if c.Status == "stopped" {
			delete(mapCrawlingInfo, id)
			return nil
		}
		return errors.New("crawling routine is not stopped")
	}
	return errors.New("crawling routine doesnt exists")
}

// StartCrawlRoutineAfter ...
func StartCrawlRoutineAfter(id int64, ending int) (*CrawlingRunInfo, error) {
	db, err := GetDownloadDB()
	if err != nil {
		return nil, err
	}
	info, err := db.GetCrawlingRunDetail(id)
	if err != nil {
		return nil, err
	}
	return StartCrawlRoutine(info.Crawler, info.Browsing, info.Current, ending)

}

func crawlBrowsing(name int64, runInfo *CrawlingRunInfo) error {
	var ok bool
	var crawler TorrentWebSiteCrawler
	if crawler, ok = crawlers[runInfo.Crawler]; !ok {
		return errors.New("Not crawler found for " + runInfo.Crawler)
	}

	updateMap(name, runInfo)

	criteria := GetCriteria(runInfo.Browsing)
	ch := make(chan Record, 5)
	if err := startDbRoutine(runInfo.Browsing, criteria, ch, runInfo.SignalCh); err != nil {
		return err
	}

	timeoutMultiplier := 1

	go func() {
		index := runInfo.Starting
		runInfo.Status = "running"
		defer func() {
			if r := recover(); r != nil {
				runInfo.Error = (r.(error)).Error()
				runInfo.Status = "stopped"
				db, err := GetDownloadDB()
				if err != nil {
					runInfo.Error = err.Error()
				} else {
					if err = db.SaveCrawlingJob(runInfo); err != nil {
						runInfo.Error = err.Error()
					} else {
						runInfo.Status = "archived"
						updateMap(name, runInfo)
						delete(mapCrawlingInfo, name)
						return
					}
				}
				updateMap(name, runInfo)
			}
		}()
		for {
			select {
			case _ = <-runInfo.SignalCh:
				panic(errors.New("manually stopped"))
			default:
				page := crawler.GetPage(runInfo.Browsing, index)
				runInfo.Current = index
				runInfo.PageInfos = append(runInfo.PageInfos, TorrentPageCrawlingInfo{
					ExtractedOn: time.Now().Unix(),
					ItemsAdded:  []string{},
				})
				updateMap(name, runInfo)
				if len(page.Items) == 0 {
					sleepWithChannel(time.Duration(timeout.AfterRejection*timeoutMultiplier)*time.Millisecond, runInfo.SignalCh)
					timeoutMultiplier++
					if timeoutMultiplier == timeout.MaxRetryAfterRejection {
						panic(errors.New("max timeout without data at index " + strconv.Itoa(index)))
					}
				}
				for _, item := range page.Items {
					sleepWithChannel(time.Duration(timeout.BetweenRequest)*time.Millisecond, runInfo.SignalCh)
					magnet := item.GetMagnet()
					ch <- Record{
						Source:    runInfo.Crawler,
						Detail:    item.GetDetail(),
						ExtractOn: time.Now().Unix(),
						Link:      *magnet,
					}
					indexPage := len(runInfo.PageInfos) - 1
					runInfo.PageInfos[indexPage].ItemsAdded = append(runInfo.PageInfos[indexPage].ItemsAdded, item.GetDetail().Name)
				}
				updateMap(name, runInfo)
				sleepWithChannel(time.Duration(timeout.BetweenPages)*time.Millisecond, runInfo.SignalCh)
				index++
				if index == runInfo.Ending {
					break
				}
			}
		}
	}()
	return nil
}

func updateMap(name int64, info *CrawlingRunInfo) {
	mapCrawlingInfo[name] = info
	FeedBack("findDownload", "onCrawlingInfoUpdate", info)
}

func sleepWithChannel(duration time.Duration, stopChannel chan interface{}) {
	select {
	case _ = <-stopChannel:
		panic(errors.New("Stop manually"))
	case <-time.After(duration):
		return
	}
}

func startDbRoutine(browsing string, criteria TorrentCriteria, ch chan Record, chOver chan interface{}) error {
	db, err := GetDownloadDB()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case _ = <-chOver:
				return
			case record := <-ch:
				entityName, tags := ExtractNameAndTag(record.Detail.Name, criteria.Tags)
				record.Detail.Flags = tags
				if !db.EntityExists(entityName) {
					if _, err := db.CreateEntity(browsing, entityName); err != nil {
						println("ERROR ", err.Error())
					}
				}
				if _, err := db.AddRecordEntity(entityName, []Record{record}); err != nil {
					println("ERROR ", err.Error())
					continue
				}
			}
		}
	}()
	return err
}

func init() {
	mapCrawlingInfo = make(map[int64]*CrawlingRunInfo)
}
