package rby

import (
	"fmt"
	"strconv"

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
	return &magnet
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

	p.Collector.OnHTML(("table[class=lista2t] tr"), func(e *colly.HTMLElement) {
		detail := &indexer.TorrentDetail{}
		e.ForEach("td", func(ii int, eee *colly.HTMLElement) {
			switch ii {
			case 1:
				detail.Name = eee.ChildText("a")
				detail.DetailLink = eee.ChildAttr("a", "href")
				break
			case 2:
				break
			case 3:
				break
			case 4:
				break
			case 5:
				break
			}
		})
		page.Items = append(page.Items, &TorrentItem{
			collector: p.Collector,
			Detail:    *detail,
		})
		return
	})

	if b, ok := p.WebSite.Browsings[browsing]; ok {
		url := p.WebSite.URL + b.Path + "&order=seeders&by=DESC&page=" + strconv.Itoa(index)
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
	indexer.AddCrawler("rarby", GetCrawler())
}
