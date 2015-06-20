package appx

import "appengine/datastore"

func From(p Persistable) *datastore.Query {
	return datastore.NewQuery(p.KeyMetadata().Kind)
}
