package appx

import (
	"appengine"
	"appengine/datastore"
	"reflect"
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

func ResolveAllKeys(c appengine.Context, s reflect.Value) ([]*datastore.Key, error) {
	keys := make([]*datastore.Key, s.Len())
	for i := 0; i < s.Len(); i++ {
		p := s.Index(i).Interface().(Entity)
		if err := ResolveKey(c, p); err != nil {
			return nil, err
		}
		keys[i] = p.EntityKey()
	}

	return keys, nil
}
