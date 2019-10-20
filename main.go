package main

import (
	"flag"

	"github.com/berlingoqc/find-download-link/indexer"

	// add crawler for the pirate bay
	_ "github.com/berlingoqc/find-download-link/indexer/pb"
)

func main() {
	dbFilePtr := flag.String("db", "./data.db", "path of the sqlite database to use")
	dbOverWritePtr := flag.Bool("overwriteDb", false, "")
	crawlerPtr := flag.String("crawler", "piratebay", "a string")
	browsingPtr := flag.String("browsing", "", "a string")
	startIndex := flag.Int("start", 0, "a int")
	endingIndex := flag.Int("end", 50, "a int")

	flag.Parse()

	indexer.DbFile = *dbFilePtr
	indexer.OverWriteDb = *dbOverWritePtr

	info, err := indexer.StartCrawlRoutine(*crawlerPtr, *browsingPtr, *startIndex, *endingIndex)
	if err != nil {
		panic(err)
	}
	<-info.SignalCh
}
