package main

import (
	"github.com/berlingoqc/find-download-link/indexer"

	// add crawler for the pirate bay
	_ "github.com/berlingoqc/find-download-link/indexer/pb"
)

func main() {
	err := indexer.CrawlBrowsing("piratebay", "hdmovie", 0, 0)
	if err != nil {
		panic(err)
	}
}
