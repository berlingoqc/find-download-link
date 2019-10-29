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
	p.Collector.OnHTML(("table[id=searchResult] tr"), func(e *colly.HTMLElement) {
		detail := &indexer.TorrentDetail{}
		e.ForEach("td", func(ii int, eee *colly.HTMLElement) {
			switch ii {
			case 1:
				one(detail, eee)
				break
			case 2:
				two(detail, eee)
				break
			case 3:
				three(detail, eee)
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
		url := p.WebSite.URL + "/browse/" + b.Path + "/" + strconv.Itoa(index) + "/7"
		println(url)
		p.Collector.Visit(url)
	}
	return page
}

func one(d *indexer.TorrentDetail, eee *colly.HTMLElement) {
	d.Name = eee.ChildText("a")
	d.DetailLink = eee.ChildAttr("a", "href")
	eee.ForEach("font", func(nbr int, e *colly.HTMLElement) {
		text := e.Text
		items := strings.Split(text, ", ")
		if len(items) == 3 {
			d.Size = strings.Replace(items[1], "Size", "", 1)
			d.By = strings.Replace(items[2], "ULed by", "", 1)
		} else {
			println("ERROR DeTAIL TORRENT")
		}
	})
}
func two(d *indexer.TorrentDetail, eee *colly.HTMLElement) {
	d.Seeders, _ = strconv.Atoi(eee.Text)
}
func three(d *indexer.TorrentDetail, eee *colly.HTMLElement) {
	d.Leechers, _ = strconv.Atoi(eee.Text)
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
