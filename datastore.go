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

func (this Datastore) Put(p Persistable) error {
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
