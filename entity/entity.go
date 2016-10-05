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

func table(record interface{}) string {
	t := fmt.Sprintf("%T", record)
	return tableNameReplacer.Replace(t)
}

func Register(record interface{}) {
	if _, err := db.TableCreate(table(record)); err == nil {
		log.Println("INFO: table created:", table)
	}
}

func Token() string {
	return string(util.URLEncode(util.Random()))
}

type Base struct {
	ID       string `gorethink:"id,omitempty"`
	table    string `gorethink:"-"`
	Created  time.Time
	Modified time.Time
}

func (b *Base) Create(record interface{}) (err error, conflict bool) {
	b.Created = time.Now()
	b.Modified = b.Created
	if _, err, conflict = db.Insert(table(record), record); conflict {
		err = ErrDuplicatePrimaryKey
	}
	return
}

func (b *Base) Read(result interface{}) (err error, found bool) {
	return db.Get(table(result), b.ID, result)
}

func (b *Base) Update(record interface{}) (err error) {
	if b.Created.IsZero() {
		b.Created = time.Now()
	}
	b.Modified = time.Now()
	_, err = db.InsertUpdate(table(record), record)
	return
}

func (b *Base) Delete(record interface{}) (err error) {
	_, err = db.Delete(table(record), b.ID)
	return
}
