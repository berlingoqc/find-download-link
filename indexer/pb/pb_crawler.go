package pb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/berlingoqc/find-download-link/indexer"
	"github.com/gocolly/colly"
)

// TorrentItem ...
type TorrentItem struct {
	collector *colly.Collector
	Detail    indexer.TorrentDetail
}

// GetDetail ...
func (i *TorrentItem) GetDetail() indexer.TorrentDetail {
	return i.Detail
}

// GetMagnet ...
func (i *TorrentItem) GetMagnet() *string {
	var magnet string
	i.collector.OnHTML("div[class=download]", func(e *colly.HTMLElement) {
		magnet = e.ChildAttr("a", "href")
	})
	i.collector.Visit(i.Detail.DetailLink)
	return &magnet
}

func getTorrentItem(c *Crawler, e *colly.HTMLElement) *TorrentItem {
	e.Text = strings.TrimSpace(e.Text)
	return &TorrentItem{
		collector: c.Collector,
		Detail: indexer.TorrentDetail{
			Name:       e.Text,
			DetailLink: c.WebSite.URL + e.ChildAttr("a", "href"),
		},
	}
}

// Crawler ...
type Crawler struct {
	Collector *colly.Collector
	WebSite   indexer.TorrentWebSite
}

// GetSettings ...
func (p *Crawler) GetSettings() indexer.TorrentWebSite {
	return p.WebSite
}

// SetSettings ...
func (p *Crawler) SetSettings(t indexer.TorrentWebSite) {
	p.WebSite = t
}

// GetPage ...
func (p *Crawler) GetPage(browsing string, index int) *indexer.TorrentBrowsingPage {
	page := &indexer.TorrentBrowsingPage{
		Pages: []int{},
		Index: index,
		Items: []indexer.TorrentItem{},
	}
	p.Collector.OnHTML("div[class=detName]", func(e *colly.HTMLElement) {
		page.Items = append(page.Items, getTorrentItem(p, e))
	})
	if b, ok := p.WebSite.Browsings[browsing]; ok {
		url := p.WebSite.URL + "/browse/" + b.Path + "/" + strconv.Itoa(index) + "/7"
		println(url)
		p.Collector.Visit(url)
	}
	return page
}

// GetCrawler ...
func GetCrawler() *Crawler {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	return &Crawler{
		Collector: c,
		WebSite:   indexer.TorrentWebSite{},
	}

}

func init() {
	indexer.AddCrawler("piratebay", GetCrawler())
}
