package appx

import (
	"appengine"
	"appengine/datastore"
	"reflect"
)

func KeyResolver(c appengine.Context) EntityHandler {
	return func(e Entity) error {
		return ResolveKey(c, e)
	}
}

func KeyAssinger(c appengine.Context) EntityHandler {
	return func(e Entity) error {
		key, err := NewKey(c, e)
		if err != nil {
			return err
		}
		e.SetEntityKey(key)
		return nil
	}
}

func DatastoreLoader(c appengine.Context) EntityHandler {
	return func(e Entity) error {
		return datastore.Get(c, e.EntityKey(), e)
	}
}

func DatastoreDeleter(c appengine.Context) EntityHandler {
	return func(e Entity) error {
		return datastore.Delete(c, e.EntityKey())
	}
}

func DatastoreUpdater(c appengine.Context) EntityHandler {
	return func(e Entity) error {
		key, err := datastore.Put(c, e.EntityKey(), e)
		if err != nil {
			return err
		}

		e.SetEntityKey(key)
		return nil
	}
}

func DatastoreBatchLoader(c appengine.Context) SliceHandler {
	return func(slice interface{}) error {
		s := reflect.ValueOf(slice)
		keys, err := ResolveAllKeys(c, s)
		if err != nil {
			return err
		}
		return datastore.GetMulti(c, keys, slice)
	}
}

func DatastoreBatchCreator(c appengine.Context) SliceHandler {
	return func(slice interface{}) error {
		s := reflect.ValueOf(slice)
		keys := make([]*datastore.Key, s.Len())
		for i := 0; i < s.Len(); i++ {
			key, err := NewKey(c, s.Index(i).Interface().(Entity))
			if err != nil {
				return err
			}
			keys[i] = key
		}

		keys, err := datastore.PutMulti(c, keys, slice)
		if err != nil {
			return err
		}

		for i, key := range keys {
			s.Index(i).Interface().(Entity).SetEntityKey(key)
		}

		return nil
	}
}

func SliceValidator() SliceHandler {
	return func (slice interface{}) error {
		s := reflect.ValueOf(slice)

		if s.Kind() != reflect.Slice {
			return datastore.ErrInvalidEntityType
		}

		return nil
	}
}
