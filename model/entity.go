package model

import (
	"errors"
	"github.com/wscherphof/essix/util"
	db "github.com/wscherphof/rethinkdb"
	"log"
	"time"
	"fmt"
)

var (
	ErrDuplicatePrimaryKey = errors.New("ErrDuplicatePrimaryKey")
	ErrUnregisteredType    = errors.New("ErrUnregisteredType")
)

type Entity struct {
	ID       string `gorethink:"id,omitempty"`
	table    string `gorethink:"-"`
	Created  time.Time
	Modified time.Time
}

var tables = make(map[string]string)

func getType(record interface{}) (t string) {
	t = fmt.Sprintf("%T", record)
	return
}

func Register(record interface{}, table string) {
	t := getType(record)
	if tables[t] == "" {
		if _, err := db.TableCreate(table); err == nil {
			log.Println("INFO: table created:", table)
		}
	}
	tables[t] = table
}

func getTable(record interface{}) (value string, err error) {
	if val, ok := tables[getType(record)]; ok {
		value = val
	} else {
		err = ErrUnregisteredType
	}
	return
}

func (n *Entity) Init(id ...string) {
	n.Created = time.Now()
	n.Modified = n.Created
	if len(id) == 1 {
		n.ID = id[0]
	}
}

func (n *Entity) Create(record interface{}) (err error, conflict bool) {
	if table, e := getTable(record); e != nil {
		err = e
	} else if _, err, conflict = db.Insert(table, record); conflict {
		err = ErrDuplicatePrimaryKey
	}
	return
}

func (n *Entity) Read(result interface{}) (err error, found bool) {
	if table, e := getTable(result); e != nil {
		err = e
	} else {
		err, found = db.Get(table, n.ID, result)
	}
	return
}

func (n *Entity) Update(record interface{}) (err error) {
	if table, e := getTable(record); e != nil {
		err = e
	} else {
		n.Modified = time.Now()
		_, err = db.InsertUpdate(table, record)
	}
	return
}

func (n *Entity) Delete(record interface{}) (err error) {
	if table, e := getTable(record); e != nil {
		err = e
	} else {
		_, err = db.Delete(table, n.ID)
	}
	return
}

func NewCode() string {
	return string(util.URLEncode(util.Random()))
}
