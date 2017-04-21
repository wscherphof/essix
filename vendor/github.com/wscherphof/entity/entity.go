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

type entity interface {
	table(record interface{}) string
}

/*
Register manages which database table to use for which type.
The table to use is derived from the type name of the entity argument. This can
be overruled by setting the Table field of the entity's Base field.
If the table to use is not present in the database, it creates the table.
If creating the table fails, it logs the error, and panics.

Call Index() on the result to ensure a secondary index for a field.
Call Index() on the result of Index() for another index on the same table.
*/
func Register(e entity) (ret *TableType) {
	table := e.table(e)
	ret = &TableType{name: table}
	if _, err := db.TableCreate(table); err != nil {
		if !strings.HasPrefix(err.Error(), "gorethink: Table `"+DB+"."+table+"` already exists") {
			log.Panicln("ERROR: failed to create table:", table, err)
		}
	} else {
		ret.new = true
		log.Println("INFO: table created:", table)
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

Later, call Base.Index() for an IndexType value:
	bus := initBus()
	busFooIndex := bus.Index(bus, "Foo")
	busBarIndex := bus.Index(bus, "Bar")
*/
func (t *TableType) Index(name string, column ...string) *TableType {
	response, err := db.IndexCreate(t.name, name, column...)
	if t.new && response.Created != 1 {
		log.Println("ERROR: failed to create index:", t.name, name, column, err)
	}
	return t
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

	// Table overrides the table to use for this business object. Leave empty to
	// have it derived from the object's type name.
	Table string `gorethink:",omitempty"`
}

var typeReplacer = strings.NewReplacer(".", "_", "[", "", "]", "", "*", "", "0", "", "1", "", "2", "", "3", "", "4", "", "5", "", "6", "", "7", "", "8", "", "9", "")

func (b *Base) table(record interface{}) string {
	if b.Table == "" {
		T := fmt.Sprintf("%T", record)
		b.Table = typeReplacer.Replace(T)
	}
	return b.Table
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
	if response, e, c := db.Insert(b.table(record), record); c {
		err, conflict = ErrDuplicatePrimaryKey, true
	} else if e != nil {
		err = e
	} else if len(response.GeneratedKeys) == 1 {
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
	if err = db.Get(b.table(result), b.ID, result); err == db.ErrEmptyResult {
		err, empty = ErrEmptyResult, true
	}
	return
}

/*
ReadAll returns a
[gorethink *Cursor](https://godoc.org/github.com/dancannon/gorethink#Cursor),
holding all records of the given type.
	bus := initBus()
	if cursor, err := bus.ReadAll(bus); err == nil {
		defer cursor.Close()
		for cursor.Next(bus) {
			// doSomethingWith(bus)
		}
		if cursor.Err() != nil {
			panic(cursor.Err())
		}
	}
*/
func (b *Base) ReadAll(record interface{}) (*db.Cursor, error) {
	return db.All(b.table(record))
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
	if response, e := db.InsertUpdate(b.table(record), record); e != nil {
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
	_, err = db.Delete(b.table(record), b.ID)
	return
}

/*
An IndexType provides index operations.
*/
type IndexType struct {
	table  string
	name string
}

/*
Index returns the IndexType for the given record's table & column. Create
indexes by calling Index() on the result of Register()
	entity.Register(&Bus{}).Index("Foo")
	...
	bus := initBus()
	busFooIndex := bus.Index(bus, "Foo")
*/
func (b *Base) Index(record interface{}, name string) *IndexType {
	return &IndexType{b.table(record), name}
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
func (i *IndexType) Read(result interface{}, values ...interface{}) (err error, empty bool) {
	if err = db.GetIndex(i.table, i.name, result, values...); err == db.ErrEmptyResult {
		err, empty = ErrEmptyResult, true
	}
	return
}

/*
Count counts the records from the database, with the given value for this
index's column.
It sets empty to true, if the error is that a record with that value doesn't
exist.
	var num int
	if err := busFooIndex.Count(1, &num); err == nil {
		// doSomethingWith(num)
	}
*/
func (i *IndexType) Count(result interface{}, values ...interface{}) error {
	return db.CountIndex(i.table, i.name, result, values...)
}

/*
Between filters the index's table on low & high limits for the value of the
index's column. It returns a Term holding the selected records.
	fooOne := busFooIndex.Between(0, false, 2, false)
Pass nil for low to use its minimum value. Pass nil for high to use its maximum
value.
*/
func (i *IndexType) Between(low interface{}, includeLow bool, high interface{}, includeHigh bool) Term {
	return Term(db.Between(i.table, i.name, low, includeLow, high, includeHigh))
}

/*
Skip cuts off the table's first (or last) n records, when ordered by this index.
Direction "asc" (default) cuts off the first n records;
Direction "desc" cuts off the last n records.
*/
func (i *IndexType) Skip(n int, opt_direction ...string) Term {
	return Term(db.Skip(i.table, i.name, n, opt_direction...))
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
