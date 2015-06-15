package ds

import (
	"appengine"
	"appengine/datastore"
	"reflect"
)

type Datastore struct {
	Context appengine.Context
}

func (this Datastore) Load(p Persistable) error {
	if err := ResolveKey(this.Context, p); err != nil {
		return err
	}

	return datastore.Get(this.Context, p.Key(), p)
}

func (this Datastore) LoadAll(slice interface{}) error {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return datastore.ErrInvalidEntityType
	}

	keys := make([]*datastore.Key, s.Len())
	for i := 0; i < s.Len(); i++ {
		p := s.Index(i).Interface().(Persistable)
		if err := ResolveKey(this.Context, p); err != nil {
			return err
		}
		keys[i] = p.Key()
	}

	return datastore.GetMulti(this.Context, keys, slice)
}

func (this Datastore) Update(p Persistable) error {
	if err := ResolveKey(this.Context, p); err != nil {
		return err
	}

	key, err := datastore.Put(this.Context, p.Key(), p)
	if err != nil {
		return err
	}

	p.SetKey(key)
	return nil
}

func (this Datastore) Create(p Persistable) error {
	key, err := NewKey(this.Context, p)
	if err != nil {
		return err
	}

	key, err = datastore.Put(this.Context, key, p)
	if err != nil {
		return err
	}

	p.SetKey(key)
	return nil
}

func (this Datastore) CreateAll(slice interface{}) error {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return datastore.ErrInvalidEntityType
	}

	keys := make([]*datastore.Key, s.Len())
	for i := 0; i < s.Len(); i++ {
		key, err := NewKey(this.Context, s.Index(i).Interface().(Persistable))
		if err != nil {
			return err
		}
		keys[i] = key
	}

	keys, err := datastore.PutMulti(this.Context, keys, slice)
	if err != nil {
		return err
	}

	for i, key := range keys {
		s.Index(i).Interface().(Entity).SetKey(key)
	}

	return nil
}

func (this Datastore) Delete(p Persistable) error {
	if err := ResolveKey(this.Context, p); err != nil {
		return err
	}

	return datastore.Delete(this.Context, p.Key())
}

func (this Datastore) Query(q *datastore.Query) QueryRunner {
	return QueryRunner{this.Context, q}
}