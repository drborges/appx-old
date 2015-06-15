package ds

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
}

func (this Model) ParentKey() *datastore.Key {
	return this.parentKey
}

func (this *Model) SetParentKey(key *datastore.Key) {
	this.parentKey = key
}

func (this *Model) ID() string {
	return this.Key().Encode()
}

func (this *Model) SetID(uuid string) error {
	key, err := datastore.DecodeKey(uuid)
	if err != nil {
		return err
	}
	this.SetKey(key)
	this.SetParentKey(key.Parent())
	return nil
}