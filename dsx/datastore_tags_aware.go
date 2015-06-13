package dsx

import (
	"appengine/datastore"
	"github.com/drborges/ds"
	"appengine"
)

type DatastoreTagsAware struct {
	Context appengine.Context
}

func (this DatastoreTagsAware) Load(e ds.Entity) error {
	if err := ds.ResolveKey(this.Context, TaggedModel{e}); err != nil {
		return err
	}
	return datastore.Get(this.Context, e.Key(), e)
}

func (this DatastoreTagsAware) Create(e ds.Entity) error {
	key, err := datastore.Put(this.Context, ds.NewKey(this.Context, TaggedModel{e}), e)

	if err != nil {
		return err
	}

	e.SetKey(key)
	return nil
}

func (this DatastoreTagsAware) Update(e ds.Entity) error {
	if err := ds.ResolveKey(this.Context, TaggedModel{e}); err != nil {
		return err
	}

	key, err := datastore.Put(this.Context, e.Key(), e)

	if err != nil {
		return err
	}

	e.SetKey(key)
	return nil
}

func (this DatastoreTagsAware) Delete(e ds.Entity) error {
	if err := ds.ResolveKey(this.Context, TaggedModel{e}); err != nil {
		return err
	}

	return datastore.Delete(this.Context, e.Key())
}
