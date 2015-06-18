package ds

import (
	"appengine"
	"appengine/datastore"
)

type DatastoreIterator struct {
	query              *datastore.Query
	context            appengine.Context
	nextCursor         datastore.Cursor
	prevCursor         datastore.Cursor
	doneProcessingPage bool
}

func NewDatastoreIterator(q *datastore.Query, c appengine.Context) *DatastoreIterator {
	return &DatastoreIterator{
		query:              q,
		context:            c,
		nextCursor:         datastore.Cursor{},
		prevCursor:         datastore.Cursor{},
		doneProcessingPage: false,
	}
}

func (this *DatastoreIterator) LoadNext(e Entity) error {
	iter := this.query.Run(this.context)
	key, err := iter.Next(e)

	if err != nil && err != datastore.Done {
		return err
	}

	this.doneProcessingPage = err != nil && err == datastore.Done

	if this.HasNext() {
		return datastore.Done
	}

	cursor, err := iter.Cursor()
	if err != nil {
		return err
	}
	this.prevCursor = this.nextCursor
	this.nextCursor = cursor
	this.query = this.query.Start(cursor)

	if this.doneProcessingPage {
		return this.LoadNext(e)
	}

	e.SetKey(key)
	return nil
}

func (this *DatastoreIterator) HasNextPage() bool {
	return this.prevCursor != this.nextCursor
}

func (this *DatastoreIterator) HasNext() bool {
	return this.doneProcessingPage && !this.HasNextPage()
}

func (this *DatastoreIterator) Cursor() string {
	return this.nextCursor.String()
}

type Iterator interface {
	LoadNext(Entity) error
	HasNext() bool
	HasNextPage() bool
	Cursor() string
	//	LoadNextPage(interface{}) bool
	//	SkipPages(int) bool
	//	CurrentCursor() string
	//	PreviousCursor() string
}
