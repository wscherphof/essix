package db

import (
	r "github.com/dancannon/gorethink"
	"github.com/wscherphof/expeertise/env"
	"log"
)

var (
	s      *r.Session
	dbname string
)

func init() {
	address := env.Get("DB_HOST") + env.Get("DB_PORT")
	dbname = env.Get("DB_NAME")
	if session, err := r.Connect(r.ConnectOpts{
		Address:  address,
		Database: dbname,
	}); err != nil {
		log.Fatalln("ERROR:", err)
	} else {
		s = session
		log.Println("INFO: DB connected to", dbname, "@", address)
	}
}

func insert(table string, record interface{}, opts ...r.InsertOpts) (r.WriteResponse, error) {
	return r.Table(table).Insert(record, opts...).RunWrite(s)
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
//   return r.Literal(args)
// }

// func Update (table, key string, arg interface{}) (r.WriteResponse, error) {
//   return r.Table(table).Get(key).Update(arg).RunWrite(s)
// }

func Delete(table, key string) (r.WriteResponse, error) {
	return r.Table(table).Get(key).Delete().RunWrite(s)
}

func Truncate(table string) (r.WriteResponse, error) {
	return r.Table(table).Delete().RunWrite(s)
}

func tableCreate(table string, opts ...r.TableCreateOpts) (*r.Cursor, error) {
	return r.Db(dbname).TableCreate(table, opts...).Run(s)
}

func TableCreate(table string) (*r.Cursor, error) {
	return tableCreate(table)
}

func TableCreatePK(table, pk string) (*r.Cursor, error) {
	return tableCreate(table, r.TableCreateOpts{
		PrimaryKey: pk,
	})
}

func IndexCreate(table, field string) (*r.Cursor, error) {
	return r.Table(table).IndexCreate(field).Run(s)
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
	return r.Table(table).Between(low, high, optArgs)
}

func DeleteTerm(term r.Term) (r.WriteResponse, error) {
	return term.Delete().RunWrite(s)
}

func Get(table, key string, result interface{}) (err error, found bool) {
	if cursor, e := r.Table(table).Get(key).Run(s); e != nil {
		err = e
	} else if e = cursor.One(result); e == nil {
		found = true
	} else if e != r.ErrEmptyResult {
		err = e
	}
	return
}

func One(table string, result interface{}) (err error, found bool) {
	if cursor, e := r.Table(table).Run(s); e != nil {
		err = e
	} else if e = cursor.One(result); e == nil {
		found = true
	} else if e != r.ErrEmptyResult {
		err = e
	}
	return
}

func All(table string) (cursor *r.Cursor, err error) {
	return r.Table(table).Run(s)
}
