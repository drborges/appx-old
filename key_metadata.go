package ds

import (
	"appengine"
	"appengine/datastore"
)

type KeyMetadata struct {
	Kind      string
	IntID     int64
	StringID  string
	ParentKey *KeyMetadata
}

func NewKey(c appengine.Context, metadata *KeyMetadata) (*datastore.Key, error) {
	var parentKey *datastore.Key
	var err error

	if metadata.ParentKey != nil {
		if parentKey, err = NewKey(c, metadata); err != nil {
			return nil, err
		}
	}

	return datastore.NewKey(c, metadata.Kind, metadata.StringID, metadata.IntID, parentKey), err
}

func ResolveKey(c appengine.Context, p Persistable) error {
	key := p.Key()
	if key == nil {
		var err error
		key, err = NewKey(c, p.KeyMetadata())

		if err != nil {
			return err
		}
	}

	p.SetKey(key)
	return nil
}