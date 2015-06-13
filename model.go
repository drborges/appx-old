package ds

import "appengine/datastore"

type Model struct {
	key *datastore.Key
}

func (this Model) Key() *datastore.Key {
	return this.key
}

func (this *Model) SetKey(key *datastore.Key) {
	this.key = key
}
