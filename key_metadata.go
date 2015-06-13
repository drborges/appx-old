package ds

import (
	"appengine"
	"appengine/datastore"
)

type KeyMetadata struct {
	Kind     string
	IntID    int64
	StringID string
}

func NewKey(c appengine.Context, p Persistable) *datastore.Key {
	metadata := p.KeyMetadata()
	return datastore.NewKey(c, metadata.Kind, metadata.StringID, metadata.IntID, p.ParentKey())
}

func ResolveKey(c appengine.Context, p Persistable) error {
	if p.Key() != nil && !p.Key().Incomplete() {
		return nil
	}

	key := NewKey(c, p)
	if key.IntID() == 0 && key.StringID() == "" {
		return ErrUnresolvableKey
	}

	p.SetKey(key)
	return nil
}
