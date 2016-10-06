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
	typeReplacer           = strings.NewReplacer("*", "", ".", "_")
	tables                 = make(map[string]string, 100)
)

func getType(record interface{}) string {
	tpe := fmt.Sprintf("%T", record)
	return typeReplacer.Replace(tpe)
}

func Register(record interface{}, table ...string) {
	tpe := getType(record)
	tbl := tpe
	if len(table) == 1 {
		tbl = table[0]
	}
	tables[tpe] = tbl
	if _, err := db.TableCreate(tbl); err == nil {
		log.Println("INFO: table created:", tbl)
	}
}

func tbl(record interface{}) string {
	tpe := getType(record)
	return tables[tpe]
}

func Token() string {
	return string(util.URLEncode(util.Random()))
}

type Base struct {
	ID       string `gorethink:"id,omitempty"`
	Created  time.Time
	Modified time.Time
}

func (b *Base) Create(record interface{}) (err error, conflict bool) {
	b.Created = time.Now()
	b.Modified = b.Created
	if _, err, conflict = db.Insert(tbl(record), record); conflict {
		err = ErrDuplicatePrimaryKey
	}
	return
}

func (b *Base) Read(result interface{}) (err error, found bool) {
	return db.Get(tbl(result), b.ID, result)
}

func (b *Base) Update(record interface{}) (err error) {
	if b.Created.IsZero() {
		b.Created = time.Now()
	}
	b.Modified = time.Now()
	_, err = db.InsertUpdate(tbl(record), record)
	return
}

func (b *Base) Delete(record interface{}) (err error) {
	_, err = db.Delete(tbl(record), b.ID)
	return
}
