package indexer

import (
	"database/sql"
	"encoding/json"
	"errors"

	// blablabla
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
	`
	addEntity = `INSERT INTO entity(type,name,data) values(?,?,?)`
	updateEntityData = `UPDATE entity SET data = ? WHERE name = ?`

	getEntityName = `SELECT name FROM entity;`
	getEntity     = `SELEcT name, data FROM entity WHERE name = ?`
)

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

// DownloadDB ...
type DownloadDB struct {
	Db *sql.DB
}

// GetEntityName ...
func (d *DownloadDB) GetEntityName() []string {
	return []string{}
}

// CreateEntity ...
func (d *DownloadDB) CreateEntity(name string) (Entity,error) {
	stmt, err := d.Db.Prepare(addEntity)
	defer stmt.Close()
	if err != nil {
		return Entity{}, err
	}
	if _ ,err = stmt.Exec("",name,"[]"); err == nil {
		return Entity{
			Records: []Record{},
			Name: name,
		}, nil
	}
	return Entity{}, err
}

// GetEntity ...
func (d *DownloadDB) GetEntity(name string) (Entity,error) {
	rows, err := d.Db.Query(getEntity,name)
	defer rows.Close()
	if err == nil {
		if(rows.Next()) {
			var data []byte
			entity := Entity{}
			if err := rows.Scan(&entity.Name, &data); err != nil {
				return Entity{}, err
			}
			if err := json.Unmarshal(data, &entity.Records); err != nil {
				return Entity{}, err
			}
			return entity,nil
		} else {
			if _, err := d.CreateEntity(name); err == nil {
				return d.GetEntity(name)
			} else {
				return Entity{}, err
			}
		}
	}
	return Entity{},errors.New("Unangdle shit")
}

// AddRecordEntity ...
func (d *DownloadDB) AddRecordEntity(name string, records []Record) (Entity,error) {
	var entity Entity
	var err error
	var data []byte
	var stmt *sql.Stmt
	if entity, err = d.GetEntity(name); err == nil {
		entity.Records = append(entity.Records, records...)
		if stmt, err = d.Db.Prepare(updateEntityData); err == nil {
			defer stmt.Close()
			if data, err = json.Marshal(entity.Records); err == nil {
				_, err = stmt.Exec(data, name)
			}
		}
	}
	return entity,err
}

// GetDownloadDB ...
func GetDownloadDB() (*DownloadDB, error) {
	if db, err := sql.Open("sqlite3", "./data.db"); err == nil && db != nil {
		if _, err = db.Exec(createTable); err != nil {
			return nil,err
		}
		return &DownloadDB{
			Db: db,
		}, nil
	} else {
		return nil, err
	}
}

