package rethinkdb

import (
	"github.com/wscherphof/env"
	r "gopkg.in/dancannon/gorethink.v2"
	"log"
)

var (
	Session        *r.Session
	DB             string
	ErrEmptyResult = r.ErrEmptyResult
)

type Cursor struct {
	*r.Cursor
}

type Term struct {
	*r.Term
}

func init() {
	DB = env.Get("DB_NAME", "essix")
	address := env.Get("DB_ADDRESS", "db1")
	var err error
	if Session, err = r.Connect(r.ConnectOpts{Address: address}); err != nil {
		log.Fatalln("ERROR:", err)
	}
	if _, err := r.DBCreate(DB).RunWrite(Session); err == nil {
		log.Println("INFO: created DB", DB, "@", address)
	} else {
		log.Println("INFO: connected to DB", DB, "@", address)
	}
}

func tableCreate(table string, opts ...r.TableCreateOpts) (r.WriteResponse, error) {
	return r.DB(DB).TableCreate(table, opts...).RunWrite(Session)
}

func TableCreate(table string) (r.WriteResponse, error) {
	return tableCreate(table)
}

func TableCreatePK(table, pk string) (r.WriteResponse, error) {
	return tableCreate(table, r.TableCreateOpts{
		PrimaryKey: pk,
	})
}

func IndexCreate(table, field string) (r.WriteResponse, error) {
	return r.DB(DB).Table(table).IndexCreate(field).RunWrite(Session)
}

func insert(table string, record interface{}, opts ...r.InsertOpts) (response r.WriteResponse, err error, conflict bool) {
	response, err = r.DB(DB).Table(table).Insert(record, opts...).RunWrite(Session)
	conflict = r.IsConflictErr(err)
	return
}

func Insert(table string, record interface{}) (response r.WriteResponse, err error, conflict bool) {
	return insert(table, record)
}

func Get(table, key string, result interface{}) (err error) {
	cursor, e := r.DB(DB).Table(table).Get(key).Run(Session)
	if cursor != nil {
		defer cursor.Close()
	}
	if e != nil {
		err = e
	} else if err = cursor.One(result); err == r.ErrEmptyResult {
		err = ErrEmptyResult
	}
	return
}

func One(table string, result interface{}) (err error) {
	cursor, e := r.DB(DB).Table(table).Run(Session)
	if cursor != nil {
		defer cursor.Close()
	}
	if e != nil {
		err = e
	} else if err = cursor.One(result); err == r.ErrEmptyResult {
		err = ErrEmptyResult
	}
	return
}

func All(table string) (*Cursor, error) {
	c, e := r.DB(DB).Table(table).Run(Session)
	return &Cursor{Cursor: c}, e
}

func InsertUpdate(table string, record interface{}, id ...string) (response r.WriteResponse, err error) {
	if len(id) == 1 {
		response, err = r.DB(DB).Table(table).Get(id[0]).Update(record).RunWrite(Session)
	} else {
		response, err, _ = insert(table, record, r.InsertOpts{
			Conflict: "update",
		})
	}
	return
}
func Delete(table, key string) (r.WriteResponse, error) {
	return r.DB(DB).Table(table).Get(key).Delete().RunWrite(Session)
}

func Truncate(table string) (r.WriteResponse, error) {
	return r.DB(DB).Table(table).Delete().RunWrite(Session)
}

func Between(table, index string, low, high interface{}, includeLeft, includeRight bool) Term {
	optArgs := r.BetweenOpts{
		LeftBound:  "closed",
		RightBound: "closed",
	}
	if len(index) > 0 {
		optArgs.Index = index
	}
	if includeLeft {
		optArgs.LeftBound = "open"
	}
	if includeRight {
		optArgs.RightBound = "open"
	}
	if low == nil {
		low = r.MinVal
	}
	if high == nil {
		low = r.MaxVal
	}
	term := r.DB(DB).Table(table).Between(low, high, optArgs)
	return Term{Term: &term}
}

func DeleteTerm(term Term) (r.WriteResponse, error) {
	return term.Delete().RunWrite(Session)
}
