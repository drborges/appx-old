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
	if metadata.HasParent && p.ParentKey() == nil {
		return nil, ErrMissingParentKey
	}

	return datastore.NewKey(c, metadata.Kind, metadata.StringID, metadata.IntID, p.ParentKey()), nil
}

func ResolveKey(c appengine.Context, p Persistable) error {
	if p.Key() != nil && !p.Key().Incomplete() {
		return nil
	}

	key, err := NewKey(c, p)
	if err != nil {
		return err
	}

	if key.IntID() == 0 && key.StringID() == "" {
		return ErrUnresolvableKey
	}

	p.SetKey(key)
	return nil
}
