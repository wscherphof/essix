package db

import (
	r "gopkg.in/dancannon/gorethink.v2"
	"github.com/wscherphof/essix/env"
	"log"
)

var (
	s      *r.Session
	dbname string
)

func init() {
	dbname = env.Get("DB_NAME", "essix")
	address := env.Get("DB_HOST", "db1") +":"+ env.Get("DB_PORT", "28015")
	if session, err := r.Connect(r.ConnectOpts{Address: address}); err != nil {
		log.Fatalln("ERROR:", err)
	} else {
		if _, err := r.DBCreate(dbname).RunWrite(session); err == nil {
			log.Println("INFO: created DB", dbname, "@", address)
		}
		s = session
		log.Println("INFO: connected to DB", dbname, "@", address)
	}
}

func insert(table string, record interface{}, opts ...r.InsertOpts) (r.WriteResponse, error) {
	return r.DB(dbname).Table(table).Insert(record, opts...).RunWrite(s)
}

func Insert(table string, record interface{}) (r.WriteResponse, error) {
	return insert(table, record)
}

func InsertUpdate(table string, record interface{}) (r.WriteResponse, error) {
	return insert(table, record, r.InsertOpts{
		Conflict: "update",
	})
}

// Unused, untested:

// func Literal (args ...interface{}) r.Term {
//   return r.DB(dbname).Literal(args)
// }

// func Update (table, key string, arg interface{}) (r.WriteResponse, error) {
//   return r.DB(dbname).Table(table).Get(key).Update(arg).RunWrite(s)
// }

func Delete(table, key string) (r.WriteResponse, error) {
	return r.DB(dbname).Table(table).Get(key).Delete().RunWrite(s)
}

func Truncate(table string) (r.WriteResponse, error) {
	return r.DB(dbname).Table(table).Delete().RunWrite(s)
}

func tableCreate(table string, opts ...r.TableCreateOpts) (r.WriteResponse, error) {
	return r.DB(dbname).TableCreate(table, opts...).RunWrite(s)
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
	return r.DB(dbname).Table(table).IndexCreate(field).RunWrite(s)
}

func Between(table, index string, low, high interface{}, includeLeft, includeRight bool) r.Term {
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
	return r.DB(dbname).Table(table).Between(low, high, optArgs)
}

func DeleteTerm(term r.Term) (r.WriteResponse, error) {
	return term.Delete().RunWrite(s)
}

func Get(table, key string, result interface{}) (err error, found bool) {
	cursor, e := r.DB(dbname).Table(table).Get(key).Run(s)
	if cursor != nil {
		defer cursor.Close()
	}
	if e != nil {
		err = e
	} else if e = cursor.One(result); e == nil {
		found = true
	} else if e != r.ErrEmptyResult {
		err = e
	}
	return
}

func One(table string, result interface{}) (err error, found bool) {
	cursor, e := r.DB(dbname).Table(table).Run(s)
	if cursor != nil {
		defer cursor.Close()
	}
	if e != nil {
		err = e
	} else if e = cursor.One(result); e == nil {
		found = true
	} else if e != r.ErrEmptyResult {
		err = e
	}
	return
}

func All(table string) (cursor *r.Cursor, err error) {
	return r.DB(dbname).Table(table).Run(s)
}
