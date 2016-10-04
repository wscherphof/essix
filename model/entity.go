package model

import (
	"errors"
	"github.com/wscherphof/essix/util"
	db "github.com/wscherphof/rethinkdb"
	"log"
	"time"
)

var (
	ErrDuplicatePrimaryKey = errors.New("ErrDuplicatePrimaryKey")
)

type Entity struct {
	ID       string `gorethink:"id,omitempty"`
	table    string `gorethink:"-"`
	Created  time.Time
	Modified time.Time
}

func (e *Entity) Init(table string, id ...string) {
	e.setTable(table)
	e.Created = time.Now()
	e.Modified = e.Created
	if len(id) == 1 {
		e.ID = id[0]
	}
}

func (e *Entity) Create(record interface{}) (err error, conflict bool) {
	if _, err, conflict = db.Insert(e.table, record); conflict {
		err = ErrDuplicatePrimaryKey
	}
	return
}

func (e *Entity) Read(result interface{}) (err error, found bool) {
	return db.Get(e.table, e.ID, result)
}

func (e *Entity) Update(record interface{}) (err error) {
	e.Modified = time.Now()
	_, err = db.InsertUpdate(e.table, record)
	return
}

func (e *Entity) Delete() (err error) {
	_, err = db.Delete(e.table, e.ID)
	return
}

var tables = make(map[string]bool)

func (e *Entity) setTable(table string) {
	if !tables[table] {
		if _, err := db.TableCreate(table); err == nil {
			log.Println("INFO: table created:", table)
		}
		tables[table] = true
	}
	e.table = table
}

func NewCode() string {
	return string(util.URLEncode(util.Random()))
}
