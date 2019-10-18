package indexer

// TorrentDetail ...
type TorrentDetail struct {
	Name       string   `json:"name"`
	DetailLink string   `json:"detail_link"`
	Size       string   `json:"size"`
	By         string   `json:"by"`
	Flags      []string `json:"flags"`
}

// TorrentBrowsingPage ...
type TorrentBrowsingPage struct {
	Index int
	Items []TorrentItem
	Pages []int
}

// TorrentItem ...
type TorrentItem interface {
	GetDetail() TorrentDetail
	GetMagnet() *string
}

// TorrentWebSiteCrawler ...
type TorrentWebSiteCrawler interface {
	GetPage(browsing string, index int) *TorrentBrowsingPage
	GetSettings() TorrentWebSite
}
