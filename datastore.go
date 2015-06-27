package appx

import (
	"appengine"
	"appengine/datastore"
	"reflect"
)

type Datastore struct {
	context   appengine.Context
}

func NewDatastore(c appengine.Context) *Datastore {
	return &Datastore{
		context: c,
	}
}

func (this *Datastore) Load(e Entity) error {
	return NewEntityHandlerChain().
		With(KeyResolver(this.context)).
		With(DatastoreLoader(this.context)).
		Handle(e)
}

func (this *Datastore) LoadAll(slice interface{}) error {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return datastore.ErrInvalidEntityType
	}

	keys, err := ResolveAllKeys(this.context, s)
	if err != nil {
		return err
	}

	return datastore.GetMulti(this.context, keys, slice)
}

func (this *Datastore) Update(e Entity) error {
	if err := ResolveKey(this.context, e); err != nil {
		return err
	}

	key, err := datastore.Put(this.context, e.EntityKey(), e)
	if err != nil {
		return err
	}

	e.SetEntityKey(key)
	return nil
}

func (this *Datastore) Create(e Entity) error {
	key, err := NewKey(this.context, e)
	if err != nil {
		return err
	}

	key, err = datastore.Put(this.context, key, e)
	if err != nil {
		return err
	}

	e.SetEntityKey(key)
	return nil
}

func (this *Datastore) CreateAll(slice interface{}) error {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return ErrInvalidSliceType
	}

	keys := make([]*datastore.Key, s.Len())
	for i := 0; i < s.Len(); i++ {
		key, err := NewKey(this.context, s.Index(i).Interface().(Entity))
		if err != nil {
			return err
		}
		keys[i] = key
	}

	keys, err := datastore.PutMulti(this.context, keys, slice)
	if err != nil {
		return err
	}

	for i, key := range keys {
		s.Index(i).Interface().(Entity).SetEntityKey(key)
	}

	return nil
}

func (this *Datastore) Delete(e Entity) error {
	if err := ResolveKey(this.context, e); err != nil {
		return err
	}

	return datastore.Delete(this.context, e.EntityKey())
}

func (this *Datastore) Query(q *datastore.Query) *QueryRunner {
	return &QueryRunner{this.context, q}
}
