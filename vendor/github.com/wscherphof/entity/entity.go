/*
Package entity provides CRUD operations for business objects, managing their
storage in a RethinkDB cluster.

It connects to the database on init, requiring two environment variables:
	DB_NAME
	DB_ADDRESS
It logs a fatal error if connecting fails.

[Embed](https://golang.org/doc/effective_go.html#embedding) the Base type in
a business object type:
	type Bus struct {
		*entity.Base
		Foo int
		Bar string
	}

Register the business object before using its entity methods:
	entity.Register(&Bus{})

Do initialise the embbedded Base:
	func initBus(opt_id string) *Bus {
		base := &entity.Base{}
		if len(opt_id) == 1 {
			base.ID = opt_id[0]
		}
		return &Bus{Base: base}
	}

Business object ready to use:
	func NewBus() (bus *Bus) {
		bus = initBus()
		if err := bus.Update(bus); err != nil {
			panic(err)
		}
		return
	}
*/
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

	// Session provides low level access to the RethinkDB session.
	Session = db.Session

	// DB provides low level access to the RethinkDB database.
	DB = db.DB

	// ErrEmptyResult is returned from read operations if no record was found.
	ErrEmptyResult = db.ErrEmptyResult

	// ErrDuplcatePrimaryKey is retuned from Create if a record with the same ID
	// already exists.
	ErrDuplicatePrimaryKey = errors.New("ErrDuplicatePrimaryKey")
)

func init() {
	name := env.Get("DB_NAME")
	address := env.Get("DB_ADDRESS")
	if err := db.Connect(name, address); err != nil {
		log.Fatalf("ERROR: connecting to DB %s@%s failed. %T %s", name, address, err, err)
	} else {
		DB, Session = db.DB, db.Session
		log.Printf("INFO: connected to DB %s@%s", name, address)
	}
}

var tables = make(map[string]string, 100)

/*
Register manages which database table to use for which type.
If the table argument is omitted, it derives the table name from the type name
of the record argument.
If the table to use is not present in the database, it creates the table.
If creating the table fails, it logs the error, and panics.

Call Index() on the result to ensure a secondary index for a field.
Call Index() on the result of Index() for another index on the same table.
*/
func Register(record interface{}, table ...string) (ret *TableType) {
	tpe := getType(record)
	tbl := tpe
	if len(table) == 1 {
		tbl = table[0]
	}
	tables[tpe] = tbl
	ret = &TableType{name: tbl}
	if _, err := db.TableCreate(tbl); err != nil {
		if !strings.HasPrefix(err.Error(), "gorethink: Table `"+DB+"."+tbl+"` already exists") {
			log.Panicln("ERROR: failed to create table:", tbl, err)
		}
	} else {
		ret.new = true
		log.Println("INFO: table created:", tbl)
	}
	return
}

/*
A TableType is the result of Register()
*/
type TableType struct {
	name string
	new  bool
}

/*
Index ensures a secondary database index on the given column.
	entity.Register(&Bus{}).Index("Foo")
or
	entity.Register(&Bus{}).Index("Foo").Index("Bar")

Later, call entity.Index() for an IndexType value:
	busFooIndex := entity.Index(&bus{}, "Foo")
	busBarIndex := entity.Index(&bus{}, "Bar")
*/
func (t *TableType) Index(column string) *TableType {
	if _, err := db.IndexCreate(t.name, column); err != nil {
		if t.new {
			log.Println("ERROR: failed to create index:", t.name, column, err)
		}
	}
	return t
}

var typeReplacer = strings.NewReplacer("*", "", ".", "_")

func getType(record interface{}) string {
	tpe := fmt.Sprintf("%T", record)
	return typeReplacer.Replace(tpe)
}

/*
Base is the base type to embed in business objects.
*/
type Base struct {

	// ID is the database record ID. If empty, the database generates a unique value
	// for it.
	ID string `gorethink:"id,omitempty"`

	// Created holds the time when the record was first created in the database.
	Created time.Time

	// Modified holds the time when the record was last modified in the database.
	Modified time.Time
}

func tbl(record interface{}) string {
	tpe := getType(record)
	return tables[tpe]
}

/*
Create saves the record in the database.
If a record with the same ID already exists, conflict is set to true.
	bus := initBus()
	err, conflict := bus.Create(bus)
*/
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

/*
Read loads a record from the database, using its ID.
It sets empty to true, if the error is that a record with that ID doesn't exist.
	bus := initBus("1")
	err, empty := bus.Read(bus)
*/
func (b *Base) Read(result interface{}) (err error, empty bool) {
	if err = db.Get(tbl(result), b.ID, result); err == db.ErrEmptyResult {
		err, empty = ErrEmptyResult, true
	}
	return
}

/*
ReadAll returns a
[gorethink *Cursor](https://godoc.org/github.com/dancannon/gorethink#Cursor),
holding all records of the given type.
	if cursor, err := entity.ReadAll(&Bus{}); err == nil {
		defer cursor.Close()
		bus := initBus()
		for cursor.Next(bus) {
			// doSomethingWith(bus)
		}
		if cursor.Err() != nil {
			panic(cursor.Err())
		}
	}
*/
func ReadAll(record interface{}) (*db.Cursor, error) {
	return db.All(tbl(record))
}

/*
Update saves the current state of the record in the database.
It creates a record with that ID, if it doesn't exist.
	bus := initBus("1")
	bus.Foo = 1
	err := bus.Update(bus)
*/
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

/*
Delete deletes a record from the database, using its ID.
	bus := initBus("1")
	err := bus.Delete(bus)
*/
func (b *Base) Delete(record interface{}) (err error) {
	_, err = db.Delete(tbl(record), b.ID)
	return
}

/*
An IndexType provides index operations.
*/
type IndexType struct {
	table  string
	column string
}

/*
Index returns the IndexType for the given record's table & column. Create
indexes by calling Index() on the result of Register()
	entity.Register(&Bus{}).Index("Foo")
	...
	busFooIndex := entity.Index(&bus{}, "Foo")
*/
func Index(record interface{}, column string) *IndexType {
	return &IndexType{tbl(record), column}
}

/*
Read loads a record from the database, with the given value for this index's
column.
It sets empty to true, if the error is that a record with that value doesn't
exist.
	bus := initBus()
	if err, empty := busFooIndex.Read(1, bus); err == nil {
		// doSomethingWith(bus)
	}
*/
func (i *IndexType) Read(value, result interface{}) (err error, empty bool) {
	if err = db.GetIndex(i.table, i.column, value, result); err == db.ErrEmptyResult {
		err, empty = ErrEmptyResult, true
	}
	return
}

/*
Between filters the index's table on low & high limits for the value of the
index's column. It returns a Term holding the selected records.
	fooOne := busFooIndex.Between(0, false, 2, false)
Pass nil for low to use its minimum value. Pass nil for high to use its maximum
value.
*/
func (i *IndexType) Between(low interface{}, includeLow bool, high interface{}, includeHigh bool) Term {
	return Term(db.Between(i.table, i.column, low, includeLow, high, includeHigh))
}

/*
A Term holds the database records selected by filter operations.
*/
type Term db.Term

/*
Delete deletes all the Term's records from the database. It returns the number
of records deleted, and an error if one occurred.
	fooOne := busFooIndex.Between(0, false, 2, false)
	deleted, err := fooOne.Delete()
*/
func (t Term) Delete() (int, error) {
	resp, e := db.DeleteTerm(db.Term(t))
	return resp.Deleted, e
}
