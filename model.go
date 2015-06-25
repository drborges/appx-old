package appx

import "appengine/datastore"

type Model struct {
	Key       *datastore.Key `datastore:"-"`
	ParentKey *datastore.Key `datastore:"-"`
}

func (this Model) EntityKey() *datastore.Key {
	return this.Key
}

func (this *Model) SetEntityKey(key *datastore.Key) {
	this.Key = key
	this.ParentKey = key.Parent()
}

func (this Model) EntityParentKey() *datastore.Key {
	return this.ParentKey
}

func (this *Model) SetEntityParentKey(key *datastore.Key) {
	this.ParentKey = key
}

func (this *Model) EncodedKey() string {
	return this.EntityKey().Encode()
}

func (this *Model) SetEncodedKey(id string) error {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		return err
	}
	this.SetEntityKey(key)
	this.SetEntityParentKey(key.Parent())
	return nil
}
