package rethinkdb

import (
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

func Connect(db, address string) {
	DB = db
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

func one(cursor *r.Cursor, err error, result interface{}) error {
	if err != nil {
		return err
	}
	defer cursor.Close()
	if cursor.IsNil() {
		return ErrEmptyResult
	}
	return cursor.One(result)
}

func Get(table, key string, result interface{}) error {
	cursor, err := r.DB(DB).Table(table).Get(key).Run(Session)
	return one(cursor, err, result)
}

func GetIndex(table, index string, value, result interface{}) error {
	cursor, err := r.DB(DB).Table(table).GetAllByIndex(index, value).Run(Session)
	return one(cursor, err, result)
}

func One(table string, result interface{}) error {
	cursor, err := r.DB(DB).Table(table).Run(Session)
	return one(cursor, err, result)
}

func All(table string) (*Cursor, error) {
	cursor, err := r.DB(DB).Table(table).Run(Session)
	return &Cursor{Cursor: cursor}, err
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

func Between(table, index string, low interface{}, includeLow bool, high interface{}, includeHigh bool) Term {
	optArgs := r.BetweenOpts{
		LeftBound:  "closed",
		RightBound: "closed",
	}
	if len(index) > 0 {
		optArgs.Index = index
	}
	if includeLow {
		optArgs.LeftBound = "open"
	}
	if includeHigh {
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
