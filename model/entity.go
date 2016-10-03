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

func (e *Entity) New(record interface{}) (err error, conflict bool) {
	e.Created = time.Now()
	e.Modified = e.Created
	if _, err, conflict = db.Insert(e.table, e); err == nil {
		_, err = db.InsertUpdate(e.table, record, e.ID)
	} else if conflict {
		err = ErrDuplicatePrimaryKey
	}
	return
}

func (e *Entity) Save(record interface{}) (err error) {
	e.Modified = time.Now()
	if _, err = db.InsertUpdate(e.table, e); err == nil {
		_, err = db.InsertUpdate(e.table, record, e.ID)
	}
	return
}

func Table(name string) {
	if _, err := db.TableCreate(name); err == nil {
		log.Println("INFO: table created:", name)
	}
}

func Code() string {
	return string(util.URLEncode(util.Random()))
}
