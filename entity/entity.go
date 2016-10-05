package entity

import (
	"errors"
	"fmt"
	"github.com/wscherphof/essix/util"
	db "github.com/wscherphof/rethinkdb"
	"log"
	"strings"
	"time"
)

var (
	ErrDuplicatePrimaryKey = errors.New("ErrDuplicatePrimaryKey")
	tableNameReplacer      = strings.NewReplacer("*", "", ".", "_")
)

func Token() string {
	return string(util.URLEncode(util.Random()))
}

type Base struct {
	ID       string `gorethink:"id,omitempty"`
	Table    string `gorethink:"-"`
	Created  time.Time
	Modified time.Time
}

func (b *Base) getTable(record interface{}) (table string) {
	table = b.Table
	if table == "" {
		T := fmt.Sprintf("%T", record)
		table = tableNameReplacer.Replace(T)
	}
	return
}

func (b *Base) Register(record interface{}) {
	table := b.getTable(record)
	if _, err := db.TableCreate(table); err == nil {
		log.Println("INFO: table created:", table)
	}
}

func (b *Base) Create(record interface{}) (err error, conflict bool) {
	b.Created = time.Now()
	b.Modified = b.Created
	table := b.getTable(record)
	if _, err, conflict = db.Insert(table, record); conflict {
		err = ErrDuplicatePrimaryKey
	}
	return
}

func (b *Base) Read(result interface{}) (err error, found bool) {
	table := b.getTable(result)
	return db.Get(table, b.ID, result)
}

func (b *Base) Update(record interface{}) (err error) {
	if b.Created.IsZero() {
		b.Created = time.Now()
	}
	b.Modified = time.Now()
	table := b.getTable(record)
	_, err = db.InsertUpdate(table, record)
	return
}

func (b *Base) Delete(record interface{}) (err error) {
	table := b.getTable(record)
	_, err = db.Delete(table, b.ID)
	return
}
