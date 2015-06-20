package appx

import (
	"appengine"
	"appengine/datastore"
)

type QueryRunner struct {
	context appengine.Context
	query   *datastore.Query
}

func NewQueryRunner(c appengine.Context, q *datastore.Query) *QueryRunner {
	return &QueryRunner{c, q}
}

func (this QueryRunner) Count() (int, error) {
	return this.query.Count(this.context)
}

func (this QueryRunner) Results(slice interface{}) error {
	return this.PagesIterator().LoadNext(slice)
}

func (this QueryRunner) Result(e Entity) error {
	return this.ItemsIterator().LoadNext(e)
}

func (this QueryRunner) StartFrom(cursor string) QueryRunner {
	c, _ := datastore.DecodeCursor(cursor)
	this.query = this.query.Start(c)
	return this
}

func (this QueryRunner) ItemsIterator() Iterator {
	return NewItemsIterator(this.query, this.context)
}

func (this QueryRunner) PagesIterator() Iterator {
	return NewPagesIterator(this.query, this.context)
}

