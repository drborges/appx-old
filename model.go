package appx

import "appengine/datastore"

type Model struct {
	key       *datastore.Key `datastore:"-"`
	parentKey *datastore.Key `datastore:"-"`
}

func (this Model) Key() *datastore.Key {
	return this.key
}

func (this *Model) SetKey(key *datastore.Key) {
	this.key = key
	this.parentKey = key.Parent()
}

func (this Model) ParentKey() *datastore.Key {
	return this.parentKey
}

func (this *Model) SetParentKey(key *datastore.Key) {
	this.parentKey = key
}

func (this *Model) ResourceID() string {
	return this.Key().Encode()
}

func (this *Model) SetResourceID(id string) error {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		return err
	}
	this.SetKey(key)
	this.SetParentKey(key.Parent())
	return nil
}
