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
	Session                = db.Session
	DB                     = db.DB
	ErrEmptyResult         = db.ErrEmptyResult
	ErrDuplicatePrimaryKey = errors.New("ErrDuplicatePrimaryKey")
	typeReplacer           = strings.NewReplacer("*", "", ".", "_")
)

func Token() string {
	return string(util.URLEncode(util.Random()))
}

type Cursor struct {
	*db.Cursor
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
	if _, err, conflict = db.Insert(tbl(record), record); conflict {
		err = ErrDuplicatePrimaryKey
	}
	return
}

func (b *Base) Read(result interface{}) (err error) {
	if err = db.Get(tbl(result), b.ID, result); err == db.ErrEmptyResult {
		err = ErrEmptyResult
	}
	return
}

func ReadAll(record interface{}) (*Cursor, error) {
	c, e := db.All(tbl(record))
	return &Cursor{Cursor: c}, e
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

// TODO: check if this type can be private
type Term struct {
	term db.Term
}

type indexType struct {
	table  string
	column string
}

func Index(record interface{}, column string) *indexType {
	return &indexType{tbl(record), column}
}

func (i *indexType) Between(low, high interface{}, includeLeft, includeRight bool) *Term {
	return &Term{db.Between(i.table, i.column, low, high, includeLeft, includeRight)}
}

func (t *Term) Delete() (err error) {
	_, err = db.DeleteTerm(t.term)
	return
}
