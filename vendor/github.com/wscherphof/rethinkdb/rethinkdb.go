package rethinkdb

import (
	r "gopkg.in/gorethink/gorethink.v3"
	"github.com/wscherphof/env"
	"strings"
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

func Connect(db, address string) (err error) {
	DB = db
	if Session, err = r.Connect(r.ConnectOpts{
		Address: address,
		InitialCap: env.GetInt("DB_POOL_INITIAL", 100),
		MaxOpen: env.GetInt("DB_POOL_MAX", 100),
	}); err != nil {
		return
	}
	if _, err = r.DBCreate(DB).RunWrite(Session); err != nil {
		if strings.HasPrefix(err.Error(), "gorethink: Database `"+DB+"` already exists") {
			err = nil
		}
	}
	return
}

func tableCreate(table string, opts ...r.TableCreateOpts) (r.WriteResponse, error) {
	opt := r.TableCreateOpts{}
	if len(opts) == 1 {
		opt = opts[0]
	}
	opt.Shards = env.GetInt("DB_SHARDS", 1)
	opt.Replicas = env.GetInt("DB_REPLICAS", 1)
	return r.DB(DB).TableCreate(table, opt).RunWrite(Session)
}

func TableCreate(table string) (r.WriteResponse, error) {
	return tableCreate(table)
}

func TableCreatePK(table, pk string) (r.WriteResponse, error) {
	return tableCreate(table, r.TableCreateOpts{
		PrimaryKey: pk,
	})
}

func IndexCreate(table, name string, opt_fields ...string) (resp r.WriteResponse, err error) {
	num := len(opt_fields)
	if num == 0 {
		resp, err = r.DB(DB).Table(table).IndexCreate(name).RunWrite(Session)
	} else {
		resp, err = r.DB(DB).Table(table).IndexCreateFunc(name, func (row r.Term) interface{} {
			values := make([]interface{}, num, num)
			for i, field := range opt_fields {
				values[i] = row.Field(field)
			}
			return values
		}).RunWrite(Session)
	}
	r.DB(DB).Table(table).IndexWait(name)
	return
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

func GetIndex(table, index string, result interface{}, values ...interface{}) error {
	cursor, err := r.DB(DB).Table(table).GetAllByIndex(index, values...).Run(Session)
	return one(cursor, err, result)
}

func CountIndex(table, index string, result interface{}, values ...interface{}) error {
	var keys interface{}
	if num := len(values); num == 1 {
		keys = values[0]
	} else {
		// Compound index. Though GetAllByIndex is defined to accept keys as
		/// ...interface{}, it really only works with an explicit []interface{}
		keys := make([]interface{}, num, num)
		for i, v := range values {
			keys[i] = v
		}
	}
	cursor, err := r.DB(DB).Table(table).GetAllByIndex(index, keys).Count().Run(Session)
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

func InsertUpdate(table string, record interface{}) (response r.WriteResponse, err error) {
	response, err, _ = insert(table, record, r.InsertOpts{
		Conflict: "update",
	})
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

func Skip(table, index string, n int, opt_direction ...string) Term {
	orderByOpts := r.OrderByOpts{
		Index:  index, // "asc"
	}
	if len(opt_direction) == 1 && opt_direction[0] == "desc" {
		orderByOpts.Index = r.Desc(index)
	}
	term := r.DB(DB).Table(table).OrderBy(orderByOpts).Skip(n)
	return Term{Term: &term}
}

func DeleteTerm(term Term) (r.WriteResponse, error) {
	return term.Delete().RunWrite(Session)
}
