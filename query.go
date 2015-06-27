package appx

import "appengine/datastore"

func From(e Entity) *datastore.Query {
	return datastore.NewQuery(e.KeyMetadata().Kind)
}
