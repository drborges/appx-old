package ds

import (
	"appengine"
	"appengine/datastore"
)

type Datastore struct {
	Context appengine.Context
}

func (this Datastore) Load(p Persistable) error {
	if err := ResolveKey(this.Context, p); err != nil {
		return err
	}

	return datastore.Get(this.Context, p.Key(), p)
}

func (this Datastore) Update(p Persistable) error {
	if err := ResolveKey(this.Context, p); err != nil {
		return err
	}

	key, err := datastore.Put(this.Context, p.Key(), p)

	if err != nil {
		return err
	}

	p.SetKey(key)
	return nil
}

func (this Datastore) Create(p Persistable) error {
	key, err := datastore.Put(this.Context, NewKey(this.Context, p), p)

	if err != nil {
		return err
	}

	p.SetKey(key)
	return nil
}

func (this Datastore) Delete(p Persistable) error {
	if err := ResolveKey(this.Context, p); err != nil {
		return err
	}

	return datastore.Delete(this.Context, p.Key())
}
