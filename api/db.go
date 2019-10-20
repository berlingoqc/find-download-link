package api

import (
	"github.com/berlingoqc/find-download-link/indexer"
)

type FindDownloadAPI struct {
	db *indexer.DownloadDB
}

func (f *FindDownloadAPI) GetEntityName() []string {
	return f.db.GetEntityName()
}

func (f *FindDownloadAPI) GetEntity(name string) indexer.Entity {
	if e, err := f.db.GetEntity(name); err == nil {
		return e
	} else {
		panic(err)
	}

}

func GetFindDownloadAPI() *indexer.DownloadDB {
	return &FindDownloadAPI{
		db: indexer.GetDownloadDB(),
	}
}
