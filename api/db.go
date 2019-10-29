package api

import (
	"github.com/berlingoqc/find-download-link/indexer"
	"github.com/mitchellh/mapstructure"
)

// PagingSearch ...
type PagingSearch struct {
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	OrderBy  []string `json:"orderBy"`
	Query    string   `json:"query"`
	Category string   `json:"category"`
}

// PagingInfo ...
type PagingInfo struct {
	ItemCount int `json:"item_count"`
}

// FindDownloadAPI ...
type FindDownloadAPI struct {
	db *indexer.DownloadDB
}

// GetEntityName ...
func (f *FindDownloadAPI) GetEntityName(in map[string]interface{}) []string {
	p := &PagingSearch{}
	err := mapstructure.Decode(in, p)
	if err != nil {
		panic(err)
	}
	s, err := f.db.GetEntityName(p.Query, p.Category, p.Limit, p.Offset)
	if err != nil {
		panic(err)
	}
	return s
}

// GetPagingInfo ...
func (f *FindDownloadAPI) GetPagingInfo() PagingInfo {
	i, err := indexer.GetCountOfTable(f.db.Db, "query_entity")
	if err != nil {
		panic(err)
	}
	return PagingInfo{
		ItemCount: i,
	}
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
