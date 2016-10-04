package model

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
)

type Entity struct {
	ID       string `gorethink:"id,omitempty"`
	table    string `gorethink:"-"`
	Created  time.Time
	Modified time.Time
}

var r = strings.NewReplacer("*", "", ".", "_")

func table(record interface{}) string {
	t := fmt.Sprintf("%T", record)
	return r.Replace(t)
}

func Register(record interface{}) {
	if _, err := db.TableCreate(table(record)); err == nil {
		log.Println("INFO: table created:", table)
	}
}

func (n *Entity) Init(id ...string) {
	n.Created = time.Now()
	n.Modified = n.Created
	if len(id) == 1 {
		n.ID = id[0]
	}
}

func (n *Entity) Create(record interface{}) (err error, conflict bool) {
	if _, err, conflict = db.Insert(table(record), record); conflict {
		err = ErrDuplicatePrimaryKey
	}
	return
}

func (n *Entity) Read(result interface{}) (err error, found bool) {
	return db.Get(table(result), n.ID, result)
}

func (n *Entity) Update(record interface{}) (err error) {
	n.Modified = time.Now()
	_, err = db.InsertUpdate(table(record), record)
	return
}

func (n *Entity) Delete(record interface{}) (err error) {
	_, err = db.Delete(table(record), n.ID)
	return
}

func NewCode() string {
	return string(util.URLEncode(util.Random()))
}
