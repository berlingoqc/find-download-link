package api

import (
	"github.com/berlingoqc/find-download-link/indexer"
)

// FindDownloadAPI ...
type FindDownloadAPI struct {
	db *indexer.DownloadDB
}

// GetEntityName ...
func (f *FindDownloadAPI) GetEntityName() []string {
	s, err := f.db.GetEntityName()
	if err != nil {
		panic(err)
	}
	return s
}

// GetEntity ...
func (f *FindDownloadAPI) GetEntity(name string) indexer.Entity {
	if e, err := f.db.GetEntity(name); err == nil {
		return e
	} else {
		panic(err)
	}

}

// GetFindDownloadAPI ...
func GetFindDownloadAPI() (*FindDownloadAPI, error) {
	db, err := indexer.GetDownloadDB()
	if err != nil {
		return nil, err
	}
	return &FindDownloadAPI{
		db: db,
	}, err
}
