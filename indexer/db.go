package indexer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"

	// drivers for the database use in the applications
	_ "github.com/mattn/go-sqlite3"
)

const (
	createTable = `
	CREATE TABLE IF NOT EXISTS entity (
		uid INTEGER PRIMARY KEY AUTOINCREMENT,
		type VARCHAR(50),
		name VARCHAR(255),
		data BLOB
	);

	CREATE TABLE IF NOT EXISTS crawlingrun (
		id INTEGER PRIMARY KEY,
		crawler VARCHAR(50),
		browsing VARCHAR(50),
		start INTEGER,
		end INTEGER,
		data BLOB
	);
	`
	addEntity        = `INSERT INTO entity(type,name,data) values(?,?,?)`
	updateEntityData = `UPDATE entity SET data = ? WHERE name = ?`

	addCrawlingRun = `INSERT INTO crawlingrun values(?,?,?,?,?,?)`

	getEntityName = `SELECT name FROM entity;`
	getEntity     = `SELECT name, data FROM entity WHERE name = ?`

	getCrawlingJobName = `SELECT id, crawler, browsing, start, end FROM crawlingrun`
	getCrawlingJob     = `SELECT data FROM crawlingrun WHERE id = ?`

	deleteCrawlingJob = `DELETE FROM crawlingrun WHERE id = ?`
)

// DBSettings ...
type DBSettings struct {
	Path      string `json:"path"`
	OverWrite bool   `json:"overwrite"`
}

var dbsettings DBSettings

// Entity ...
type Entity struct {
	Name    string   `json:"name"`
	Records []Record `json:"records"`
}

// Record ...
type Record struct {
	ExtractOn int64         `json:"extract_on"`
	Source    string        `json:"source"`
	Detail    TorrentDetail `json:"detail"`
	Link      string        `json:"link"`
}

// CrawlingRun ...
type CrawlingRun struct {
	ExtractOn int64
	Crawler   string
	Browsing  string
	Starting  int
	Ending    int
}

// DownloadDB ...
type DownloadDB struct {
	Db *sql.DB
}

// GetCrawlingRunInfos ...
func (d *DownloadDB) GetCrawlingRunInfos() (a []CrawlingRunInfo, e error) {
	var rows *sql.Rows
	rows, e = d.Db.Query(getCrawlingJobName)
	if e != nil {
		return a, e
	}
	var aa CrawlingRunInfo
	aa.Status = "archived"
	for rows.Next() {
		if e = rows.Scan(&aa.ID, &aa.Crawler, &aa.Browsing, &aa.Starting, &aa.Ending); e != nil {
			return a, e
		}
		a = append(a, aa)
	}
	return a, e
}

// GetCrawlingRunDetail ...
func (d *DownloadDB) GetCrawlingRunDetail(id int64) (c *CrawlingRunInfo, e error) {
	r := d.Db.QueryRow(getCrawlingJob, id)
	var b []byte
	if e = r.Scan(&b); e != nil {
		return nil, e
	}
	c = &CrawlingRunInfo{}
	return c, json.Unmarshal(b, c)
}

// RemoveCrawlingRun ...
func (d *DownloadDB) RemoveCrawlingRun(id int64) error {
	_, err := d.Db.Exec(deleteCrawlingJob, id)
	return err
}

// SaveCrawlingJob ...
func (d *DownloadDB) SaveCrawlingJob(job *CrawlingRunInfo) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	_, err = d.Db.Exec(addCrawlingRun, job.ID, job.Crawler, job.Browsing, job.Starting, job.Current, data)
	return err
}

// GetEntityName ...
func (d *DownloadDB) GetEntityName() (a []string, e error) {
	rows, err := d.Db.Query(getEntityName)
	if err != nil {
		return nil, err
	}
	var s string
	var r []string
	for rows.Next() {
		if err = rows.Scan(&s); err != nil {
			return nil, err
		}
		r = append(r, s)
	}
	return r, nil
}

// CreateEntity ...
func (d *DownloadDB) CreateEntity(name string) (Entity, error) {
	stmt, err := d.Db.Prepare(addEntity)
	defer stmt.Close()
	if err != nil {
		return Entity{}, err
	}
	if _, err = stmt.Exec("", name, "[]"); err == nil {
		return Entity{
			Records: []Record{},
			Name:    name,
		}, nil
	}
	return Entity{}, err
}

// GetEntity ...
func (d *DownloadDB) GetEntity(name string) (Entity, error) {
	rows, err := d.Db.Query(getEntity, name)
	defer rows.Close()
	if err == nil {
		if rows.Next() {
			var data []byte
			entity := Entity{}
			if err := rows.Scan(&entity.Name, &data); err != nil {
				return Entity{}, err
			}
			if err := json.Unmarshal(data, &entity.Records); err != nil {
				return Entity{}, err
			}
			return entity, nil
		}
		if _, err = d.CreateEntity(name); err == nil {
			return d.GetEntity(name)
		}
		return Entity{}, err
	}
	return Entity{}, errors.New("Unangdle shit")
}

// AddRecordEntity ...
func (d *DownloadDB) AddRecordEntity(name string, records []Record) (Entity, error) {
	var entity Entity
	var err error
	var data []byte
	var stmt *sql.Stmt
	if entity, err = d.GetEntity(name); err == nil {
		for _, rr := range records {
			index := -1
			for i, r := range entity.Records {
				if rr.Link == r.Link {
					index = i
					break
				}
			}
			if index > -1 {
				entity.Records[index] = rr
			} else {
				entity.Records = append(entity.Records, rr)
			}

		}
		if stmt, err = d.Db.Prepare(updateEntityData); err == nil {
			defer stmt.Close()
			if data, err = json.Marshal(entity.Records); err == nil {
				_, err = stmt.Exec(data, name)
			}
		}
	}
	return entity, err
}

// GetDownloadDB ...
func GetDownloadDB() (*DownloadDB, error) {
	if dbsettings.OverWrite {
		os.Remove(dbsettings.Path)
	}
	var db *sql.DB
	var err error
	println("Openning database " + dbsettings.Path)
	if db, err = sql.Open("sqlite3", dbsettings.Path); err == nil && db != nil {
		db.SetMaxOpenConns(1)
		if _, err = db.Exec(createTable); err != nil {
			return nil, err
		}
		return &DownloadDB{
			Db: db,
		}, nil
	}
	return nil, err
}
