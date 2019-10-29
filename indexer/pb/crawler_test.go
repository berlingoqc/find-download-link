package pb_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/berlingoqc/find-download-link/indexer"
	"github.com/gocolly/colly"
)

func TestCrawler(t *testing.T) {
	c := colly.NewCollector()
	c.OnHTML(("table[id=searchResult] tr"), func(e *colly.HTMLElement) {
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
		println(detail.Name, " ", detail.Seeders, " ", detail.Leechers, detail.Size, detail.By)
		return
	})
	c.Visit("https://proxtpb.art/browse/207/0/7")
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
