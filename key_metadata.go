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

func NewKey(c appengine.Context, p Persistable) (*datastore.Key, error) {
	metadata := p.KeyMetadata()
	if metadata.HasParent && p.EntityParentKey() == nil {
		return nil, ErrMissingParentKey
	}

	return datastore.NewKey(c, metadata.Kind, metadata.StringID, metadata.IntID, p.EntityParentKey()), nil
}

func ResolveKey(c appengine.Context, p Persistable) error {
	if p.EntityKey() != nil && !p.EntityKey().Incomplete() {
		return nil
	}

	key, err := NewKey(c, p)
	if err != nil {
		return err
	}

	if key.Incomplete() {
		return ErrUnresolvableKey
	}

	p.SetEntityKey(key)
	return nil
}
