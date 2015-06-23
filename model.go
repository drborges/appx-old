package appx

import "appengine/datastore"

type Model struct {
	key       *datastore.Key `datastore:"-"`
	parentKey *datastore.Key `datastore:"-"`
}

func (this Model) EntityKey() *datastore.Key {
	return this.key
}

func (this *Model) SetEntityKey(key *datastore.Key) {
	this.key = key
	this.parentKey = key.Parent()
}

func (this Model) EntityParentKey() *datastore.Key {
	return this.parentKey
}

func (this *Model) SetEntityParentKey(key *datastore.Key) {
	this.parentKey = key
}

func (this *Model) ResourceID() string {
	return this.EntityKey().Encode()
}

func (this *Model) SetResourceID(id string) error {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		return err
	}
	this.SetEntityKey(key)
	this.SetEntityParentKey(key.Parent())
	return nil
}
