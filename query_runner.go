package appx

import (
	"appengine"
	"appengine/datastore"
)

type QueryRunner struct {
	Context appengine.Context
	Query   *datastore.Query
}

func (this QueryRunner) Count() (int, error) {
	return this.Query.Count(this.Context)
}

func (this QueryRunner) Results(slice interface{}) error {
	return this.PagesIterator().LoadNext(slice)
}

func (this QueryRunner) Result(e Entity) error {
	return this.ItemsIterator().LoadNext(e)
}

func (this QueryRunner) StartFrom(cursor string) QueryRunner {
	c, _ := datastore.DecodeCursor(cursor)
	this.Query = this.Query.Start(c)
	return this
}

func (this QueryRunner) ItemsIterator() Iterator {
	return NewItemsIterator(this.Query, this.Context)
}

func (this QueryRunner) PagesIterator() Iterator {
	return NewPagesIterator(this.Query, this.Context)
}

