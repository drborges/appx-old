package appx

import (
	"appengine"
	"appengine/datastore"
)

type KeyMetadata struct {
	Kind      string
	IntID     int64
	StringID  string
	HasParent bool
}

func NewKey(c appengine.Context, e Entity) (*datastore.Key, error) {
	metadata := e.KeyMetadata()
	if metadata.HasParent && e.EntityParentKey() == nil {
		return nil, ErrMissingParentKey
	}

	return datastore.NewKey(c, metadata.Kind, metadata.StringID, metadata.IntID, e.EntityParentKey()), nil
}

func ResolveKey(c appengine.Context, e Entity) error {
	if e.EntityKey() != nil && !e.EntityKey().Incomplete() {
		return nil
	}

	key, err := NewKey(c, e)
	if err != nil {
		return err
	}

	if key.Incomplete() {
		return ErrUnresolvableKey
	}

	e.SetEntityKey(key)
	return nil
}
