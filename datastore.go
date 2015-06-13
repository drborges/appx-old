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
