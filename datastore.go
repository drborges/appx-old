package appx

import (
	"appengine"
	"appengine/datastore"
)

type Datastore struct {
	context   appengine.Context
}

func NewDatastore(c appengine.Context) *Datastore {
	return &Datastore{
		context: c,
	}
}

func (this *Datastore) Create(e Entity) error {
	return NewEntityHandlerChain().
		With(KeyAssinger(this.context)).
		With(DatastoreUpdater(this.context)).
		Handle(e)
}

func (this *Datastore) Load(e Entity) error {
	return NewEntityHandlerChain().
		With(KeyResolver(this.context)).
		With(DatastoreLoader(this.context)).
		Handle(e)
}

func (this *Datastore) Update(e Entity) error {
	return NewEntityHandlerChain().
		With(KeyResolver(this.context)).
		With(DatastoreUpdater(this.context)).
		Handle(e)
}

func (this *Datastore) Delete(e Entity) error {
	return NewEntityHandlerChain().
		With(KeyResolver(this.context)).
		With(DatastoreDeleter(this.context)).
		Handle(e)
}

func (this *Datastore) CreateAll(slice interface{}) error {
	return NewSliceHandlerChain().
		With(SliceValidator()).
		With(DatastoreBatchCreator(this.context)).
		Handle(slice)
}

func (this *Datastore) LoadAll(slice interface{}) error {
	return NewSliceHandlerChain().
		With(SliceValidator()).
		With(DatastoreBatchLoader(this.context)).
		Handle(slice)
}

func (this *Datastore) Query(q *datastore.Query) *QueryRunner {
	return &QueryRunner{this.context, q}
}
