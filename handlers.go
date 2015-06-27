package appx

import (
	"appengine"
	"appengine/datastore"
)

func KeyResolver(c appengine.Context) EntityHandler {
	return func(e Entity) error {
		return ResolveKey(c, e)
	}
}

func DatastoreLoader(c appengine.Context) EntityHandler {
	return func(e Entity) error {
		return datastore.Get(c, e.EntityKey(), e)
	}
}
