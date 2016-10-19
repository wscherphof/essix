package entity

import (
	"errors"
	"fmt"
	"github.com/wscherphof/env"
	db "github.com/wscherphof/rethinkdb"
	"log"
	"strings"
	"time"
)

var (
	Session                = db.Session
	DB                     = db.DB
	ErrEmptyResult         = db.ErrEmptyResult
	ErrDuplicatePrimaryKey = errors.New("ErrDuplicatePrimaryKey")
	typeReplacer           = strings.NewReplacer("*", "", ".", "_")
)

func init() {
	name := env.Get("DB_NAME")
	address := env.Get("DB_ADDRESS")
	if err := db.Connect(name, address); err != nil {
		log.Fatalf("ERROR: connecting to DB %s@%s failed. %T %s", name, address, err, err)
	} else {
		log.Printf("INFO: connected to DB %s@%s", name, address)
	}
}

type tableType struct {
	name string
	new  bool
}

var tables = make(map[string]string, 100)

func getType(record interface{}) string {
	tpe := fmt.Sprintf("%T", record)
	return typeReplacer.Replace(tpe)
}

func Register(record interface{}, table ...string) (ret *tableType) {
	tpe := getType(record)
	tbl := tpe
	if len(table) == 1 {
		tbl = table[0]
	}
	tables[tpe] = tbl
	ret = &tableType{name: tbl}
	if _, err := db.TableCreate(tbl); err == nil {
		ret.new = true
		log.Println("INFO: table created:", tbl)
	}
	return
}

func (t *tableType) Index(column string) *tableType {
	if _, err := db.IndexCreate(t.name, column); err != nil {
		if t.new {
			log.Println("ERROR: failed to create index:", t.name, column, err)
		}
	}
	return t
}

type Base struct {
	ID       string `gorethink:"id,omitempty"`
	Created  time.Time
	Modified time.Time
}

func tbl(record interface{}) string {
	tpe := getType(record)
	return tables[tpe]
}

func (b *Base) Create(record interface{}) (err error, conflict bool) {
	b.Created = time.Now()
	b.Modified = b.Created
	if response, e, c := db.Insert(tbl(record), record); c {
		err, conflict = ErrDuplicatePrimaryKey, true
	} else if e != nil {
		err = e
	} else if b.ID == "" {
		b.ID = response.GeneratedKeys[0]
	}
	return
}

func (b *Base) Read(result interface{}) (err error, empty bool) {
	if err = db.Get(tbl(result), b.ID, result); err == db.ErrEmptyResult {
		err, empty = ErrEmptyResult, true
	}
	return
}

func ReadAll(record interface{}) (*db.Cursor, error) {
	return db.All(tbl(record))
}

func (b *Base) Update(record interface{}) (err error) {
	if b.Created.IsZero() {
		b.Created = time.Now()
	}
	b.Modified = time.Now()
	if response, e := db.InsertUpdate(tbl(record), record); e != nil {
		err = e
	} else if b.ID == "" {
		b.ID = response.GeneratedKeys[0]
	}
	return
}

func (b *Base) Delete(record interface{}) (err error) {
	_, err = db.Delete(tbl(record), b.ID)
	return
}

type indexType struct {
	table  string
	column string
}

func Index(record interface{}, column string) *indexType {
	return &indexType{tbl(record), column}
}

func (i *indexType) Between(low interface{}, includeLow bool, high interface{}, includeHigh bool) *term {
	return &term{db.Between(i.table, i.column, low, includeLow, high, includeHigh)}
}

func (i *indexType) Read(value, result interface{}) (err error, empty bool) {
	if err = db.GetIndex(i.table, i.column, value, result); err == db.ErrEmptyResult {
		err, empty = ErrEmptyResult, true
	}
	return
}

type term struct {
	term db.Term
}

func (t *term) Delete() (int, error) {
	resp, e := db.DeleteTerm(t.term)
	return resp.Deleted, e
}
